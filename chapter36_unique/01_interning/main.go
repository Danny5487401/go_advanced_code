package main

import (
	"fmt"
	"unique"
)

func main() {

	intInterning()
	strInterning()
	structInterning()

}

func intInterning() {
	var a, b int = 5, 6
	h1 := unique.Make(a)
	h2 := unique.Make(a)
	h3 := unique.Make(b)
	fmt.Println(h1 == h2) // true
	fmt.Println(h1 == h3) // false
}

func strInterning() {
	// 创建唯一Handle
	s1 := unique.Make("hello")
	s2 := unique.Make("world")
	s3 := unique.Make("hello")

	// s1和s3是相等的，因为它们是同一个字符串值
	fmt.Println(s1 == s3) // true
	fmt.Println(s1 == s2) // false

	// 从Handle获取原始值
	fmt.Println(s1.Value()) // hello
	fmt.Println(s2.Value()) // world
}

type UserType struct {
	a int
	z float64
	s string
}

func structInterning() {
	var u1 = UserType{
		a: 5,
		z: 3.14,
		s: "golang",
	}
	var u2 = UserType{
		a: 5,
		z: 3.15,
		s: "golang",
	}
	h1 := unique.Make(u1)
	h2 := unique.Make(u1)
	h3 := unique.Make(u2)
	fmt.Println(h1 == h2) // true
	fmt.Println(h1 == h3) // false
}
