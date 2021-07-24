package main

import (
	"fmt"
)

/*
nil介绍
	nil是预定义标识，可以在Go语言标准文档builtin/builtin.go标准库中找到nil的定义。
	nil代表了指针pointer、通道channel、函数func、接口interface、map、切片slice类型变量的零值
源码
	// nil is a predeclared identifier representing the zero value for a
	// pointer, channel, func, interface, map, or slice type.
	var nil Type // Type must be a pointer, channel, func, interface, map, or slice type

	// Type is here for the purposes of documentation only. It is a stand-in
	// for any Go type, but represents the same type for any given function
	// invocation.
	type Type int
*/

func main() {
	// 1. 各类型为nil时的地址
	var p *int = nil
	var c chan int = nil
	var f func() = nil
	var m map[int]int = nil
	var s []int = nil
	var i interface{} = nil
	fmt.Printf("*int地址是%p\n", p)      // 0x0
	fmt.Printf("chan int地址是%p\n", c)  // 0x0
	fmt.Printf("函数地址是%p\n", f)        // 0x0
	fmt.Printf("map地址是%p\n", m)       // 0x0
	fmt.Printf("切片地址是%p\n", s)        // 0x0
	fmt.Printf("interface地址是%p\n", i) // %!p(<nil>)
	/*
		从代码中可以直观的看到指针、管道、函数、map、切片slice为nil时输出的地址都为0x0，可以验证不同类型nil值地址都是相同的。
		而其中比较特殊的是接口，输出的是%!p(<nil>)，大致的原因是因为nil的接口经由reflect.ValueOf()函数输出的类型为<invalid reflect.Value>，
		针对于这种类型Printf函数进行了特别的拼接最终得到%!p(<nil>)


	*/
	//实例化是有地址的
	var p1 = &People{}
	var p2 = People{}
	fmt.Printf("地址%p\n", p1) // 0xc00000c060
	fmt.Printf("%p\n", &p2)  // 0xc00000c080
	//fmt.Println(p1 == nil, &p2 == nil) // false false

}

type People struct {
	name string
	age  int
}
