package main

import (
	"fmt"
)

func main() {
	x := 200

	// 1. 不带方法 eface
	var anyInfo interface{} = x
	fmt.Println(anyInfo) //200

	// 2. 带方法的interface iface
	g := Gopher{"Go", 1}
	var c coder = g
	fmt.Println(c) //{Go 1}

}

type coder interface {
	code()
	debug()
}

type Gopher struct {
	language string
	Level    int
}

func (p Gopher) code() {
	fmt.Printf("I am coding %s language\n", p.language)
}

func (p Gopher) debug() {
	fmt.Printf("I am debuging %s language\n", p.language)
}

/*
汇编：
	go tool compile -S -N -L pointer.go > main.s 2>&1
可以看到，main 函数里调用了两个函数
	func convT2E64(t *_type, elem unsafe.Pointer) (e eface)
	func convT2I(tab *itab, elem unsafe.Pointer) (i iface)
*/
