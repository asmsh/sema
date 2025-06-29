// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Note: this file is copied from: golang.org/x/sync/semaphore.semaphore_test.go file.
// It's been modified to run the same tests that are done on the [semaphore.Weighted] type,
// but on the [sema.Group] type instead.
// Also, the benchmarks are moved to a matching file in the benchmarks sub-package.

package sema_test

import (
	"context"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/asmsh/sema"
)

const maxSleep = 1 * time.Millisecond

func HammerGroup(sg *sema.Group, n int64, loops int) {
	for i := 0; i < loops; i++ {
		sg.ReserveN(context.Background().Done(), int(n))
		time.Sleep(time.Duration(rand.Int63n(int64(maxSleep/time.Nanosecond))) * time.Nanosecond)
		sg.FreeN(int(n))
	}
}

func TestGroup(t *testing.T) {
	t.Parallel()

	n := runtime.GOMAXPROCS(0)
	loops := 10000 / n
	sg := sema.NewGroup(n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 1; i <= n; i++ {
		i := i
		go func() {
			defer wg.Done()
			HammerGroup(sg, int64(i), loops)
		}()
	}
	wg.Wait()
}

func TestGroupPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("free of an unreserved semaphore group did not panic")
		}
	}()
	sg := sema.NewGroup(1)
	sg.FreeN(1)
}

func TestGroupTryReserveN(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	sem := sema.NewGroup(2)
	tries := []bool{}
	sem.ReserveN(ctx.Done(), 1)
	tries = append(tries, sem.TryReserveN(1))
	tries = append(tries, sem.TryReserveN(1))

	sem.FreeN(2)

	tries = append(tries, sem.TryReserveN(1))
	sem.ReserveN(ctx.Done(), 1)
	tries = append(tries, sem.TryReserveN(1))

	want := []bool{true, false, true, false}
	for i := range tries {
		if tries[i] != want[i] {
			t.Errorf("tries[%d]: got %t, want %t", i, tries[i], want[i])
		}
	}
}

func TestGroupReserve(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	sem := sema.NewGroup(2)
	tryReserve := func(n int) bool {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		defer cancel()
		return sem.ReserveN(ctx.Done(), n)
	}

	tries := []bool{}
	sem.ReserveN(ctx.Done(), 1)
	tries = append(tries, tryReserve(1))
	tries = append(tries, tryReserve(1))

	sem.FreeN(2)

	tries = append(tries, tryReserve(1))
	sem.ReserveN(ctx.Done(), 1)
	tries = append(tries, tryReserve(1))

	want := []bool{true, false, true, false}
	for i := range tries {
		if tries[i] != want[i] {
			t.Errorf("tries[%d]: got %t, want %t", i, tries[i], want[i])
		}
	}
}

func TestLargeReserveDoesntStarveMulti(t *testing.T) {
	t.SkipNow()
	t.Parallel()
	n := 100
	for range n {
		TestLargeReserveDoesntStarve(t)
	}
}

// TestLargeReserveDoesntStarve times out if a large call to Reserve starves.
// Merely returning from the test function indicates success.
func TestLargeReserveDoesntStarve(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	n := runtime.GOMAXPROCS(0)
	sg := sema.NewGroup(n)
	running := true

	var wg sync.WaitGroup
	wg.Add(n)
	for i := n; i > 0; i-- {
		sg.ReserveN(ctx.Done(), 1)
		go func() {
			defer func() {
				sg.FreeN(1)
				wg.Done()
			}()
			for running {
				time.Sleep(1 * time.Millisecond)
				sg.FreeN(1)
				sg.ReserveN(ctx.Done(), 1)
			}
		}()
	}

	sg.ReserveN(ctx.Done(), n)
	running = false
	sg.FreeN(n)
	wg.Wait()
}

func TestAllocCancelDoesntStarveMulti(t *testing.T) {
	t.SkipNow()
	t.Parallel()
	for range 100 {
		TestAllocCancelDoesntStarve(t)
	}
}

// translated from https://github.com/zhiqiangxu/util/blob/master/mutex/crwmutex_test.go#L43
func TestAllocCancelDoesntStarve(t *testing.T) {
	sg := sema.NewGroup(10)

	// Block off a portion of the semaphore so that ReserveN(_, 10) can eventually succeed.
	sg.ReserveN(context.Background().Done(), 1)

	// In the background, ReserveN(_, 10).
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		sg.ReserveN(ctx.Done(), 10)
	}()

	// Wait until the ReserveN(_, 10) call blocks.
	for sg.TryReserveN(1) {
		sg.FreeN(1)
		runtime.Gosched()
	}

	// Now try to grab a read lock, and simultaneously unblock the ReserveN(_, 10) call.
	// Both Reserve calls should unblock and return, in either order.
	go cancel()

	sg.Reserve()
	sg.FreeN(1)
}

func TestWeightedAcquireCanceled(t *testing.T) {
	// https://go.dev/issue/63615
	sg := sema.NewGroup(2)
	ctx, cancel := context.WithCancel(context.Background())
	sg.ReserveN(context.Background().Done(), 1)
	ch := make(chan struct{})
	go func() {
		// Synchronize with the ReserveN(2) below.
		for sg.TryReserveN(1) {
			sg.FreeN(1)
		}
		// Now cancel ctx, and then free the token.
		cancel()
		sg.FreeN(1)
		close(ch)
	}()
	// Since the context closing happens before enough tokens become available,
	// this ReserveN must fail.
	if ret := sg.ReserveN(ctx.Done(), 2); ret != false {
		t.Errorf("ReserveN with canceled context chan returned wrong result: want false, got %v", ret)
	}
	// There must always be two tokens in the semaphore after the other
	// goroutine releases the one we held at the start.
	<-ch
	if !sg.TryReserveN(2) {
		t.Fatal("TryReserveN after canceled ReserveN failed")
	}
	// Additionally verify that we don't acquire with a done context even when
	// we wouldn't need to block to do so.
	sg.FreeN(2)
	if ret := sg.ReserveN(ctx.Done(), 1); ret != false {
		t.Errorf("ReserveN with canceled context chan returned result: want false, got %v", ret)
	}
}
