package main

import "fmt"

type Coder interface {
	code()
}

type Gopher struct {
	name string
}

func (g Gopher) code() {
	fmt.Printf("%s is coding\n", g.name)
}

func main() {
	var c Coder
	fmt.Println(c == nil)
	fmt.Printf("c: %T, %v\n", c, c)

	var g *Gopher
	fmt.Println(g == nil)
	fmt.Printf("g: %T, %v\n", g, g)

	// 接口值 iface 的零值是指动态类型 iface.tab._type 和动态值 iface.data  都为 nil。
	// 当仅且当这两部分的值都为 nil 的情况下，这个接口值就才会被认为 接口值 == nil
	c = g
	fmt.Println(c == nil)
	fmt.Printf("c: %T, %v\n", c, c)
}
