package main

import "fmt"

// C1 是我们定义的约束，它声明了一个方法M1，以及两个可用作类型实参的类型（~int | ~int32）。
type C1 interface {
	~int | ~int32
	M1()
}

// 两个自定义类型T1和T2，两个类型都实现了M1方法

// T1 类型的底层类型为struct{}，这样就导致了虽然T1类型满足了约束C1的方法集合，但类型T因为底层类型并不是int或int32而不满足约束C1，这也就会导致foo(t)调用在编译阶段报错。
type T1 struct{}

func (t T1) M1() {
}

// T2 类型的底层类型为int
type T2 int

func (t T2) M1() {
}

func foo[P C1](t P) {
	fmt.Println("泛型初步测试", t)
}

func main() {
	//var t1 T1
	//foo(t1) // ./pointer.go:27:5: T does not implement C1

	var t2 T2
	foo(t2)
}
