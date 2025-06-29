// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Note: this file is copied from: sync/waitgroup_test.go standard library file.
// It's been modified to run the same benchmarks that are done on the [sync.WaitGroup] type,
// but on the [sema.Group] type instead.
// The benchmarks have been updated slightly to have an extra version where
// the SetParallelism is set to a high number.

package benchmarks

import (
	. "sync"
	"testing"

	"github.com/asmsh/sema"
)

const highParallelismNum = 10_000

func benchmarkWaitGroupUncontended(b *testing.B) {
	type PaddedWaitGroup struct {
		WaitGroup
		pad [128]uint8
	}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var wg PaddedWaitGroup
		for pb.Next() {
			wg.Add(1)
			wg.Done()
			wg.Wait()
		}
	})
}

func BenchmarkWaitGroupUncontended(b *testing.B) {
	benchmarks := []struct {
		name            string
		highParallelism bool
	}{
		{
			name: "Uncontended",
		},
		{
			name:            "Uncontended-HighParallelism",
			highParallelism: true,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			if bm.highParallelism {
				b.SetParallelism(highParallelismNum)
			}
			benchmarkWaitGroupUncontended(b)
		})
	}
}

func benchmarkSemaGroupUncontended(b *testing.B) {
	type PaddedSemaGroup struct {
		sema.Group
		pad [128]uint8
	}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var sg PaddedSemaGroup
		for pb.Next() {
			sg.Reserve()
			sg.Free()
			sg.Wait()
		}
	})
}

func BenchmarkSemaGroupUncontended(b *testing.B) {
	benchmarks := []struct {
		name            string
		highParallelism bool
	}{
		{
			name: "Uncontended",
		},
		{
			name:            "Uncontended-HighParallelism",
			highParallelism: true,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			if bm.highParallelism {
				b.SetParallelism(highParallelismNum)
			}
			benchmarkSemaGroupUncontended(b)
		})
	}
}

func benchmarkWaitGroupAddDone(b *testing.B, localWork int) {
	var wg WaitGroup
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		foo := 0
		for pb.Next() {
			wg.Add(1)
			for i := 0; i < localWork; i++ {
				foo *= 2
				foo /= 2
			}
			wg.Done()
		}
		_ = foo
	})
}

func BenchmarkWaitGroupAddDone(b *testing.B) {
	benchmarks := []struct {
		name            string
		localWork       int
		highParallelism bool
	}{
		{
			name: "no work",
		},
		{
			name:      "with work",
			localWork: 100,
		},
		{
			name:            "with work-HighParallelism",
			localWork:       100,
			highParallelism: true,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			if bm.highParallelism {
				b.SetParallelism(highParallelismNum)
			}
			benchmarkWaitGroupAddDone(b, bm.localWork)
		})
	}
}

func benchmarkSemaGroupAddDone(b *testing.B, localWork int) {
	var sg sema.Group
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		foo := 0
		for pb.Next() {
			sg.Reserve()
			for i := 0; i < localWork; i++ {
				foo *= 2
				foo /= 2
			}
			sg.Free()
		}
		_ = foo
	})
}

func BenchmarkSemaGroupAddDone(b *testing.B) {
	benchmarks := []struct {
		name            string
		localWork       int
		highParallelism bool
	}{
		{
			name: "no work",
		},
		{
			name:      "with work",
			localWork: 100,
		},
		{
			name:            "with work-HighParallelism",
			localWork:       100,
			highParallelism: true,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			if bm.highParallelism {
				b.SetParallelism(highParallelismNum)
			}
			benchmarkSemaGroupAddDone(b, bm.localWork)
		})
	}
}

func benchmarkWaitGroupWait(b *testing.B, localWork int) {
	var wg WaitGroup
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		foo := 0
		for pb.Next() {
			wg.Wait()
			for i := 0; i < localWork; i++ {
				foo *= 2
				foo /= 2
			}
		}
		_ = foo
	})
}

func BenchmarkWaitGroupWait(b *testing.B) {
	benchmarks := []struct {
		name            string
		localWork       int
		highParallelism bool
	}{
		{
			name: "no work",
		},
		{
			name:      "with work",
			localWork: 100,
		},
		{
			name:            "with work-HighParallelism",
			localWork:       100,
			highParallelism: true,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			if bm.highParallelism {
				b.SetParallelism(highParallelismNum)
			}
			benchmarkWaitGroupWait(b, bm.localWork)
		})
	}
}

func benchmarkSemaGroupWait(b *testing.B, localWork int) {
	var sg sema.Group
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		foo := 0
		for pb.Next() {
			sg.Wait()
			for i := 0; i < localWork; i++ {
				foo *= 2
				foo /= 2
			}
		}
		_ = foo
	})
}

func BenchmarkSemaGroupWait(b *testing.B) {
	benchmarks := []struct {
		name            string
		localWork       int
		highParallelism bool
	}{
		{
			name: "no work",
		},
		{
			name:      "with work",
			localWork: 100,
		},
		{
			name:            "with work-HighParallelism",
			localWork:       100,
			highParallelism: true,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			if bm.highParallelism {
				b.SetParallelism(highParallelismNum)
			}
			benchmarkSemaGroupWait(b, bm.localWork)
		})
	}
}

func benchmarkWaitGroupActuallyWait(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var wg WaitGroup
			wg.Add(1)
			go func() {
				wg.Done()
			}()
			wg.Wait()
		}
	})
}

func BenchmarkWaitGroupActuallyWait(b *testing.B) {
	benchmarks := []struct {
		name            string
		highParallelism bool
	}{
		{
			name: "",
		},
		{
			name:            "HighParallelism",
			highParallelism: true,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			if bm.highParallelism {
				b.SetParallelism(highParallelismNum)
			}
			benchmarkWaitGroupActuallyWait(b)
		})
	}
}

func benchmarkSemaGroupActuallyWait(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var sg sema.Group
			sg.Reserve()
			go func() {
				sg.Free()
			}()
			sg.Wait()
		}
	})
}

func BenchmarkSemaGroupActuallyWait(b *testing.B) {
	benchmarks := []struct {
		name            string
		highParallelism bool
	}{
		{
			name: "",
		},
		{
			name:            "HighParallelism",
			highParallelism: true,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			if bm.highParallelism {
				b.SetParallelism(highParallelismNum)
			}
			benchmarkSemaGroupActuallyWait(b)
		})
	}
}
