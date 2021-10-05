package main

import (
	"fmt"
	"time"
)

func countTo(max int) (<-chan int, func()) {
	ch := make(chan int)
	done := make(chan struct{})
	// return a closure rather than the done channel directly.
	cancel := func() {
		fmt.Println("closing done")
		close(done)
	}
	go func() {
		for i := 0; i < max; i++ {
			select {
			case <-done:
				fmt.Println("<- done")
				goto exit
			default:
				fmt.Printf("ch <- %d\n", i)
				ch <- i
			}
		}
	exit:
		fmt.Println("closing ch")
		close(ch)
	}()
	return ch, cancel
}

// this sample shows how you can use a done channel wrapped in a closure
// to allow the caller to exit a go routine. It generates 10 numbers, but
// after 5 we quit the loop, so it would block writing the 6th. But we
// call cancel which triggers in the select, exiting the go routine.
func main() {
	ch, cancel := countTo(10)
	for v := range ch {
		if v == 5 {
			break
		}
		fmt.Println(v)
	}
	cancel()
	// slight sleep so we can see <- done
	time.Sleep(time.Millisecond * 50)
}
