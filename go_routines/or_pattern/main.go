package main

import (
	"fmt"
	"time"
)

// This is an example of the Or patten from the Go concurrency book
// The or function takes a varadic slice of channels and when one
// of them completes, it closes the orDone channel which is returned
// a client can monitor the orDone to detect when any of the channels is complete.

func main() {
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() {
			defer close(orDone)

			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()
		return orDone
	}

	sig := func(done <-chan struct{}, after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			select {
			case <-done:
				fmt.Printf("cancelled %v\n", after)
				return
			case <-time.After(after):
				break
			}
			fmt.Printf("completed %v\n", after)
		}()
		return c
	}

	start := time.Now()

	done := make(chan struct{})

	<-or(
		sig(done, time.Millisecond*20),
		sig(done, time.Millisecond*250),
		sig(done, time.Millisecond*30),
		sig(done, time.Millisecond*10),
	)

	fmt.Printf("done after %v\n", time.Since(start))

	// Close the other in progress channels.
	close(done)
	time.Sleep(time.Millisecond * 10)
}
