package main

import (
	"fmt"
	"reflect"
)

func InvertSlice(args []reflect.Value) (result []reflect.Value) {
	inSlice, n := args[0], args[0].Len() // 入参 只有一个，出参 也只有一个
	outSlice := reflect.MakeSlice(inSlice.Type(), 0, n)
	for i := n - 1; i >= 0; i-- {
		element := inSlice.Index(i)
		outSlice = reflect.Append(outSlice, element)
	}
	return []reflect.Value{outSlice}
}

func Bind(p interface{}, f func([]reflect.Value) []reflect.Value) {

	invert := reflect.ValueOf(p).Elem()

	//Use of MakeFunc() method
	invert.Set(reflect.MakeFunc(invert.Type(), f))
}

// Main function
func main() {
	// 1 自己构建
	var invertInts func([]int) []int
	Bind(&invertInts, InvertSlice)
	fmt.Println(invertInts([]int{1, 2, 3, 4, 5, 6}))

	// 2 使用内置 swapper调换部分顺序
	s := []int{1, 2, 3, 4, 5} // 声明一个切片，元素排列为 [1 2 3 4 5]
	f := reflect.Swapper(s)   // 调用reflect.Swapper()方法，出参是一个方法
	f(0, 1)                   // 调用方法，将索引位 0、1的元素互换
	fmt.Println(s)            // 结果为[2 1 3 4 5]

}
