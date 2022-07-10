package main

import (
	"fmt"
	"reflect"
)

func main() {
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
