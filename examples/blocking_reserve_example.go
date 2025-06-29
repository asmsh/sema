package main

import (
	"context"
	"log"
	"time"

	"github.com/asmsh/sema"
)

func main() {
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
