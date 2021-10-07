package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/timc4662/learning_go/errors/panic"
	"github.com/timc4662/learning_go/errors/sentinel"
)

// Normal error handling

func doSomething() error {
	return errors.New("this never worked")
}

// by convention the second argument is err and type error
func somethingElse(msg string) (ret string, err error) {
	if msg == "valid string" {
		ret = "great stuff"
		return
	}
	if msg == "invalid string" {
		err = errors.New("the string cannot be the value 'invalid string'")
		return
	}
	return
}

// NormalErrorHandlingDemo demos std error handling in go.
func NormalErrorHandlingDemo() {
	// The happy path should be on the left hand side
	err := doSomething()
	if err != nil {
		fmt.Println(err)
	}
	ret, err2 := somethingElse("valid string")
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println(ret)
	_, err3 := somethingElse("invalid string")
	if err3 != nil {
		fmt.Println(err3)
		return
	}
}

// Custom errors.

// Anything is an error if it implements the Error() interface

type MyError struct {
	lineNo int
	msg    string
}

// Now make it implement error
func (err MyError) Error() string {
	return "This is the actual error. Line number: " + strconv.Itoa(err.lineNo)
}

func doWork() (string, error) {
	return "", MyError{
		lineNo: 10,
		msg:    "this would have failed on line 10",
	}
}

func CustomErrorDemo() {
	ret, err := doWork()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret)
}

func main() {
	NormalErrorHandlingDemo()
	CustomErrorDemo()
	if sentinel.DetectVideoLoss() {
		fmt.Println("video loss")
	}
	fmt.Println(panic.DoSomeWork())
}
