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
	// 指针类型转换一
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

	// 指针类型转换二
	/*
		[]byte和string其实内部的存储结构都是一样的，但 Go 语言的类型系统禁止他俩互换。如果借助unsafe.Pointer，我们就可以实现在零拷贝的情况下，
			将[]byte数组直接转换成string类型	、
	*/
	fmt.Println("-----------")
	bytes := []byte{104, 101, 108, 108, 111}

	g := unsafe.Pointer(&bytes) //强制转换成unsafe.Pointer，编译器不会报错
	str := *(*string)(g)        //然后强制转换成string类型的指针，再将这个指针的值当做string类型取出来
	fmt.Println(str)            //输出 "hello"
}
