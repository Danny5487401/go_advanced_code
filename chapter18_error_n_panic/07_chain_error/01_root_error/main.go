package main

import (
	"errors"
	"fmt"
)

func rootCause(err error) error {
	for {
		e, ok := err.(interface{ Unwrap() error })
		if !ok {
			return err
		}
		err = e.Unwrap()
		if err == nil {
			return nil
		}
	}
}

func main() {
	err1 := errors.New("error1")

	err2 := fmt.Errorf("2nd err: %w", err1)
	err3 := fmt.Errorf("3rd err: %w", err2)

	fmt.Println(err3) // 3rd err: 2nd err: error1

	fmt.Println(rootCause(err1)) // error1
	fmt.Println(rootCause(err2)) // error1
	fmt.Println(rootCause(err3)) // error1
}
