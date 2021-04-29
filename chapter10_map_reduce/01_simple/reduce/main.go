package main

import "fmt"

func Reduce(arr []string, fn func(s string) int) int {
	sum := 0
	for _, it := range arr {
		sum += fn(it)
	}
	return sum
}

func main()  {
	var list = []string{"Danny", "Joy", "Michael"}
	// 求字符串总长
	x := Reduce(list, func(s string) int {
		return len(s)
	})
	fmt.Printf("%v\n", x)
}