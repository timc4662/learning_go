// a heartbeat is a way for an idle go routine to
// signal that it is actually ok state
// Useful in unit tests as makes them determistic
// For any long running routines or where testing is needed

package main

import (
	"fmt"
	"time"
)

func main() {

	// If set to true we fakle an error in the go routine.
	simError := false

	doWork := func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) (<-chan interface{}, <-chan time.Time) {
		heartbeat := make(chan interface{})
		results := make(chan time.Time)
		go func() {
			// simulate but by not closing
			if !simError {
				defer close(heartbeat)
				defer close(results)
			}

			pulse := time.Tick(pulseInterval)
			workGen := time.Tick(2 * pulseInterval)

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default: // guard against nothing listening to heartbeats
				}
			}

			var count int
			for {
				select {
				case <-done:
					return
				case <-pulse: // send the heartbeat
					sendPulse()
				case r := <-workGen:
					results <- r
				}
				// sim error by exiting the loop and go routine without
				// closing
				if simError {
					count++
					if count == 2 {
						break
					}
				}
			}
		}()
		return heartbeat, results
	}

	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			// ok indicates whether it has been closed
			if ok == false {
				return
			}
			fmt.Println("heartbeat")
		case r, ok := <-results:
			// ok indicates whether it has been closed
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r.Second())
			// each time round the loop this is reset
			// so as we get a heartbeat normally within timeout/2
			// this should NEVER be hit, unless something is wrong the
			// the goroutine.
		case <-time.After(timeout):
			fmt.Println("worker goroutine is not healthy!")
			return
		}
	}
}
