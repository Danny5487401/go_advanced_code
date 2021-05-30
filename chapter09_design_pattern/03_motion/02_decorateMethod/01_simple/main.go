package main

import "fmt"

type Decorator func(i int,s string) bool

// 被装饰的函数
func foo(i int,s string) bool {
	fmt.Println("___foo___")
	return true
}

// 增加的功能，装饰器
func withDeco(fn Decorator)Decorator  {
	return func(i int, s string) bool {
		fmt.Println("开始装饰")
		result := fn(i,s)
		fmt.Println("结束装饰")
		return result
	}
}

// 开始使用
func main()  {
	// 开始用装饰器修饰自己
	foo := withDeco(foo)
	// 调用自己
	foo(1,"hello danny")
	
}
