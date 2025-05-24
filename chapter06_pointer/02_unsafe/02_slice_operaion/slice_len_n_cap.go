package main

import (
	"fmt"
	"unsafe"
)

func main() {
	// 切片操作
	sliceOperation()

}

func sliceOperation() {
	// 有一个内存分配相关的事实：结构体会被分配一块连续的内存，结构体的地址也代表了第一个成员的地址
	/* runtime/slice.go
	type slice struct{
		array unsafe.Pointer
		len int
		cap int

	}
	func makeslice() slice  返回slice 结构体
	*/
	s := make([]int, 9, 20)

	var len1 = *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + unsafe.Sizeof(int(0))))
	fmt.Println("长度", len1, len(s))

	var cap1 = *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + uintptr(16)))
	fmt.Println("容量", cap1, cap(s))
	// 转换过程 Len: &s => pointer => uintptr => pointer => *int => int
}
