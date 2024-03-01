package main

import "fmt"

func main() {

	// 1 如果要清空一个slice，那么可以简单的赋值为nil，垃圾回收器会自动回收原有的数据
	nilSlice := []int{1, 2, 3}
	fmt.Printf("nilSlice %p\n", &nilSlice)
	nilSlice = nil
	fmt.Println(nilSlice, len(nilSlice), cap(nilSlice), nilSlice == nil) // [] 0 0 true

	// 2 还需要使用 slice 底层内存，那么最佳的方式是 re-slice [0:0]
	emptySlice := []int{1, 2, 3}
	emptySlice = emptySlice[:0]
	fmt.Printf("emptySlice %p\n", &emptySlice)
	fmt.Println(emptySlice, len(emptySlice), cap(emptySlice)) // [] 0 3
	fmt.Println(emptySlice[:1])                               // [1]
}
