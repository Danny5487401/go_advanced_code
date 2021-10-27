package main

import (
	"fmt"
	"reflect"
)

/*
实现：
	通过反射实现装饰任意的函数
*/

// 需要装饰的函数1--一个参数
func help1(name string) {
	fmt.Println("help1():", name)
}

// 需要装饰的函数2--两个参数
func help2(age int, name string) {
	fmt.Printf("help2():name:%v,age:%v\n", name, age)
}

// 反射实现装饰器
// 第一个是出参 decoPtr ，就是完成修饰后的函数；
// 第二个是入参 fn ，就是需要修饰的函数。
func decorator(decoP, fn interface{}) {
	var decoratedFunc, targetFunc reflect.Value

	// 1.作为输出的装饰完毕后的函数
	decoratedFunc = reflect.ValueOf(decoP).Elem()

	//2.获取将要装饰的函数
	targetFunc = reflect.ValueOf(fn)

	// 3.通过反射生成函数
	v := reflect.MakeFunc(targetFunc.Type(), func(in []reflect.Value) (out []reflect.Value) {
		fmt.Println("start decorate self-code") //添加的逻辑
		out = targetFunc.Call(in)               //调用被装饰函数
		fmt.Println("end decorate self-code")   //添加的逻辑

		return
	})

	// 4.将装饰过的存储在输出函数中
	decoratedFunc.Set(v)
	return
}

func main() {
	testFunc := help1
	decorator(&testFunc, help1) //需要输出的函数testFunc
	testFunc("Joy")
	fmt.Println("------")

	testFunc1 := help2
	decorator(&testFunc1, help2)
	testFunc1(18, "danny")

}
