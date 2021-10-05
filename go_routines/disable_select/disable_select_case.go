package main

import "fmt"

func generate() {
	ch := make(chan int)
	quit := make(chan struct{})
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
		quit <- struct{}{}
	}()
	var values []int
	for {
		select {
		case i, ok := <-ch:
			fmt.Println(ok)
			// When you need to disable a select, set the channel
			// variable to nil and it will never be hit again. Without this
			// in here, you get random amounts of 0's appended to the slice
			if ok == false {
				ch = nil
			}
			values = append(values, i)
			fmt.Println(values)
		case <-quit:
			goto exit
		}
	}
exit:
	close(quit)
}

func main() {
	generate()
}
