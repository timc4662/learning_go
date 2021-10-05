package main

import (
	"fmt"
	"sync"
	"time"
)

var once sync.Once

// demos how to use a WaitGroup to wait for 3 go routines to finish
func main() {
	var wg sync.WaitGroup

	// 3 as we will create three go routines
	wg.Add(3)

	for i := 0; i < 3; i++ {

		once.Do(func() {
			fmt.Println("starting the first routine")
		})

		i := i
		go func() {
			time.Sleep(time.Second * time.Duration(i))
			fmt.Printf("%d done\n", i)
			wg.Done()
		}()
	}

	// wait for all done
	wg.Wait()

	fmt.Println("all done")
}
