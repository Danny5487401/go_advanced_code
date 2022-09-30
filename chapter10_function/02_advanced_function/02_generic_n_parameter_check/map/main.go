package main

import (
	"fmt"
	"reflect"
)

// 带有类型检查

// 方式一：有返回值,返回一个全新的数组 – Transform()
func Transform(slice, function interface{}) interface{} {
	return transform(slice, function, false)
}

// 方式二：无返回值，“就地完成”
func TransformInPlace(slice, function interface{}) interface{} {
	return transform(slice, function, true)
}

func transform(slice, function interface{}, inPlace bool) interface{} {

	// 1. check the `slice` type is Slice
	sliceInType := reflect.ValueOf(slice)
	if sliceInType.Kind() != reflect.Slice {
		panic("transform: not slice")
	}
	// 2. check the function signature
	fn := reflect.ValueOf(function)
	elemType := sliceInType.Type().Elem()
	if !verifyFuncSignature(fn, elemType, nil) {
		panic("trasform: function must be of type func(" + sliceInType.Type().Elem().String() + ") outputElemType")
	}

	sliceOutType := sliceInType

	if !inPlace {
		// 新生成一个Slice
		sliceOutType = reflect.MakeSlice(reflect.SliceOf(fn.Type().Out(0)), sliceInType.Len(), sliceInType.Cap())
	}

	// 修改原先Slice
	for i := 0; i < sliceInType.Len(); i++ {
		sliceOutType.Index(i).Set(fn.Call([]reflect.Value{sliceInType.Index(i)})[0])
	}
	return sliceOutType.Interface()
}

// 检验是否是函数
func verifyFuncSignature(fn reflect.Value, types ...reflect.Type) bool {

	//Check it is a function
	if fn.Kind() != reflect.Func {
		return false
	}
	// 函数的入参数量和出参 数量 检查
	// NumIn() - returns a function type's input parameter count.
	// NumOut() - returns a function type's output parameter count.
	if (fn.Type().NumIn() != len(types)-1) || (fn.Type().NumOut() != 1) {
		return false
	}

	// In() - returns the type of a function type's i'th input parameter.
	for i := 0; i < len(types)-1; i++ {
		if fn.Type().In(i) != types[i] {
			return false
		}
	}
	// Out() - returns the type of a function type's i'th output parameter.
	outType := types[len(types)-1]
	if outType != nil && fn.Type().Out(0) != outType {
		return false
	}
	return true
}

type Employee struct {
	Name     string
	Age      int
	Vacation int
	Salary   int
}

// 入口函数
func main() {
	// 可以用于字符串数组
	list := []string{"1", "2", "3", "4", "5", "6"}
	result := Transform(list, func(str string) string {
		return str + str + str
	})
	fmt.Printf("原始切片:%+v\n", list)
	fmt.Printf("%+v\n", result)

	// 可以用于整形数组
	list1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	TransformInPlace(list1, func(a int) int {
		return a * 3
	})
	fmt.Printf("%+v\n", list1)

	//
	//// 用于结构体
	//var list = []Employee{
	//	{"Hao", 44, 0, 8000},
	//	{"Bob", 34, 10, 5000},
	//	{"Alice", 23, 5, 9000},
	//	{"Jack", 26, 0, 4000},
	//	{"Tom", 48, 9, 7500},
	//}
	//
	//result := TransformInPlace(list, func(e Employee) Employee {
	//	e.Salary += 1000
	//	e.Age += 1
	//	return e
	//})
}
