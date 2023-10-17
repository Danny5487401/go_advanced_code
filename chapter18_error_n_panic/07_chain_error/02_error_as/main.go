package main

import (
	"errors"
	"fmt"
)

type MyError struct {
	err string
}

func (e *MyError) Error() string {
	return e.err
}

func main() {
	err1 := &MyError{"temp error"}
	err2 := fmt.Errorf("2nd err: %w", err1)
	err3 := fmt.Errorf("3rd err: %w", err2)

	fmt.Println(err3)

	var e *MyError
	ok := errors.As(err3, &e)
	if ok {
		fmt.Println(e)
		return
	}
}
