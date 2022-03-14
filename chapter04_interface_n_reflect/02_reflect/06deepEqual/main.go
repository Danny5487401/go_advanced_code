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

}

func sliceEqual() {
	src := reflect.ValueOf([]int{10, 20, 32})
	dest := reflect.ValueOf([]int{1, 2, 3})

	cnt := reflect.Copy(dest, src)
	cnt += 1

	// DeepEqual is used to check
	// two interfaces are eual or not
	res1 := reflect.DeepEqual(src, dest)
	fmt.Println("Is dest is equal to src:", res1)
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

	// DeepEqual is used to check
	// two interfaces are eual or not
	res1 := reflect.DeepEqual(map_1, map_2)
	fmt.Println("Is Map 1 is equal to Map 2:", res1)
}

type MyInt int
type YourInt int

func intDeepEqual() {
	m := MyInt(1)
	y := YourInt(1)

	fmt.Println("myInt and YrInt ", reflect.DeepEqual(m, y)) // false
}
