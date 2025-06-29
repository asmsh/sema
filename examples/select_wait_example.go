package main

import (
	"context"
	"log"
	"time"

	"github.com/asmsh/sema"
)

func main() {
	var sg sema.Group

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		sg.Reserve()
		go func(i int) {
			defer sg.Free()

			// Do some work using the ctx...
			ctx = ctx
			log.Println("some work...")
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
