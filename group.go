// Copyright 2025 Ahmad Sameh(asmsh)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sema

import (
	"runtime"
	"sync/atomic"
)

// Group guards concurrent access to a resource by providing methods to
// control concurrency, observe usage counters, and wait for in-flight
// operations to complete.
//
// The zero value is a ready to use [Group], with no concurrency limit,
// comparable to a [sync.WaitGroup].
//
// Group is best suited for concurrent tasks of equal weight or cost.
type Group struct {
	// waitChan is created lazily in Wait, only if it hasn't already,
	// and there are some active calls.
	// it's an unbuffered channel and is closed once the Group zeros.
	waitChan atomic.Value // chan struct{}

	// blockChan is created lazily in [NewGroup] or [Group.SetSize],
	// only if size > 0.
	// it's an unbuffered channel that's never closed.
	blockChan atomic.Value // chan struct{}

	// high 32 bits are pending count, low 32 bits are active count.
	counter atomic.Uint64

	// size is the maximum number of N that can be reserved.
	size atomic.Uint32
}

// NewGroup creates a new [Group] with the provided size.
// The [Group] size is the concurrency limit that it can handle.
func NewGroup(size int) *Group {
	g := &Group{}
	g.setSize(size)
	return g
}

func (g *Group) setSize(size int) {
	// normalize negative size to 0.
	if size < 0 {
		size = 0
	}
	// only create the blockChan if size is not zero.
	if size > 0 {
		// make sure the size isn't too big.
		s := uint32(size)
		if int(s) != size {
			panic("sema.Group: incorrect group size")
		}

		// save the size and create the block chan.
		g.size.Store(s)
		g.blockChan.Store(make(chan struct{}))
	}
}

func counterParts(counter uint64) (pending uint32, active int32) {
	return uint32(counter >> 32), int32(counter)
}

func (g *Group) counterUpdate(
	oldCounter uint64,
	pendingDelta int,
	activeDelta int,
) (newCounter uint64, ok bool) {
	oldPending, oldActive := counterParts(oldCounter)

	newPending := oldPending + uint32(pendingDelta)
	newActive := oldActive + int32(activeDelta)

	// Check for overflow/underflow from low
	if activeDelta > 0 && newActive < oldActive {
		newPending += 1 // Carry
	} else if activeDelta < 0 && newActive > oldActive {
		newPending -= 1 // Borrow
	}

	newCounter = uint64(newPending)<<32 | uint64(uint32(newActive))

	if g.counter.CompareAndSwap(oldCounter, newCounter) {
		return newCounter, true
	}

	return newCounter, false
}

// SetSize sets the [Group.Size] to the passed value.
//
// It panics if it's called on a non-zero [Group].
// It panics if the [Group.Size] size was set before, either through this
// method or through the [NewGroup] function.
// This means that it should be called once before any other methods, and
// only on the zero value of a [Group].
//
// If size is zero or negative, then the [Group] has no limits, and no
// [Group.Reserve] or [Group.ReserveN] calls will block.
//
// Calling it with 0 on a zero [Group] has no effect.
func (g *Group) SetSize(size int) {
	// check if the Group has already started counting.
	if g.counter.Load() != 0 {
		panic("sema.Group: concurrent Reserve calls while initializing group")
	}

	// check if the Group is already initialized.
	if g.blockChan.Load() != nil {
		panic("sema.Group: group already initialized")
	}

	g.setSize(size)

	// re-check if the Group has already started counting.
	if g.counter.Load() != 0 {
		panic("sema.Group: concurrent Reserve calls while initializing group")
	}
}

// Size is the current limit of this [Group], which is the maximum
// N resources allowed to be active at the same time.
//
// The [Group] size is the concurrency limit that it can handle.
// If it's 0, then the [Group] has no limit.
func (g *Group) Size() int {
	return int(g.size.Load())
}

// ActiveCount is the total number of successfully reserved N resources
// via calling either [Group.Reserve] or [Group.TryReserve].
// It represents the number of N that's currently used from this [Group]'s
// size, which can never be greater than size.
func (g *Group) ActiveCount() int {
	_, active := counterParts(g.counter.Load())
	return int(active)
}

// PendingCount is the total number of N resources that's pending and
// blocking their [Group.Reserve] calls, waiting for matching [Group.Free]
// calls to unblock them.
func (g *Group) PendingCount() int {
	pending, _ := counterParts(g.counter.Load())
	return int(pending)
}

// Reserve increments [Group.ActiveCount] by 1, blocking if needed until
// there's room made available by [Group.Free] or [Group.FreeN] calls.
//
// It returns immediately if the [Group.PendingCount] is 0, and there's
// available room for 1 (in [Group.ActiveCount] against the [Group.Size]).
// Otherwise, it increments the [Group.PendingCount] instead while being
// blocked, until there's room.
//
// If there's a room made available, it wakes up in random order with other
// blocked [Group.Reserve] and [Group.ReserveN] calls.
//
// It always returns immediately, without blocking, if the [Group.Size] is 0.
//
// It always updates the [Group.ActiveCount] and [Group.PendingCount] before
// returning.
func (g *Group) Reserve() {
	g.ReserveN(nil, 1)
}

// ReserveN increments [Group.ActiveCount] by n, blocking if needed, as long
// as n is within [Group.Size], until there's room made available by [Group.Free]
// or [Group.FreeN] calls, and returns true if the increment was successful,
// or false if it was aborted via the provided doneChan.
//
// It returns immediately if the [Group.PendingCount] is 0, and there's
// available room for n (in [Group.ActiveCount] against the [Group.Size]).
// Otherwise, it increments the [Group.PendingCount] instead while being
// blocked, until there's room, or aborts if the provided doneChan becomes
// receive-ready if it's non-nil.
//
// If there's a room made available, it wakes up in random order with other
// blocked [Group.Reserve] and [Group.ReserveN] calls.
//
// It always returns immediately, without blocking, if the [Group.Size] is 0.
//
// It always updates the [Group.ActiveCount] and [Group.PendingCount] before
// returning.
// Such that, if it returns false, the [Group.PendingCount] will no longer
// include the provided n.
// This means that, once the provided doneChan becomes receive-ready,
// the provided n will not move to the [Group.ActiveCount], and will be
// removed from the [Group.PendingCount] before returning.
//
// It panics if n is less than or equal to 0.
//
// Note: The doneChan becomes receive-ready when it's closed or sent to.
func (g *Group) ReserveN(doneChan <-chan struct{}, n int) (reserved bool) {
	if n <= 0 {
		// n can't be 0 or negative, as 0 isn't a valid resource,
		// and for negative values, [Group.FreeN] should be used.
		panic("sema.Group: invalid group reserve N value")
	}

	// return and don't update any counters if the provided doneChan
	// is already closed.
	if doneChan != nil {
		select {
		case <-doneChan:
			return false
		default:
		}
	}

	// if the size is 0, then there's no limitation on the Reserve
	// calls, and the call should succeed right away.
	size := g.size.Load()
	if size == 0 {
		for ok := false; !ok; {
			_, ok = g.counterUpdate(g.counter.Load(), 0, n)
		}

		return true
	}

	// if the requested N is greater than the set size, then this
	// Reserve call is destined to fail, so wait for the done chan,
	// if it's provided, and return failure.
	if n > int(size) {
		if doneChan != nil {
			<-doneChan
		}

		return false
	}

	// if the Reserve call can be made with the size limit, then
	// the call should succeed right away.
	if g.tryReserve(size, n, false) {
		return true
	}

	// otherwise, block until matching FreeN calls are made.
	return g.reserveNSlow(size, doneChan, n)
}

func (g *Group) reserveNSlow(size uint32, doneChan <-chan struct{}, reserveN int) bool {
	// at this point, the blockChan shouldn't be nil.
	// because we enter this method only if the size is not 0,
	// and the blockChan will always be set if the size is not 0,
	// unless the [Group.SetSize] method has been called concurrently
	// with the [Group.ReserveN] method, and a context switch to
	// this method happened before the [Group.SetSize] method has set
	// the blockChan value, but right after setting the size value.
	blockChan := g.blockChan.Load()

	if blockChan == nil {
		panic("sema.Group: concurrent Reserve calls while initializing group")
	}

	blockChanVal := blockChan.(chan struct{})

	// wait for a suitable freed tickets, or keep looping.
	for {
		select {
		case <-blockChanVal:
			// block for a FreeN call.
			reloop, ok := g.reserveNSuccessWait(size, doneChan, reserveN, blockChanVal)
			if ok {
				return true
			}
			if !reloop {
				return false
			}
		case <-doneChan:
			// or abort on demand.
			g.reserveNAbortWait(blockChanVal, reserveN)
			return false
		}
	}
}

func (g *Group) reserveNSuccessWait(
	size uint32,
	doneChan <-chan struct{},
	reserveN int,
	blockChan chan struct{},
) (reloop, ok bool) {
	// execute in a loop, because counterUpdate might lose the CAS.
	for {
		counter := g.counter.Load()
		pending, active := counterParts(counter)
		diffN := int(size) - int(active) - reserveN

		// if we got what we need, update the counter and return true.
		if diffN >= 0 {
			// only move forward if the doneChan wasn't closed.
			select {
			case <-doneChan:
				g.reserveNAbortWait(blockChan, reserveN)
				return false, false
			default:
			}

			_, ok = g.counterUpdate(counter, -reserveN, reserveN)
			if !ok {
				// the counter got changed, re-loop and try again.
				continue
			}

			return false, true
		}

		// this call still needs more N to succeed...
		// if this is the only blocked call, then return and wait for
		// the next free call.
		if int(pending) == reserveN {
			return true, false
		}

		// if there are still potentially active calls, then return and wait
		// for the next free call.
		if int(pending)+int(active) > int(size) {
			return true, false
		}

		// notify another blocked call to check its goal...
		// if we didn't notify any reserve call successfully, it means that
		// we unblocked another call, which means that the other call already
		// updated the counter, so we need to re-loop and check it again,
		// because if we return and wait for a FreeN call, it might never
		// come because this call might be the last one.
		//_, notifiedReserve := g.notifyReserve(blockChan, reserveN)
		//if !notifiedReserve {
		//	continue
		//}

		// only proceed if there are other pending calls.
		if int(pending)-reserveN <= 0 {
			continue
		}

		// attempt to wake up a Reserve call, or unblock another that's trying the same.
		select {
		case blockChan <- struct{}{}:
			// woke up 1 blocked Reserve call...
			return true, false
		case <-blockChan:
			// unblock another call.
			continue
		}
	}
}

func (g *Group) notifyReserve(
	blockChan chan struct{},
	excludeN int,
) (counter uint64, notified bool) {
	counter = g.counter.Load()
	pending, _ := counterParts(counter)

	// only proceed if there are other pending calls.
	if int(pending)-excludeN <= 0 {
		return counter, false
	}

	// attempt to wake up a Reserve call, or unblock another that's trying the same.
	select {
	case blockChan <- struct{}{}:
		// woke up 1 blocked Reserve call...
		return counter, true
	case <-blockChan:
		return counter, false
	}
}

func (g *Group) reserveNAbortWait(blockChan chan struct{}, reserveN int) {
	for ok := false; !ok; {
		counter := g.counter.Load()
		counter, ok = g.counterUpdate(counter, -reserveN, 0)
	}

	counter := g.notifyFree(blockChan)
	g.notifyWait(counterParts(counter))
}

// TryReserveN tries to increment the [Group.ActiveCount] by n without
// blocking and returns true if it was successful.
// It returns false if the [Group.PendingCount] is not 0 or there's no
// room in the [Group.ActiveCount] against the [Group.Size].
//
// It always returns true if the [Group.Size] is 0.
//
// It panics if n is less than or equal to 0.
func (g *Group) TryReserveN(n int) bool {
	if n <= 0 {
		// n can't be 0 or negative, as 0 isn't a valid resource,
		// and for negative values, [Group.FreeN] should be used.
		panic("sema.Group: invalid group reserve N value")
	}

	size := g.size.Load()
	if size == 0 {
		// no limitation on the Reserve calls, so the call is allowed.
		for ok := false; !ok; {
			_, ok = g.counterUpdate(g.counter.Load(), 0, n)
		}

		return true
	}

	if n > int(size) {
		return false
	}

	return g.tryReserve(size, n, true)
}

func (g *Group) tryReserve(size uint32, reserveN int, tryCall bool) bool {
	for {
		counter := g.counter.Load()
		pending, active := counterParts(counter)

		// if the group has room and doesn't have any waiters
		if pending == 0 && int(active)+reserveN <= int(size) {
			_, ok := g.counterUpdate(counter, 0, reserveN)
			if ok {
				return true
			}
			continue
		} else if !tryCall {
			_, ok := g.counterUpdate(counter, reserveN, 0)
			if ok {
				return false
			}
			continue
		} else {
			return false
		}
	}
}

// Free decrements the [Group.ActiveCount] by 1, making it available for other
// reserve calls, and attempting to wake up a single blocked [Group.Reserve]
// or [Group.ReserveN] call, in random order, if there's any blocked.
//
// If the [Group.ActiveCount] reaches zero by this call, it will wake up any
// calls blocked on [Group.Wait], and closes the [Group.WaitChan].
//
// If the [Group.ActiveCount] goes below zero by this call, it panics.
func (g *Group) Free() {
	g.FreeN(1)
}

// FreeN decrements the [Group.ActiveCount] by n, making it available for other
// reserve calls, and attempting to wake up a single blocked [Group.Reserve]
// or [Group.ReserveN] call, in random order, if there's any blocked.
//
// If the [Group.ActiveCount] reaches zero by this call, it will wake up any
// calls blocked on [Group.Wait], and closes the [Group.WaitChan].
//
// If the [Group.ActiveCount] goes below zero by this call, it panics.
//
// It panics if n is less than or equal to 0.
func (g *Group) FreeN(n int) {
	if n <= 0 {
		// n can't be 0 or negative, as 0 isn't a valid resource,
		// and for negative values, [Group.FreeN] should be used.
		panic("sema.Group: invalid group free N value")
	}

	// check if we should wake up any blocked [Group.ReserveN] calls.
	// note: if the blockChan is nil, then pending must be 0, as
	// pending is used to track blocked calls, and blockChan is
	// only set when the Group can have blocked calls (size != 0).
	blockChan := g.blockChan.Load()

	// update the counter, and read its values.
	var active int32
	for ok := false; !ok; {
		counter := g.counter.Load()
		counter, ok = g.counterUpdate(counter, 0, -n)
		_, active = counterParts(counter)
	}

	// notify any blocked ReserveN calls of the counter update,
	// and make sure the counter values are updated.
	var pending uint32
	if blockChan != nil {
		counter := g.notifyFree(blockChan.(chan struct{}))
		pending, _ = counterParts(counter)
	}

	// attempt to wake up any blocked [Group.Wait] calls.
	g.notifyWait(pending, active)

	// handle any misuse, assuming valid usage so far.
	if active < 0 {
		panic("sema.Group: negative group counter")
	}
}

func (g *Group) notifyFree(blockChan chan struct{}) (counter uint64) {
	counter = g.counter.Load()
	pending, _ := counterParts(counter)

	// this will avoid Free missing an opportunity to wake up a Reserve.
	for int(pending) > 0 {
		// attempt to wakeup a Reserve call, or update pending until it's 0.
		select {
		case blockChan <- struct{}{}:
			// woke up 1 blocked Reserve call...
			counter = g.counter.Load()
			return counter
		default:
			// the Reserve call aborted, so re-sync.
			runtime.Gosched()

			counter = g.counter.Load()
			pending, _ = counterParts(counter)
		}
	}

	return counter
}

func (g *Group) notifyWait(pending uint32, active int32) {
	// if there are still blocked calls, as this means we still don't need
	// to wake up any [Group.Wait] calls.
	// return if there are still active calls, as this means we still don't
	// need to wake up any [Group.Wait] calls.
	if pending > 0 || active > 0 {
		return
	}

	// waitChan will be nil only if no Wait calls have been made.
	waitChan := g.waitChan.Load()
	if waitChan == nil || waitChan == nilChan {
		return
	}

	if !g.waitChan.CompareAndSwap(waitChan, nilChan) {
		// if it didn't succeed, return, as it means that this value is
		// an waitChan value, and a newer once has been set, after the old
		// one got already closed.
		return
	}

	close(waitChan.(chan struct{}))
}

// Wait blocks until the [Group] reaches zero.
// It waits only for [Group.Reserve] and [Group.ReserveN] calls that are
// either made before this function is called, or made while the [Group] is non-zero.
// A zero [Group] means both [Group.ActiveCount] and [Group.PendingCount] are zero.
func (g *Group) Wait() {
	waitChan := g.initWaitChan()
	<-waitChan
}

// WaitChan returns a channel that will be closed once the [Group] reaches zero.
// It waits only for [Group.Reserve] and [Group.ReserveN] calls that are
// either made before this function is called, or made while the [Group] is non-zero.
// A zero [Group] means both [Group.ActiveCount] and [Group.PendingCount] are zero.
func (g *Group) WaitChan() <-chan struct{} {
	return g.initWaitChan()
}

var closedChan = make(chan struct{})
var nilChan chan struct{}

func init() {
	close(closedChan)
}

func (g *Group) initWaitChan() chan struct{} {
	waitChan := g.waitChan.Load()

	pending, active := counterParts(g.counter.Load())
	if pending <= 0 && active <= 0 {
		return closedChan
	}

	if waitChan != nil && waitChan != nilChan {
		return waitChan.(chan struct{})
	}

	return g.initWaitChanSlow(waitChan)
}

func (g *Group) initWaitChanSlow(waitChan any) chan struct{} {
	// we entered this method with waitChan either nil or nilChan,
	// and it can only go to nilChan or a valid chan value.
	newWaitChan := make(chan struct{})

	if g.waitChan.CompareAndSwap(waitChan, newWaitChan) {
		// we need to be sure that the swapped waitChan will be closed by
		// a Free call, which happens only if the counter is still not 0.
		pending, active := counterParts(g.counter.Load())
		if pending <= 0 && active <= 0 {
			return closedChan
		}

		return newWaitChan
	}

	// the waitChan value got changed...
	// if it's not a nilChan, then it must be a valid chan value.
	// 1- if it's nilChan, then it must have been set by a concurrent Free
	// call.
	// meaning that, all Reserve calls that have been done before this Wait,
	// their Free calls have been executed, so we unblock these Wait calls,
	// even if there's now new Reserve calls.
	// 2- if it's a valid chan value, then it must have been set by another
	// competing Wait call.
	// note: from point no. 1 above, it means that all Reserve calls must
	// be executed before their respective Wait calls.
	// note: comparing only against nilChan, because it can't be nil
	// again at this point.
	loadedWaitChan := g.waitChan.Load()
	if loadedWaitChan == nilChan {
		return closedChan
	}

	return loadedWaitChan.(chan struct{})
}
