package main

import "fmt"

type A struct {
	s string
}

// 这是上面提到的 "在方法内把局部变量指针返回" 的情况
func foo(s string) *A {
	a := new(A)
	a.s = s
	return a //返回局部变量a,在C语言中妥妥野指针，但在go则ok，但a会逃逸到堆
}

//go build -gcflags=-m gcflags_main.go
func main() {
	a := foo("hello")
	b := a.s + " world"
	c := b + "!"
	fmt.Println(c) //这个fmt参数是个interface，所以c会逃逸到堆上
}

// go build -gcflags '-m -l -N' 02_params_ptr_return.go

/*
./01_params_ptr_return.go:10:10: leaking param: s
./01_params_ptr_return.go:11:10: new(A) escapes to heap
./01_params_ptr_return.go:19:11: a.s + " world" does not escape
./01_params_ptr_return.go:20:9: b + "!" escapes to heap
./01_params_ptr_return.go:21:13: ... argument does not escape
./01_params_ptr_return.go:21:13: c escapes to heap


*/
