package main

import "fmt"

type C1 interface {
	~int | ~int32
	M1()
}

type T struct{}

func (t T) M1() {
}

type T1 int

func (t T1) M1() {
}

func foo[P C1](t P) {
	fmt.Println("泛型初步测试", t)
}

func main() {
	var t1 T1
	foo(t1)

	//var t2 T
	//foo(t2)  // ./pointer.go:27:5: T does not implement C1
}
