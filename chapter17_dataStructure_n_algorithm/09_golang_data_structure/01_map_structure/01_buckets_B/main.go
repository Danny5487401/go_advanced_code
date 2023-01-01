package main

import "fmt"

const (
	loadFactorNum = 13
	loadFactorDen = 2
	bucketCntBits = 3
	bucketCnt     = 1 << bucketCntBits
)

func bucketShift(b uint8) uintptr {
	return uintptr(1) << b
}

func overLoadFactor(count int, B uint8) bool {
	return count > bucketCnt && uintptr(count) > loadFactorNum*(bucketShift(B)/loadFactorDen)
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
