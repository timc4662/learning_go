package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// Some example 'generic' generators. A generator can be used in a pipeline.
// A generator is a function whch takes a set of values and pushes them through a channel.

func main() {

	repeat := func(
		done <-chan interface{},
		values ...interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			defer fmt.Println("Closed repeat")
			// this generator takes the values and
			// writes them to the valueStream forever until cancelled
			// with a write on the done channel.
			for {
				for _, v := range values {
					fmt.Printf("writing %v\n", v)
					select {
					case <-done:
						fmt.Println("Got done on repeat")
						return

					case valueStream <- v:
					}
				}
			}

		}()
		return valueStream
	}

	take := func(
		done <-chan interface{},
		valueStream <-chan interface{},
		num int,
		name string,
	) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			defer fmt.Println("Closed take")
			// this generator will take num values from the valueStream
			// and write them to the takeStream channel. Cancelled when done.
			for i := 0; i < num; i++ {
				select {
				case <-done:
					fmt.Println("got done on take")
					return
				case v := <-valueStream:
					fmt.Printf("take %v by %v\n", v, name)
					if v != nil {
						takeStream <- v
					}
				}
			}
		}()
		return takeStream
	}

	// this pipeline stage will convert values into strings
	// using the conv function. The output of this function will therefore
	// be strings.
	toString := func(
		done <-chan interface{},
		conv func(interface{}) string,
		valueStream <-chan interface{},
	) <-chan string {
		stringStream := make(chan string)
		go func() {
			defer close(stringStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case stringStream <- conv(v):
				}
			}
		}()
		return stringStream
	}

	done := make(chan interface{})

	v := []int{2, 3, 4, 5, 3, 1000}

	c := func(slice []int) []interface{} {
		r := make([]interface{}, len(slice))
		for i := 0; i < len(slice); i++ {
			r[i] = interface{}(slice[i])
		}
		return r
	}

	iiToStr := func(v interface{}) string {
		return strconv.FormatInt(int64(v.(int)), 10)
	}

	// fanIn takes a varadic slice of channels and returns a single channel
	fanIn := func(
		done <-chan interface{},
		channels ...<-chan interface{},
	) <-chan interface{} {
		var wg sync.WaitGroup
		multiplexedStream := make(chan interface{})

		multiplex := func(c <-chan interface{}) {
			defer wg.Done()
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		// Select from all the channels
		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		// Wait for all the reads to complete
		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}

	valueStream := repeat(done, c(v)...)

	// here is an example fanning out. We use two take generators
	// to read from the value stream. In real life this would perhaps
	// be expensive work.
	var takers []<-chan interface{}
	takers = append(takers, take(done, valueStream, 5, "t1"))
	takers = append(takers, take(done, valueStream, 5, "t2"))

	// And here we use the fanIn algorithm to combine the channels back together again
	for num := range toString(done, iiToStr, fanIn(done, takers...)) {
		fmt.Printf("Reading %q\n ", num)

		if num == "1000" {
			fmt.Println("Closing...")
			close(done)
			break
		}

	}

	time.Sleep(time.Millisecond * 50)
}
