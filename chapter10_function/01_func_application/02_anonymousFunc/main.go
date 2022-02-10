package main

import (
	"errors"
	"fmt"
)

/* 闭包应用
闭包经常用于回调函数，当IO操作（例如从网络获取数据、文件读写)完成的时候，会对获取的数据进行某些操作，这些操作可以交给函数对象处理

*/

// Traveser 定义函数类型 用于排序
type Traveser func(ele interface{})

// SortByDescending 具体操作:降序排序数组元素
func SortByDescending(ele interface{}) {
	intSlice, ok := ele.([]int)
	if !ok {
		return
	}
	length := len(intSlice)
	for i := 0; i < length-1; i++ {
		isChange := false
		for j := 0; j < length-1-i; j++ {
			if intSlice[j] < intSlice[j+1] {
				isChange = true
				intSlice[j], intSlice[j+1] = intSlice[j+1], intSlice[j]
			}
		}

		if isChange == false {
			return
		}

	}
}

//  SortByAscending具体操作:升序排序数组元素
func SortByAscending(ele interface{}) {
	intSlice, ok := ele.([]int)
	if !ok {
		return
	}
	length := len(intSlice)

	for i := 0; i < length-1; i++ {
		isChange := false
		for j := 0; j < length-1-i; j++ {

			if intSlice[j] > intSlice[j+1] {
				isChange = true
				intSlice[j], intSlice[j+1] = intSlice[j+1], intSlice[j]
			}
		}

		if isChange == false {
			return
		}

	}
}

func Process(array interface{}, traveser Traveser) error {
	if array == nil {
		return errors.New("nil pointer")
	}
	var length int // 定义数组长度
	switch array.(type) {
	case []int:
		length = len(array.([]int))
	case []string:
		length = len(array.([]string))
	case []float32:
		length = len(array.([]float32))
	default:
		return errors.New("error type")
	}
	if length == 0 {
		return errors.New("len is zero")
	}
	traveser(array)
	return nil

}

//在一些公共的操作中经常会包含一些差异性的特殊操作，而这些差异性的操作可以用函数来进行封装。
func main() {
	// 1. int类型切片
	intSlice := make([]int, 0)
	intSlice = append(intSlice, 3, 1, 4, 2)

	Process(intSlice, SortByDescending)
	fmt.Println(intSlice) //[4 3 2 1]
	Process(intSlice, SortByAscending)
	fmt.Println(intSlice) //[1 2 3 4]

	// 2. string类型切片
	stringSlice := make([]string, 0)
	stringSlice = append(stringSlice, "hello", "world", "china")

	/*
	   具体操作:使用匿名函数封装输出操作
	*/
	Process(stringSlice, func(elem interface{}) {

		if slice, ok := elem.([]string); ok {
			for index, value := range slice {
				fmt.Println("index:", index, "  value:", value)
			}
		}
	})

	// 3. float32类型切片
	floatSlice := make([]float32, 0)
	floatSlice = append(floatSlice, 1.2, 3.4, 2.4)

	/*
	   具体操作:使用匿名函数封装自定义操作
	*/
	Process(floatSlice, func(elem interface{}) {

		if slice, ok := elem.([]float32); ok {
			for index, value := range slice {
				slice[index] = value * 2
			}
		}
	})
	fmt.Println(floatSlice) //[2.4 6.8 4.8]

}
