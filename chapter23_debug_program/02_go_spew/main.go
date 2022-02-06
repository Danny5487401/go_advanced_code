package main

import (
	"github.com/davecgh/go-spew/spew"
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
	spew.Dump(ins)

	// 2. 数组或者map中是指针对象时
	arr := [...]*Demo{{100, "Python"}, {200, "Golang"}}
	spew.Dump(arr)

	// 3. 循环结构
	c := &Circular{1, nil}
	c.next = &Circular{2, c}
	spew.Dump(c)

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
