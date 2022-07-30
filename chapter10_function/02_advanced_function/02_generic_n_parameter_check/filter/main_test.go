package filter

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	var intSet = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// 1. 原地替换
	FilterInPlace(&intSet, func(n int) bool {
		return n%2 == 1
	})
	fmt.Println(intSet)

	// 2. 非原地替换
	newSlice := Filter(intSet, func(n int) bool {
		return n%3 == 0
	})
	fmt.Println(newSlice)

}
