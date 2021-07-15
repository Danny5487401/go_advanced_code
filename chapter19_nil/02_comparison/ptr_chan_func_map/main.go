package main

import (
	"fmt"
	"unsafe"
)

/*
管道、函数、map定义后编译器会返回指针类型，这些类型与nil进行比较等价与指针与nil进行比较，而指针与nil进行比较就是地址间的比较。
这里需要注意，如果使用make()函数为map或管道分配了空间，则不为nil。
*/

func main() {
	var a int = 0 // 分配了地址
	var p *int = &a
	fmt.Println(p == nil) // false
	p = (*int)(unsafe.Pointer(uintptr(0x0)))
	fmt.Println(p == nil) // true

	// ==============分配了空间====================

	m := make(map[int]int)
	fmt.Println(m == nil) // false
	c := make(chan int)
	fmt.Println(c == nil) // false
}
