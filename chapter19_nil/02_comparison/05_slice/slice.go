package main

import (
	"fmt"
	"unsafe"
)

/*
slice底层的数据结构有三个属性，分别是指向底层数组的指针（数据存放的地址）、切片的长度和切片的容量，那么slice和nil比较究竟是比较什么呢？

	答案是，slice和nil进行比较实质上比较的是slice结构体中指向数据的指针是否为nil，本质上也是指针的地址比较
*/
type sliceTest struct {
	array unsafe.Pointer //指向底层数组的指针

	len int //切片的长度

	cap int //切片的容量

}

func main() {
	var s []byte
	(*sliceTest)(unsafe.Pointer(&s)).len = 10 // 赋值len
	fmt.Println(s == nil)                     // true

	(*sliceTest)(unsafe.Pointer(&s)).cap = 10 // 赋值cap
	fmt.Println(s == nil)                     // true

	(*sliceTest)(unsafe.Pointer(&s)).array = unsafe.Pointer(uintptr(0x1))
	fmt.Println(s == nil) // false

	(*sliceTest)(unsafe.Pointer(&s)).array = unsafe.Pointer(uintptr(0x0))
	fmt.Println(s == nil) // true
}
