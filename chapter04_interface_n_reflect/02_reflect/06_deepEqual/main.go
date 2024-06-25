package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 内置 ==
	builtinEqual()
	// 反射 equal
	deepEqualUse()
}

func builtinEqual() {
	// 1.无隐式类型转换
	// int8AndInt16()

	// 2 浮点数的比较问题
	floatCompare()

	// 3 数组的长度视为类型的一部分，长度不同的两个数组是不同的类型，不能直接比较。
	arrayCompare()

	// 4 结构体来说，依次比较各个字段的值
	structCompare()

	// 5 使用type可以基于现有类型定义新的类型。新类型会根据它们的底层类型来比较
	typeNewCompare()
}

func typeNewCompare() {
	type myInt int
	var a myInt = 10
	var b myInt = 20
	var c myInt = 10
	fmt.Println(a == b) // false
	fmt.Println(a == c) // true
}

func structCompare() {
	aa := A{a: 1, b: "test1"}
	bb := A{a: 1, b: "test1"}
	cc := A{a: 1, b: "test2"}
	fmt.Println(aa == bb)
	fmt.Println(aa == cc)
}

func arrayCompare() {
	a := [4]int{1, 2, 3, 4}
	b := [4]int{1, 2, 3, 4}
	c := [4]int{1, 3, 4, 5}
	fmt.Println(a == b) // true
	fmt.Println(a == c) // false
}

type A struct {
	a int
	b string
}

func floatCompare() {
	var a float64 = 0.1
	var b float64 = 0.2
	var c float64 = 0.3
	fmt.Println("a+b == c 结果", a+b == c) // false
}

func int8AndInt16() {
	// var a int8
	// var b int16
	// 编译错误：invalid operation a == b (mismatched types int8 and int16)
	// fmt.Println(a == b)
}

func deepEqualUse() {
	// 1. 切片对比
	sliceEqual()

	// 2. map对比
	mapEqual()

	// 3. 自定义int比较
	intDeepEqual()

	// 4.带有环的数据对比
	circleEqual()
}

type link struct {
	value string
	tail  *link
}

func circleEqual() {
	// Circular linked lists a -> b -> a and c -> c.
	a, b, c := &link{value: "a"}, &link{value: "b"}, &link{value: "c"}
	a.tail, b.tail, c.tail = b, a, c
	fmt.Println(reflect.DeepEqual(a, b)) // "false"
	fmt.Println(reflect.DeepEqual(a, c)) // "false"
	fmt.Println(reflect.DeepEqual(c, c)) // "true"
}

func sliceEqual() {
	src := reflect.ValueOf([]int{10, 20, 32})
	dest := reflect.ValueOf([]int{1, 2, 3})

	cnt := reflect.Copy(dest, src)
	cnt += 1

	// DeepEqual is used to check two interfaces are equal or not
	res1 := reflect.DeepEqual(src, dest)
	fmt.Println("Is dest is equal to src:", res1)

	var a, b []string = nil, []string{}
	fmt.Println("nil值的slice 和非nil但是空的slice", reflect.DeepEqual(a, b)) // "false"

}

func mapEqual() {
	map_1 := map[int]string{
		200: "Anita",
		201: "Neha",
		203: "Suman",
		204: "Robin",
		205: "Rohit",
	}
	map_2 := map[int]string{
		200: "Anita",
		201: "Neha",
		203: "Suman",
		204: "Robin",
		205: "Rohit",
	}

	res1 := reflect.DeepEqual(map_1, map_2)
	fmt.Println("Is Map 1 is equal to Map 2:", res1)

	var c, d map[string]int = nil, make(map[string]int)
	fmt.Println("一个nil值的map和非nil值但是空的map", reflect.DeepEqual(c, d)) // "false"
}

type MyInt int
type YourInt int

func intDeepEqual() {
	m := MyInt(1)
	y := YourInt(1)

	fmt.Println("myInt and YrInt ", reflect.DeepEqual(m, y)) // false
}
