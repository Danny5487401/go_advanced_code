package main

import "fmt"

// Map为nil时能够进行读取，但不能进行写入。

func main() {
	var m map[int]int
	// 读
	_ = m[10]

	// 写
	//m[10] = 10 // panic: assignment to entry in nil map

	var m1 = make(map[int]int)
	m1[20] = 20
	fmt.Println(m1[20])

}
