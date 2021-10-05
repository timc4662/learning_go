package main

import (
	"fmt"
)

// This sample shows a common bug with capture variables and go routines
// but it is detected by go vet so no worries. The solution is to use shadowing
// or alternatively to use an argument.
func main() {

	a := []int{2, 4, 6, 8, 10}
	ch := make(chan int, len(a))
	for _, v := range a {
		// v := v             // shadowing (unused)
		go func(val int) { // or pass as argument
			ch <- val * 2
		}(v)
	}
	// this should the values of a *2 (in a random order)
	for i := 0; i < len(a); i++ {
		fmt.Println(<-ch)
	}
}
