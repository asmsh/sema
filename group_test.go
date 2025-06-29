package sema_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/asmsh/sema"
)

func helperTestGroupRace() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	sg := sema.NewGroup(1)

	sg.Reserve()

	go func() {
		if sg.ReserveN(ctx.Done(), 1) {
			sg.Free()
		}
	}()
	go func() {
		if sg.ReserveN(ctx.Done(), 1) {
			sg.Free()
		}
	}()

	context.AfterFunc(ctx, func() { sg.Free() })

	sg.Wait()
}

func TestGroupRace(t *testing.T) {
	t.Parallel()
	n := 100
	for range n {
		helperTestGroupRace()
	}
}

func helperTestGroupRacePanic(t *testing.T) bool {
	numGoroutines := 5
	recoverChan := make(chan any, numGoroutines)
	handleRecover := func() {
		v := recover()
		if v != nil && v != "sema.Group: negative group counter" {
			t.Errorf("Free misue went unnoticed: %v", v)
		}
		recoverChan <- v
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	sg := sema.NewGroup(1)

	sg.Reserve()

	go func() {
		defer handleRecover()
		if sg.ReserveN(ctx.Done(), 1) {
			sg.Free()
		}
	}()

	// only one of the below Free calls below should succeed,
	// the rest should panic.
	for range numGoroutines - 1 {
		context.AfterFunc(ctx, func() {
			defer handleRecover()
			sg.Free()
		})
	}

	sg.Wait()

	// retrieve the panic results.
	var panics []any
	for range numGoroutines {
		v := <-recoverChan
		if v == nil {
			continue
		}
		panics = append(panics, v)
	}

	// check the results for the expected number of panics.
	// note: we are expecting 2 successes.
	success := false
	if len(panics) == numGoroutines-2 {
		success = true
	}

	return success
}

func TestGroupRacePanic(t *testing.T) {
	t.Parallel()
	success := true
	for range 100 {
		ret := helperTestGroupRacePanic(t)
		success = success && ret
	}
	if !success {
		t.Errorf("Free failed when racing")
	}
}

// merely returning from it without panics is a success.
func helperTestGroupRaceWait(t *testing.T) {
	numGoroutines := 10

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	sg := sema.NewGroup(numGoroutines)

	sg.ReserveN(nil, numGoroutines)

	// arrange to free the reserved N, 1 by 1, once ctx is closed.
	for range numGoroutines {
		context.AfterFunc(ctx, func() { sg.Free() })
	}

	// arrange for multiple goroutines to wait for sg,
	// competing with the arranged free calls above.
	// it could be really any number of goroutines, but sticking to n.
	wg := sync.WaitGroup{}
	wg.Add(numGoroutines)
	for range numGoroutines {
		context.AfterFunc(ctx, func() {
			defer wg.Done()

			sg.Wait()
		})
	}

	// wait for all ongoing sg.Wait calls, without pre-initiating the WaitChan.
	wg.Wait()
}

// merely returning from it without any panics is a success.
func TestGroupRaceWait(t *testing.T) {
	t.Parallel()
	for range 100 {
		helperTestGroupRaceWait(t)
	}
}

func helperGroupReserveFree(wg *sync.WaitGroup, sg *sema.Group) {
	defer wg.Done()
	sg.Reserve()
	go func() {
		time.Sleep(10 * time.Millisecond)
		sg.Free()
	}()
}

func helperGroupTryReserveFree(wg *sync.WaitGroup, sg *sema.Group) bool {
	defer wg.Done()
	if sg.TryReserveN(1) {
		go func() {
			time.Sleep(10 * time.Millisecond)
			sg.Free()
		}()
		return true
	}
	return false
}

func TestGroupFlow(t *testing.T) {
	t.Parallel()

	n := 10
	wg := sync.WaitGroup{}
	sg := sema.NewGroup(n)

	wg.Add(n)
	for range n {
		go helperGroupReserveFree(&wg, sg)
	}
	wg.Wait()

	wg.Add(1)
	go func() {
		if helperGroupTryReserveFree(&wg, sg) {
			t.Errorf("TryReserve() should be false")
		}
	}()
	wg.Wait()

	sg.Wait()

	wg.Add(1)
	go func() {
		if !helperGroupTryReserveFree(&wg, sg) {
			t.Errorf("TryReserve() should be true")
		}
	}()
	wg.Wait()

	sg.Wait()
}

func TestGroupCounters(t *testing.T) {
	t.Parallel()
	n := 10

	t.Run("the size has to be equal to init size", func(t *testing.T) {
		sg := sema.NewGroup(n)

		if sg.Size() != n {
			t.Errorf("Group size should be %d, got %d", n, sg.Size())
		}
	})

	t.Run("the counters must be zero before any usage", func(t *testing.T) {
		sg := sema.NewGroup(n)

		if active := sg.ActiveCount(); active != 0 {
			t.Errorf("Group active count should be 0, got %d", active)
		}
		if pending := sg.PendingCount(); pending != 0 {
			t.Errorf("Group pending count should be 0, got %d", pending)
		}
	})

	t.Run("the counters must be accurate after non-blocking reserve", func(t *testing.T) {
		sg := sema.NewGroup(n)

		sg.Reserve()
		if active := sg.ActiveCount(); active != 1 {
			t.Errorf("Group active count should be 1, got %d", active)
		}
		if pending := sg.PendingCount(); pending != 0 {
			t.Errorf("Group pending count should be 0, got %d", pending)
		}

		sg.Free()
		if active := sg.ActiveCount(); active != 0 {
			t.Errorf("Group active count should be 0, got %d", active)
		}
		if pending := sg.PendingCount(); pending != 0 {
			t.Errorf("Group pending count should be 0, got %d", pending)
		}
	})

	t.Run("the counters must be accurate after blocking reserve", func(t *testing.T) {
		sg := sema.NewGroup(n)

		// take 1 to make the below ReserveN blocks.
		sg.Reserve()

		reserveNDone := make(chan struct{})
		freeReady := make(chan struct{})
		freeDone := make(chan struct{})
		go func() {
			sg.ReserveN(nil, n) // this blocks.

			close(reserveNDone)
			<-freeReady

			sg.FreeN(n)

			close(freeDone)
		}()

		// wait until the above ReserveN blocks.
		for sg.TryReserveN(1) {
			sg.Free()
		}

		if active := sg.ActiveCount(); active != 1 {
			t.Errorf("Group active count should be %d, got %d", 1, active)
		}
		if pending := sg.PendingCount(); pending != n {
			t.Errorf("Group pending count should be %d, got %d", n, pending)
		}

		// wake up the blocked ReserveN call.
		sg.Free()

		// wait until the Free wakes up the blocked ReserveN.
		<-reserveNDone

		if active := sg.ActiveCount(); active != n {
			t.Errorf("Group active count should be %d, got %d", n, active)
		}
		if pending := sg.PendingCount(); pending != 0 {
			t.Errorf("Group pending count should be %d, got %d", 0, pending)
		}

		// unblock the Free call and wait for it to finish.
		close(freeReady)
		<-freeDone

		if active := sg.ActiveCount(); active != 0 {
			t.Errorf("Group active count should be %d, got %d", 0, active)
		}
		if pending := sg.PendingCount(); pending != 0 {
			t.Errorf("Group pending count should be %d, got %d", 0, pending)
		}
	})

	t.Run("the counters must be as expected after canceled reserve", func(t *testing.T) {
		sg := sema.NewGroup(n)

		sg.Reserve()

		// this blocks until the ReserveN below is executed,
		// then closes the doneChan and returns.
		doneChan := make(chan struct{})
		go func() {
			// wait until the ReserveN blocks.
			for sg.TryReserveN(1) {
				sg.Free()
			}

			// since the ReserveN is executed, and there's no room
			// for it, it must block.
			if active := sg.ActiveCount(); active != 1 {
				t.Errorf("Group active count should be 1, got %d", active)
			}
			// at least n,
			if pending := sg.PendingCount(); pending < n || pending > 2*n {
				t.Errorf("Group pending count should be %d, got %d", n, pending)
			}

			close(doneChan)

			sg.Free()
		}()

		// this blocks until the ReserveN below is executed,
		// which triggers the above TryReserveN to return,
		// which causes the doneChan to be closed.
		go func() {
			// wait for the TryReserveN to return.
			<-doneChan

			// make sure the counter can only go down.
			if active := sg.ActiveCount(); active > 1 {
				t.Errorf("Group active count should not be more than 1, got %d", active)
			}
		}()

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()

			if got := sg.ReserveN(doneChan, n); got != false {
				t.Errorf("ReserveN with closed done chan returned wrong result: want false, got %v", got)
			}
		}()

		if got := sg.ReserveN(doneChan, n); got != false {
			t.Errorf("ReserveN with closed done chan returned wrong result: want false, got %v", got)
		}

		// wait for the ReserveN to execute, block, abort, and return.
		wg.Wait()

		if active := sg.ActiveCount(); active != 0 {
			t.Errorf("Group active count should be 0, got %d", active)
		}
		if pending := sg.PendingCount(); pending != 0 {
			t.Errorf("Group pending count should be 0, got %d", pending)
		}
	})

	t.Run("the counters must be as expected after wait calls", func(t *testing.T) {
		sg := sema.NewGroup(n)

		sg.Reserve()

		ctx, cancel := context.WithCancel(context.Background())
		context.AfterFunc(ctx, func() {
			sg.Free()
		})

		go func() {
			for sg.TryReserveN(1) {
				sg.FreeN(1)
			}

			cancel()
		}()

		if got := sg.ReserveN(ctx.Done(), n); got != false {
			t.Errorf("ReserveN with canceled context chan returned wrong result: want false, got %v", got)
		}

		sg.Wait()

		if active := sg.ActiveCount(); active != 0 {
			t.Errorf("Group active count should be 0, got %d", active)
		}
		if pending := sg.PendingCount(); pending != 0 {
			t.Errorf("Group pending count should be 0, got %d", pending)
		}
	})
}
