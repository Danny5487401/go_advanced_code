package main

import "fmt"

func main() {

	// 1 如果要清空一个slice，那么可以简单的赋值为nil，垃圾回收器会自动回收原有的数据
	nilSlice := []int{1, 2, 3}
	fmt.Printf("nilSlice %p\n", &nilSlice)
	nilSlice = nil
	fmt.Printf("nilSlice=nil: %v,len: %d, cap: %d, elems: %v \n", nilSlice == nil, len(nilSlice), cap(nilSlice), nilSlice) // nilSlice=nil: true ,len: 0, cap: 0, elems: []

	// 2 还需要使用 slice 底层内存，那么最佳的方式是 re-slice [0:0]
	emptySlice := []int{1, 2, 3}
	emptySlice = emptySlice[:0]
	fmt.Printf("emptySlice %p\n", &emptySlice)
	fmt.Printf("len: %d, cap: %d, elems: %v \n", len(emptySlice), cap(emptySlice), emptySlice) // len: 0, cap: 3, elems: []
	fmt.Println(emptySlice[:1])                                                                // [1]

}
