package main

import (
	"fmt"
	"time"
)

// simulate doing some work
func doWork(i int) int {
	time.Sleep(time.Millisecond * 500)
	return i * 2
}

// generates go routines for each input value
// returns a readable channel for the results
func generate(input []int) <-chan int {
	results := make(chan int, len(input))
	for i := 0; i < len(input); i++ {
		go func(i int) {
			results <- doWork(input[i])
		}(i)
	}
	return results
}

// takes a readable channel and returns the first len ints from it
func readResults(ch <-chan int, len int) []int {
	var out []int
	for i := 0; i < len; i++ {
		out = append(out, <-ch)
	}
	return out
}

func main() {
	// input data to be processed by doWork
	input := []int{2, 4, 6, 8}
	// process each input on a separate go routing
	// results are added to a buffered channel with the buffer
	// size the same size as the input dataset
	ch := generate(input)
	// We then read all the values back from the channel
	result := readResults(ch, len(input))
	fmt.Println(result)
}
