// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Note: this file is copied from: sync/waitgroup_test.go standard library file.
// It's been modified to run the same tests that are done on the [sync.WaitGroup] type,
// but on the [sema.Group] type instead.
// Also, the benchmarks are moved to a matching file in the benchmarks sub-package.

package sema_test

import (
	"sync/atomic"
	"testing"

	"github.com/asmsh/sema"
)

func testSemaGroup(t *testing.T, sg1 *sema.Group, sg2 *sema.Group) {
	n := 16
	sg1.ReserveN(nil, n)
	sg2.ReserveN(nil, n)
	exited := make(chan bool, n)
	for i := 0; i != n; i++ {
		go func() {
			sg1.Free()
			sg2.Wait()
			exited <- true
		}()
	}
	sg1.Wait()
	for i := 0; i != n; i++ {
		select {
		case <-exited:
			t.Fatal("sema.Group: released Group too soon")
		default:
		}
		sg2.Free()
	}
	for i := 0; i != n; i++ {
		<-exited // Will block if barrier fails to unlock someone.
	}
}

func TestSemaGroup(t *testing.T) {
	sg1 := &sema.Group{}
	sg2 := &sema.Group{}

	// Run the same test a few times to ensure barrier is in a proper state.
	for i := 0; i != 8; i++ {
		testSemaGroup(t, sg1, sg2)
	}
}

func TestSemaGroupMisuse(t *testing.T) {
	defer func() {
		err := recover()
		if err != "sema.Group: negative group counter" {
			t.Fatalf("Unexpected panic: %#v", err)
		}
	}()
	wg := &sema.Group{}
	wg.Reserve()
	wg.Free()
	wg.Free()
	t.Fatal("Should panic")
}

func TestSemaGroupRace(t *testing.T) {
	// Run this test for about 1ms.
	for i := 0; i < 1000; i++ {
		sg := &sema.Group{}
		n := new(int32)
		// spawn goroutine 1
		sg.ReserveN(nil, 1)
		go func() {
			atomic.AddInt32(n, 1)
			sg.Free()
		}()
		// spawn goroutine 2
		sg.ReserveN(nil, 1)
		go func() {
			atomic.AddInt32(n, 1)
			sg.Free()
		}()
		// Wait for goroutine 1 and 2
		sg.Wait()
		if atomic.LoadInt32(n) != 2 {
			t.Fatalf("Spurious wakeup from Wait @ i = %d", i)
		}
	}
}

func TestSemaGroupAlign(t *testing.T) {
	type X struct {
		x  byte
		sg sema.Group
	}
	var x X
	x.sg.Reserve()
	go func(x *X) {
		x.sg.Free()
	}(&x)
	x.sg.Wait()
}
