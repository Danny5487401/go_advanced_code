package main

import (
	"fmt"
	"strings"
)

func filter(arr []int, fn func(n int) bool) []int {
	var newArray = []int{}
	for _, it := range arr {
		if fn(it) {
			newArray = append(newArray, it)
		}
	}
	return newArray
}

func mapStrToStr(arr []string, fn func(s string) string) []string {
	var newArray = make([]string, 0)
	for _, it := range arr {
		newArray = append(newArray, fn(it))
	}
	return newArray
}

func mapStrToInt(arr []string, fn func(s string) int) []int {
	var newArray = make([]int, 0)
	for _, it := range arr {
		newArray = append(newArray, fn(it))
	}
	return newArray

}

func reduce(arr []string, fn func(s string) int) int {
	sum := 0
	for _, it := range arr {
		sum += fn(it)
	}
	return sum
}

func main() {
	// 1. filter
	var intSet = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	out := filter(intSet, func(n int) bool {
		return n%2 == 1
	})
	fmt.Printf("过滤后数据：%v\n", out)
	out = filter(intSet, func(n int) bool {
		return n > 5
	})
	fmt.Printf("过滤后数据：%v\n", out)

	// 2. map
	var list = []string{"Danny", "Joy", "Michael"}
	// 转换大写
	x := mapStrToStr(list, func(s string) string {
		return strings.ToUpper(s)
	})
	fmt.Printf("转换大写后的数据：%v\n", x)
	// 求字符串长度
	y := mapStrToInt(list, func(s string) int {
		return len(s)
	})
	fmt.Printf("分别求长度后的数据：%v\n", y)

	// 3. reduce
	// 求字符串总长
	x := reduce(list, func(s string) int {
		return len(s)
	})
	fmt.Printf("%v\n", x)
}
