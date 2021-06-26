package main

import (
	"fmt"
	"reflect"
)

/*
装饰器模式：
	装饰器模式创建了一个装饰类，用来包装原有的类，并在保持类方法签名完整的前提下，提供了额外的功能。
优点
	装饰器类和被装饰类可以独立发展，不会互相耦合
缺点
	多层装饰比较复杂

实现：
	通过反射实现装饰任意的函数
*/

// 需要装饰的函数1
func help1(name string) {
	fmt.Println("help1():", name)
}

// 需要装饰的函数2
func help2(age int, name string) {
	fmt.Printf("help2():name:%v,age:%v\n", name, age)
}

// 反射实现装饰器
func decorator(decoP, fn interface{}) {
	var decoratedFunc, targetFunc reflect.Value
	decoratedFunc = reflect.ValueOf(decoP).Elem() // 1.作为输出的装饰完毕后的函数
	targetFunc = reflect.ValueOf(fn)              //2.获取将要装饰的函数
	// 3.通过反射生成函数
	v := reflect.MakeFunc(targetFunc.Type(), func(in []reflect.Value) (out []reflect.Value) {
		fmt.Println("decorate self-code") //添加的逻辑
		out = targetFunc.Call(in)         //调用被装饰函数
		return
	})
	decoratedFunc.Set(v) // 4.将装饰过的存储在输出函数中
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
