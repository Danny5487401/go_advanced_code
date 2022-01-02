package main

import (
	"fmt"
	"unsafe"
)

type Foo struct {
	A int8 // 1
	B int8 // 1
	C int8 // 1
}

type Bar struct {
	x int32 // 4
	y *Foo  // 8
	z bool  // 1
}

type Demo1 struct {
	m struct{} // 0
	n int8     // 1
}

type Demo2 struct {
	n int8     // 1
	m struct{} // 0
}

func main() {
	var b1 Bar
	fmt.Println(unsafe.Sizeof(b1)) // 24

	// 结构体变量b1的对齐系数
	fmt.Println(unsafe.Alignof(b1)) // 8
	// b1每一个字段的对齐系数
	fmt.Println(unsafe.Alignof(b1.x)) // 4：表示此字段须按4的倍数对齐
	fmt.Println(unsafe.Alignof(b1.y)) // 8：表示此字段须按8的倍数对齐
	fmt.Println(unsafe.Alignof(b1.z)) // 1：表示此字段须按1的倍数对齐

	// 空结构体
	var d1 Demo1
	fmt.Println(unsafe.Sizeof(d1)) // 1

	var d2 Demo2
	fmt.Println(unsafe.Sizeof(d2)) // 2

}
