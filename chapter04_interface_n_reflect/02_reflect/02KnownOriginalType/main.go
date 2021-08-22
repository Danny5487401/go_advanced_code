// 从reflect.Value中获取接口interface的信息

// realValue := value.Interface().(已知的类型)

package main

import (
	"fmt"
	"reflect"
)

func main() {
	var num float64 = 1.2345

	pointer := reflect.ValueOf(&num)
	value := reflect.ValueOf(num)
	fmt.Println(pointer) // 0xc00009c058
	fmt.Println(value) // 1.2345

	// 可以理解为“强制转换”，但是需要注意的时候，转换的时候，如果转换的类型不完全符合，则直接panic
	// Golang 对类型要求非常严格，类型一定要完全符合
	// 如下两个，一个是*float64，一个是float64，如果弄混，则会panic
	/*
	1. 从 Value 到实例
		该方法最通用，用来将 Value 转换为空接口，该空接口内部存放具体类型实例
		可以使用接口类型查询去还原为具体的类型
		//func (v Value) Interface() （i interface{})
	 */

	convertPointer := pointer.Interface().(*float64)
	convertValue := value.Interface().(float64)

	fmt.Println(convertPointer) // 0xc00009c058
	fmt.Println(convertValue)   // 1.2345
}
