package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"time"
)

type CustomError struct {
	Inner      error
	Message    string
	StackTrace string
}

func (err CustomError) Error() string {
	return err.Message
}

func wrapError(err error, message string) error {
	return CustomError{
		Message:    message,
		Inner:      err,
		StackTrace: string(debug.Stack()),
	}
}

// This example shows how you may use a context
// rather than a done channel to cancel an operation.
func someWorker(ctx context.Context, i int) (int, error) {
	deadline, ok := ctx.Deadline()
	if ok {
		if deadline.Sub(time.Now().Add(time.Second*2)) < 0 {
			fmt.Printf("would have exceeded deadline of: %q\n", deadline)
			// it hasn't enough time to compleet the next wait, cancel immediately.
			return 0, wrapError(context.DeadlineExceeded, "Would have exceeded timeout")
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
				// example of how to deal with custom errors
				if err != nil {
					// Note this does not fire as this is not a custom error
					// is there a way to implement the Is?  Out of scope of this test
					if errors.Is(err, CustomError{}) {
						log.Println("This won't fire with the CustomError!")
					}
					// But the comma ok idiom can be used to determine the error type
					if err, ok := err.(CustomError); ok {
						log.Println(err.StackTrace)
					}
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
