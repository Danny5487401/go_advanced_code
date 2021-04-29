package main

import (
	"fmt"
	"reflect"
	"strings"
)

// 泛型

// 背景： 上面的Map-Reduce都因为要处理数据的类型不同而需要写出不同版本的Map-Reduce，虽然他们的代码看上去是很类似的
// 做法：Go泛型(type parameter)将在Go 1.18版本落地，即2022.2月份.目前的Go语言的泛型只能用 interface{} + reflect来完成，interface{}
//	可以理解为C中的 void*，Java中的 Object ，reflect是Go的反射机制包，用于在运行时检查类型。

// 1, 无类型检查

func Map(data interface{}, fn interface{}) []interface{} {
	vfn := reflect.ValueOf(fn)

	vData := reflect.ValueOf(data)
	result := make([]interface{}, vData.Len())

	for i := 0; i < vData.Len(); i++ {
		// 调用函数
		// 通过 vfn.Call() 方法来调用函数，通过 []refelct.Value{vData.Index(i)}来获得数据
		result[i] = vfn.Call([]reflect.Value{vData.Index(i)})[0].Interface()
	}
	return result
}

// 主体逻辑
func main()  {
	// 乘方
	square := func(x int) int {
		return x * x
	}

	nums := []int{1, 2, 3, 4}
	squaredArr := Map(nums,square)
	fmt.Println(squaredArr)
	//[1 4 9 16]


	// 大写
	upcase := func(s string) string {
		return strings.ToUpper(s)
	}
	strs := []string{"Danny", "Joy", "Michael"}
	upstrs := Map(strs, upcase)
	fmt.Println(upstrs)
}
