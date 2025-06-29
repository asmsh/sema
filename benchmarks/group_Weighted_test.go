// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Note: this file is copied from: golang.org/x/sync/semaphore.semaphore_test.go file.
// It's been modified to run the same tests that are done on the [semaphore.Weighted] type,
// but on the [sema.Group] type instead.
// Note: this is the only remaining test from the mentioned file that's not
// in the base package.
// It's move here to avoid adding a dependency on the golang.org/x/sync in
// the base package.

package benchmarks

import (
	"context"
	"testing"
	"time"

	"github.com/asmsh/sema"
	"golang.org/x/sync/errgroup"
)

func TestGroupDoesntBlockIfTooBig(t *testing.T) {
	t.Parallel()

	const n = 2
	sg := sema.NewGroup(n)
	{
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go sg.ReserveN(ctx.Done(), n+1)
	}

	g, ctx := errgroup.WithContext(context.Background())
	for i := n * 3; i > 0; i-- {
		g.Go(func() error {
			if sg.ReserveN(ctx.Done(), 1) {
				time.Sleep(1 * time.Millisecond)
				sg.Free()
				return nil
			} else {
				return ctx.Err()
			}
		})
	}
	if err := g.Wait(); err != nil {
		t.Errorf("semaphore.NewWeighted(%v) failed to AcquireCtx(_, 1) with AcquireCtx(_, %v) pending", n, n+1)
	}
}
