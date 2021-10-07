package panic

import "fmt"

func AFunctionWhichPanics() {
	fmt.Println("about to panic")
	panic("something very bad went wrong")
}

// Can we return anything from this function even with a panic
func DoSomeWork() (result string) {

	defer func() {
		r := recover()
		result = "that was close.. the function panicked with " + r.(string)
	}()

	AFunctionWhichPanics()

	return "no chance to make it here"
}
