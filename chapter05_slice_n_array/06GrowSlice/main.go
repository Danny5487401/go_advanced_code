package main

import (
	"fmt"
)

func main() {
	// 看扩容策略
	slice := []int{10, 20, 30, 40}
	newSlice := append(slice, 50)
	// Before slice = [10 20 30 40], Pointer = 0xc000098420, len = 4, cap = 4
	fmt.Printf("Before slice = %v, Pointer = %p, len = %d, cap = %d\n", slice, &slice, len(slice), cap(slice))

	// Before newSlice = [10 20 30 40 50], Pointer = 0xc000098440, len = 5, cap = 8
	fmt.Printf("Before newSlice = %v, Pointer = %p, len = %d, cap = %d\n", newSlice, &newSlice, len(newSlice), cap(newSlice))

	newSlice[1] += 10

	// After slice = [10 20 30 40], Pointer = 0xc000098420, len = 4, cap = 4
	fmt.Printf("After slice = %v, Pointer = %p, len = %d, cap = %d\n", slice, &slice, len(slice), cap(slice))

	// After newSlice = [10 30 30 40 50], Pointer = 0xc000098440, len = 5, cap = 8
	fmt.Printf("After newSlice = %v, Pointer = %p, len = %d, cap = %d\n", newSlice, &newSlice, len(newSlice), cap(newSlice))

}
