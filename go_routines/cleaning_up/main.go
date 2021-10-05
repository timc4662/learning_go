package main

import (
	"fmt"
	"time"
)

// Create a generator function which generates ints on the channel
// using a go routine up to max.
func timer(max int) <-chan int {
	ch := make(chan int)
	go func() {
		for i := 0; i < max; i++ {
			ch <- i * 100
			time.Sleep(time.Millisecond * 100)
		}
		// without the close, we get a deadlock.
		// always close on the writer side, never the reader side.
		fmt.Println("go routing is about to close")
		close(ch)
	}()
	return ch
}

func main() {
	// note: for loops over channels do not have indexes (obv)
	started := time.Now()
	for v := range timer(10) {
		fmt.Printf("%v ms\n", v)
		// if we exit the go routing here, then it never closes
		// as it blocks waiting for something to read. Fix this with the
		// done channel
		if false && v == 500 {
			break
		}
	}
	duration := time.Since(started)
	fmt.Println(duration)
}
