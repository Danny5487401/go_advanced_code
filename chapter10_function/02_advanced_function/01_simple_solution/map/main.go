package main

import (
	"fmt"
	"strings"
)

// 定义两个函数 函数func(s string_test) string_test 或 func(s string_test) in

func MapStrToStr(arr []string, fn func(s string) string) []string {
	var newArray = make([]string, 0)
	for _, it := range arr {
		newArray = append(newArray, fn(it))
	}
	return newArray
}

func MapStrToInt(arr []string, fn func(s string) int) []int {
	var newArray = make([]int, 0)
	for _, it := range arr {
		newArray = append(newArray, fn(it))
	}
	return newArray
}
func main() {
	var list = []string{"Danny", "Joy", "Michael"}

	// 转换大写
	x := MapStrToStr(list, func(s string) string {
		return strings.ToUpper(s)
	})

	fmt.Printf("%v\n", x)

	// 求字符串长度
	y := MapStrToInt(list, func(s string) int {
		return len(s)
	})
	fmt.Printf("%v\n", y)

}