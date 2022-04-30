package main

import (
	"fmt"
	"golang.org/x/exp"
	"math/rand"
)

// 为函数声明类型参数T, 其类型为any
func ZeroValue[T any]() {
	var x T
	fmt.Printf("zero value of %T type: %v\n", x, x)
}

func main() {
	ZeroValue[bool]()           // zero value of bool type: false
	ZeroValue[complex64]()      // zero value of complex64 type: (0+0i)
	ZeroValue[[3]int]()         // zero value of [3]int type: [0 0 0]
	ZeroValue[map[string]int]() // zero value of map[string]int type: map[]
	ZeroValue[struct {
		b bool
		e error
	}]() // zero value of struct { b bool; e error } type: {false <nil>}
}

// 作为函数输入参数或返回值的类型
func Fn0[T, U any](t T, u U) {
	fmt.Println(t, u)
}

func Fn1[T constraints.Integer]() (t T) {
	return T(rand.Int()) // 使用类型参数强制转换类型
}

func Fn2[K comparable, V any](m map[K]V) {
	fmt.Println(m)
}
