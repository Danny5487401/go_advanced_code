package main

import (
	"fmt"
)

func main() {
	// 1. 结构体中含有指针对象
	ins := Instance{
		A: "AAAA",
		B: 1000,
		C: &Inner{
			D: "DDDD",
			E: "EEEE",
		},
	}
	fmt.Println(ins) // {AAAA 1000 0x14000070000}
	// 由于 C 字段是指针，所以打印出来的是一个地址0xc000054020，而地址背后的数据却被隐藏了。显然，这对程序排查非常不友好

	// 2. 数组或者map中是指针对象时
	arr := [...]*Demo{{100, "Python"}, {200, "Golang"}}
	fmt.Printf("%v\n-----------------分割线-----------\n", arr)

	// 3. 循环结构
	c := &Circular{1, nil}
	c.next = &Circular{2, c}

	fmt.Printf("%+v\n----------------分割线-------------------\n", c)
}

type Circular struct {
	a    int
	next *Circular
}

type Instance struct {
	A string
	B int
	C *Inner
}

type Inner struct {
	D string
	E string
}

type Demo struct {
	a int
	b string
}
