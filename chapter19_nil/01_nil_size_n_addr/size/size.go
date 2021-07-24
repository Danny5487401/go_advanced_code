package main

import (
	"fmt"
	"unsafe"
)

func main() {
	//2. 各类型为nil时的大小
	var p *int = nil
	fmt.Println("int: ", unsafe.Sizeof(p)) // 8
	var c chan int = nil
	fmt.Println("chan int: ", unsafe.Sizeof(c)) //8
	var f func() = nil
	fmt.Println("func: ", unsafe.Sizeof(f)) //8
	var m map[int]int = nil
	fmt.Println("map: ", unsafe.Sizeof(m)) // 8
	var s []int = nil
	fmt.Println("slice: ", unsafe.Sizeof(s)) // 24
	var i interface{} = nil
	fmt.Println("interface: ", unsafe.Sizeof(i)) // 16

}

/*
以map、slice作为两个特例来进行说明。map定义时编译器返回的是指针类型，在64位操作系统上，指针类型会分配8个字节大小的空间。
slice输出是24，让我们来看一下slice底层的数据结构
	// slice底层数据结构
	type slice struct {

		array unsafe.Pointer //指向底层数组的指针

		len int //切片的长度

		cap int //切片的容量

	}
*/
