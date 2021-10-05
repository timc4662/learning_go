package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

func someExpensiveWorker() (int, error) {
	time.Sleep(time.Second)
	return 4, nil
}

func processWithLimit() (int, error) {
	var result int
	var err error
	done := make(chan struct{})
	go func() {
		result, err = someExpensiveWorker()
		close(done)
	}()
	select {
	case <-done:
		return result, err
	case <-time.After(500 * time.Millisecond):
		return 0, errors.New("timeout")
	}
}

// this code will always exit with fatal timeout
// because we use the time.After() with 500ms. this returns a channel
// which writes after this period, so this will always have data before
// the doWork completes
func main() {
	r, err := processWithLimit()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r)
}
