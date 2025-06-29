# Sema: A feature-rich concurrency manager for Go.

[![PkgGoDev](https://pkg.go.dev/badge/github.com/asmsh/sema)](https://pkg.go.dev/github.com/asmsh/sema)
[![Go Report Card](https://goreportcard.com/badge/github.com/asmsh/sema)](https://goreportcard.com/report/github.com/asmsh/sema)
[![Tests](https://github.com/asmsh/sema/workflows/Tests/badge.svg)](https://github.com/asmsh/sema/actions)
[![Go Coverage](https://github.com/asmsh/sema/wiki/coverage.svg)](https://raw.githack.com/wiki/asmsh/sema/coverage.html)

It bridges the gap between using channels, [sync.WaitGroup](https://pkg.go.dev/sync#WaitGroup) and [semaphore.Weighted](https://pkg.go.dev/golang.org/x/sync/semaphore#Weighted) for
managing concurrency.

### Features

* Fast, low-overhead implementation with performance comparable to `sync.WaitGroup` and `semaphore.Weighted`.
* Enables `select`-based waiting, in addition to standard `sync.WaitGroup` wait behavior.
* Exposes the internal counters in a concurrent-safe way.
* Clear API that's compatible with `sync.WaitGroup` and `semaphore.Weighted` with minimal changes.
* Usable zero value that's comparable to a `sync.WaitGroup`.

### Notes

* It wakes up blocked calls in random order, unlike `semaphore.Weighted` which preserves the order of calls.
* It's more suitable as a replacement for `semaphore.Weighted` when all the weights are of equal size.
* It relies on a single channel for blocking and wake-ups, so the Go runtime guarantees no starvation.

### Examples

#### Using it as a `sync.WaitGroup`:

```go
func example() {
	var sg sema.Group

	for i := 0; i < 10; i++ {
		sg.Reserve()
		go func(i int) {
			defer sg.Free()

			// Do some work...
		}(i)
	}

	sg.Wait()
}
```

#### Using `select`-based waiting:

```go
func example() {
	var sg sema.Group

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		sg.Reserve()
		go func(i int) {
			defer sg.Free()

			// Do some work using the ctx...
			ctx = ctx
		}(i)
	}

	waitChan := sg.WaitChan()

	select {
	case <-waitChan:
		// all the work is done.
	case <-ctx.Done():
		// the work times out.
	}
}
```

#### Using it to limit concurrency (instead of a `semaphore.Weighted` or a channel):

```go
func example() {
	var sg sema.Group
	sg.SetSize(10)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		sg.ReserveN(ctx.Done(), 5)
		go func(i int) {
			defer sg.FreeN(5)

			// Do some work...
			log.Println("some work...")
		}(i)
	}

	sg.Wait()
}
```