package main

import "fmt"

func bucketShift(b uint8) uintptr {
	return uintptr(1) << b
}

func overLoadFactor(count int, B uint8) bool {
	return count > 8 && uintptr(count) > 13*(bucketShift(B)/2)
}

func bucketLen(hint int) uint8 {
	B := uint8(0)
	for overLoadFactor(hint, B) {
		B++
	}

	return B
}

func main() {
	for i := 0; i < 100; i++ {
		fmt.Println(i, bucketLen(i))
	}
}
