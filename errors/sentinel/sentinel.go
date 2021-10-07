package sentinel

import (
	"errors"
	"io"
)

// sentinel errors are declared at package level
// by convention their names start with Err (io.EOF:))
// treated as read only, but no compiler enforcement.

// A possible alternative would be constant errors https://dave.cheney.net/2016/04/07/constant-errors
// however this has disadvantages too in that they are strings, so the error would be equal to a string
// with the same value, or if two packages created errors with the same string value a.Err and b.Err would
// be equal.

// Better just to avoid sentinel errors if possible.

// ErrVideoLoss simulates an error generated on video loss.
var ErrVideoLoss = errors.New("custom error 1")

func Generate() error {
	// randomly legal, but very bad. See below for constant errors idiom
	// type error is not constant.
	io.EOF = errors.New("whoops")
	return ErrVideoLoss
}

// DetectVideoLoss is simply a test of comparing against sentinel errors
func DetectVideoLoss() bool {
	err := Generate()
	return err == ErrVideoLoss
}

// constant errors idiom
// anything which implements the Error interface is an error, even a string
type MyError string

func (e MyError) Error() string {
	return string(e)
}

const ErrVideoLoss2 = MyError("custom error 2")

// Cannot assign to it
// ErrVideoLoss2 = MyError("xyz");
func DetectVideoLoss2() (bool, error) {
	err := ErrVideoLoss2
	if err == ErrVideoLoss2 {
		return true, nil
	}
	return false, err
}
