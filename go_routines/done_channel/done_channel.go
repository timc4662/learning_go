package main

import "fmt"

type searchFunc func(string) []string

func searchData(s string, searchers []searchFunc) []string {
	done := make(chan struct{})
	result := make(chan []string)
	for _, searcher := range searchers {
		go func(searcher searchFunc) {
			defer fmt.Println("go routine is now closing")
			found := searcher(s)
			select {
			case result <- found:
			case <-done:
			}
		}(searcher)
	}
	r := <-result
	// This is the key part of this sample. Without this close
	// the second go routine never exits as we read one value from result channel
	// (line above), so the second one is blocked waiting to write
	// and never close. When we close, we get a read on done, which frees off the go routine.
	close(done)
	return r
}

func s1(term string) []string {
	return []string{"abc"}
}

func s2(term string) []string {
	return []string{"abc"}
}

func main() {
	sf1 := searchFunc(s1)
	sf2 := searchFunc(s2)
	result := searchData("abc", []searchFunc{sf1, sf2})
	fmt.Println(result)
}
