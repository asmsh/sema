package main

import (
	"log"

	"github.com/asmsh/sema"
)

func main() {
	var sg sema.Group

	for i := 0; i < 10; i++ {
		sg.Reserve()
		go func(i int) {
			defer sg.Free()

			// Do some work...
			log.Println("some work...")
		}(i)
	}

	sg.Wait()
}
