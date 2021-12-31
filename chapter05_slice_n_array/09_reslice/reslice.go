package main

import "fmt"

func main() {
	a := []int{0, 1, 2, 3, 4, 5, 6}
	fmt.Println(lastNumsBySlice(a))
	fmt.Println(lastNumsByCopy(a))

}

// 原始切片上操作，底层数组没有发生变化，内存一直占用，直到没有变量引用该数组，这种操作不推荐
func lastNumsBySlice(origin []int) []int {
	return origin[len(origin)-2:]
}

// 推荐做法：copy
func lastNumsByCopy(origin []int) []int {
	result := make([]int, 2)
	copy(result, origin[len(origin)-2:])
	return result

}
