package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

// This example shows how you may use a context
// rather than a done channel to cancel an operation.
func someWorker(ctx context.Context, i int) (int, error) {
	deadline, ok := ctx.Deadline()
	if ok {
		if deadline.Sub(time.Now().Add(time.Second*2)) < 0 {
			fmt.Printf("would have exceeded deadline of: %q\n", deadline)
			// it hasn't enough time to compleet the next wait, cancel immediately.
			return 0, context.DeadlineExceeded
		}
	}
	select {
	case <-time.After(time.Second * 2):
		fmt.Print(".")
		return i, nil
	case <-ctx.Done():
		fmt.Println("Context is done")
		return 0, ctx.Err()
	}
}

func main() {

	a := func(ctx context.Context) <-chan struct{} {
		finished := make(chan struct{})
		go func() {
			defer close(finished)
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			var i int
			fmt.Println("Hit ctrl+c to quit, or it will quit in 5 seconds")
			for {
				var err error
				i, err = someWorker(ctx, i)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}()
		return finished
	}

	ctx := context.Background()
	// We can really easily hook the signal interrupt up
	// to the context which will cancel this on Ctrl+C
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	finished := a(ctx)
	<-finished
}
