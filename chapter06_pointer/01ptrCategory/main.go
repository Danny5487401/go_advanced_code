package main

// 指针基本使用

import (
	"fmt"
	"reflect"
	"unsafe"
)

func test(m int) {
	var y int = 66
	y += m
}

func main() {
	// 通过指针修改值
	var x int = 99
	var p *int = &x

	fmt.Println(p)

	x = 100

	fmt.Println("x: ", x)
	fmt.Println("*p: ", *p)

	test(11)

	*p = 999

	fmt.Println("x: ", x)
	fmt.Println("*p: ", *p)

	/*
		*p 称为 解引用 或者 间接引用.

		*p = 999 是通过借助 x 变量的地址, 来操作 x 对应的空间.

		不管是 x 还是 *p , 我们操作的都是同一个空间.
	*/
	// 指针类型转换
	fmt.Println("-----------")
	v1 := uint(12)
	v2 := int(13)

	fmt.Println("值")
	fmt.Println(reflect.TypeOf(v1)) //uint
	fmt.Println(reflect.TypeOf(v2)) //int

	fmt.Println("指针")
	fmt.Println(reflect.TypeOf(&v1)) //*uint
	fmt.Println(reflect.TypeOf(&v2)) //*int

	//使用unsafe.Pointer进行类型的转换
	q := (*uint)(unsafe.Pointer(&v2)) // &v2 *int-->*unit转换

	fmt.Println(reflect.TypeOf(q)) // *unit
	fmt.Println(*q)                //13值不变
}
