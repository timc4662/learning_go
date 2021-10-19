package main

import (
	"fmt"
	"sync"
	"time"
)

// run this with go run -race to detect data races

// No protection, this has a data race
func exp1NoProtection() {
	var i int
	go func() {
		i = i + 1
	}()
	fmt.Println(i)
}

// critical sections protected with a mutex
// note race cond still present
func exp2MutexSync() {
	var i int
	var m sync.Mutex
	go func() {
		// calls to lock can make the program slow
		m.Lock()
		i = i + 1
		m.Unlock()
	}()
	m.Lock()
	fmt.Println(i)
	m.Unlock()
}

func exp3Deadlock() {
	type value struct {
		mu    sync.Mutex
		value int
	}

	var wg sync.WaitGroup
	printSum := func(v1, v2 *value) {
		defer wg.Done()
		v1.mu.Lock()
		print("a\n")
		defer v1.mu.Unlock()

		time.Sleep(2 * time.Second)
		print("b\n")
		// deadlock here as gr 1 is waiting for b.v1 to be released
		// and gr 2 is waiting for a.v1 to be released before releasing b.v1

		// Interestingly if run with the -race flag, we don't see the goroutines asleep warning
		// it hangs forever, but without the flag we do.

		v2.mu.Lock()
		defer v2.mu.Unlock()
		print("c\n")
		fmt.Printf("sum=%v\n", v1.value+v2.value)
	}

	var a, b value
	wg.Add(2)
	go printSum(&a, &b)
	go printSum(&b, &a)
	wg.Wait()
}

// live locks.
// starvation
func exp4Starvation() {

	var wg sync.WaitGroup
	var sharedLock sync.Mutex
	const runtime = 1 * time.Second

	greedyWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(3 * time.Nanosecond)
			sharedLock.Unlock()
			count++
		}

		fmt.Printf("Greedy worker was able to execute %v work loops\n", count)
	}

	politeWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			count++
		}

		fmt.Printf("Polite worker was able to execute %v work loops.\n", count)
	}

	wg.Add(2)

	go greedyWorker()
	go politeWorker()

	wg.Wait()
}

func exp5() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		//salutation := salutation
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
		}()
		// this is interesting in that this works. although horrible.
		// the problem with the above is that the go routines are not started
		// immediately, so whilst they capture salutation, it is updated to good day
		// before they are started, thus they all print the same.
		// the shadowing trick just lets them all capture different variables and not the same
		// one.
		time.Sleep(time.Millisecond * 50)
	}
	wg.Wait()
}

func exp6() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Removed from queue")
		c.L.Unlock()
		c.Signal()
	}

	test := func() bool {
		print("testing length\n")
		return len(queue) == 2
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()
		for test() {
			c.Wait()
			print("something happened - Cond signalled\n")
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()
	}
}

func expBroadcastOnCond() {
	var m sync.Mutex
	cond := sync.NewCond(&m)
	var wg, wg2 sync.WaitGroup
	wg.Add(200)
	wg2.Add(200)
	var done bool = false

	for i := 0; i < 200; i++ {
		go func() {
			defer wg2.Done()

			print("waiting on bc...\n")
			wg.Done()

			m.Lock()
			if !done {
				cond.Wait()
			}

			fmt.Println("broadcast occurred...")
			m.Unlock()

		}()
	}

	wg.Wait()
	print("started, broadcasting...\n")
	m.Lock()
	done = true
	cond.Broadcast()
	m.Unlock()

	wg2.Wait()
}

func expPool() {
	var c int

	p := sync.Pool{
		New: func() interface{} {
			newCount := c
			c++
			fmt.Printf("in here: %d\n", newCount)
			return c
		},
	}
	inst := p.Get()
	inst2 := p.Get()
	fmt.Printf("equals = %v %v\n", inst.(int), inst2.(int))
	p.Put(inst)
	inst3 := p.Get()
	fmt.Printf("equals = %v %v\n", inst.(int), inst3.(int))
}

type MyTester struct {
}

func (mt *MyTester) test() {
	print("here\n")
}

func messing() {

	type Tester interface {
		test()
	}

	type myfunc func()

	x := struct {
		myfunc
	}{
		myfunc: func() {
			print("here 2\n")
		},
	}
	var t Tester = &MyTester{}
	t.test()
	x.myfunc()

}

func exp7() {
	doWork := func(
		done <-chan interface{},
		strings <-chan string,
	) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					// Do something interesting
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done, nil)

	go func() {
		// Cancel the operation after 1 second.
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated
	fmt.Println("Done.")
}
func main() {
	//exp1NoProtection()
	//exp2MutexSync()
	//exp3Deadlock()
	//exp4Starvation()
	//exp5()
	//exp6()
	//expBroadcastOnCond()
	//expPool()
	//messing()

}
