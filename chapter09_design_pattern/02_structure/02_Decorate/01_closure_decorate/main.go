package main

import (
	"fmt"
)

// 装饰模式
// 方法一：使用闭包

type Decorator func(i int, s string) bool

// 被装饰的函数
func foo(i int, s string) bool {
	fmt.Println("___foo___")
	return true
}

// 增加的功能，装饰器1
func withDeco1(fn Decorator) Decorator {
	return func(i int, s string) bool {
		fmt.Println("开始装饰1")
		result := fn(i, s)
		fmt.Println("结束装饰1")
		return result
	}
}

// 增加的功能，装饰器2
func withDeco2(fn Decorator) Decorator {
	return func(i int, s string) bool {
		fmt.Println("开始装饰2")
		result := fn(i, s)
		fmt.Println("结束装饰2")
		return result
	}
}

// 多个装饰器的使用

type HttpHandlerDecorator func(Decorator) Decorator

func Handler(h Decorator, decors ...HttpHandlerDecorator) Decorator {
	for i := range decors {
		d := decors[len(decors)-1-i] // iterate in reverse
		h = d(h)
	}
	return h
}

func main() {
	//简单使用方式一
	// 开始用装饰器修饰自己
	//foo := withDeco2(withDeco1(foo))
	//// 调用自己
	//foo(1, "hello danny")

	//更优雅使用方式二
	handler := Handler(foo, withDeco1, withDeco2)
	handler(1, "hello danny")

}
