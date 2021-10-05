package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Counter int64

func (c *Counter) Inc() {
	(*c)++
}

func (c *Counter) Dec() {
	(*c)--
}

func (c *Counter) Print() string {
	return strconv.FormatInt(int64(count), 10)
}

var count Counter

// we only want this to be run 5 times, we can build this with go routines.
func limitedJob() string {
	count.Inc()
	fmt.Println("limited job running: " + count.Print())
	defer fmt.Println("limited job finished: " + count.Print())
	defer count.Dec()
	time.Sleep(2 * time.Second)
	return "done"
}

type PressureGauge struct {
	ch chan struct{}
}

func (pg *PressureGauge) Process(f func()) error {
	select {
	case <-pg.ch:
		// if we can read from the channel, then it is acceptable to run the function
		f()
		// pop back.
		pg.ch <- struct{}{}
		return nil
	default:
		// if we cannot read the channel, the default clause will fire
		// so we must have exceeded the limit and return an error
		return errors.New("capacity reached")
	}
}

// creates a new pressure gauge with the limit
// fills the channel with this number of empty structs
// so the process func can at most read this many to start
// this whole thing acts much like a semaphore
func NewPG(limit int) *PressureGauge {
	ch := make(chan struct{}, limit)
	for i := 0; i < limit; i++ {
		ch <- struct{}{}
	}
	return &PressureGauge{ch}
}

// run this and open localhost:8080/request and hit refresh
// it only allows 1 request to be processed at once.
func main() {
	pg := NewPG(1)
	http.HandleFunc("/request", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("starting..."))
		err := pg.Process(func() {
			w.Write([]byte(limitedJob()))
		})
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	})
	http.ListenAndServe(":8080", nil)
}
