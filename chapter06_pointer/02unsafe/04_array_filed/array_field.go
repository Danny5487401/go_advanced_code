package main

import (
	"fmt"
	"unsafe"
)

func main() {
	array := [...]int{0, 1, -2, 3, 4}

	// pointer变量指向array[0]的地址，array[0]是整型数组的第一个元素
	pointer := &array[0]
	fmt.Println(*pointer, " ")

	memoryAddress := uintptr(unsafe.Pointer(pointer)) + unsafe.Sizeof(array[0])
	for i := 0; i < len(array)-1; i++ {
		pointer = (*int)(unsafe.Pointer(memoryAddress))
		fmt.Print(*pointer, " ")
		memoryAddress = uintptr(unsafe.Pointer(pointer)) + unsafe.Sizeof(array[0])
	}

	// Note: 当你尝试访问无效的数组元素，程序并不会出错而是会返回一个随机的数字。
	pointer = (*int)(unsafe.Pointer(memoryAddress))
	fmt.Print("访问无效的数组元素，One more：", *pointer, " ")
}
