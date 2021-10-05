package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type processor struct {
	outA chan string
	outB chan string
	outC chan string
	inC  chan string
	errs chan error
}

func simulateRemoteRequest(ctx context.Context, input string) (string, error) {
	fmt.Println("start sim request: " + input)
	defer fmt.Println("end sim request: " + input)
	for i := 0; i < 10; i++ {
		if ctx.Err() != nil {
			fmt.Println("cancelled: " + input)
			return "", ctx.Err()
		}
		time.Sleep(time.Millisecond)
	}
	return "hello: " + input, nil
}

func makeBgRequest(ctx context.Context, input string, ch chan<- string, errs chan<- error) {
	go func() {
		defer fmt.Println("exit go routine")
		r, err := simulateRemoteRequest(ctx, input)
		// if an error, add this to the errors channel
		if err != nil {
			errs <- err
			return
		}
		ch <- r
	}()
}

func (p *processor) launch(ctx context.Context, data string) {
	// make the calls to the three remote services (simulated)
	makeBgRequest(ctx, "input a: "+data, p.outA, p.errs)
	makeBgRequest(ctx, "input b: "+data, p.outB, p.errs)
	makeBgRequest(ctx, "input c: "+data, p.outC, p.errs)
}

func (p *processor) wait(ctx context.Context) ([]string, error) {
	var results []string
	// gather the results
	for len(results) < 3 {
		select {
		case a := <-p.outA:
			results = append(results, a)
		case b := <-p.outB:
			results = append(results, b)
		case c := <-p.outC:
			results = append(results, c)
		case err := <-p.errs:
			return []string{}, err
		case <-ctx.Done():
			return nil, errors.New("timeout c")
		}
	}
	return results, nil
}

func makeRequests(ctx context.Context) ([]string, error) {
	// create a derived ctx which times out after 50ms
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	p := processor{
		outA: make(chan string, 1),
		outB: make(chan string, 1),
		outC: make(chan string, 1),
		inC:  make(chan string, 1),
		errs: make(chan error, 3),
	}

	p.launch(ctx, "some data")

	// wait for completion
	return p.wait(ctx)
}

func main() {
	var ctx context.Context
	ctx = context.Background()
	results, err := makeRequests(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(results)
}
