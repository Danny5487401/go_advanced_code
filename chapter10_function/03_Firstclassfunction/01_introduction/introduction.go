package main

import (
	"fmt"
	"time"
)

/*
一等公民：
	函数作为变量对待。也就说，函数与变量没有差别，它们是一样的，变量出现的地方都可以替换成函数，并且编译也是可以通过的，没有任何语法问题。

函数使用：
	1. 函数可以定义函数类型
	2. 函数可以赋值给变量
	3. 高阶函数---可以作为入参也可以作为返回值
	4. 动态创建函数
	5. 匿名函数
	6. 闭包

1.定义函数类型：
	type Operation func(a,b int) int
-----Operation :type name类型名称
-----func(a,b int) int:signature函数签名
	func Add func(a,b int) int{
		return a+b
		}符合函数签名的函数

2.声明函数类型的变量和为变量赋值：
	var op Operation
	op = Add
	fmt.Println(op(1,2))
	变量op是Operation类型的，可以把Add作为值赋值给变量op，执行op等价于执行Add。

3.函数作为其他函数入参
	type Calculator struct {
		v int
	}
	func (c Calculator)Do(op Operation,a int){
		c.v = op(c.v,a)
	}
	func main(){
		var calc Calculator
		calc.Do(add,1)

4. 函数作为返回值+动态创建
	type Operation func(b int)int
	func Add(b int) Operation{
		addB := func(a int)int{
			return a + b
		}
		return addB
	}

	type Calculator struct {
		v int
	}
	func (c Calculator)Do(op Operation){
		c.v = op(c.v)
	}
	func main(){
		var calc Calculator
		calc.Do(add(1))  //c.v = 1

5. 匿名函数
	func(a int)int{}
	func Add(b int) Operation{
		return func(a int)int{
			return a + b
		}
	}

6. 闭包
	定义：
		闭包指有权访问另一个函数作用域中的变量的函数。大白话就是，可以创建1个函数，它可以访问其他函数遍历，但不需要传值。

	type Operation func(b int)int
	func Add(b int) Operation{
		addB := func(a int)int{
			return a + b
		}
		return addB
	}

	比如匿名函数里直接使用了变量b，该匿名函数也是闭包函数。
	Note:一个函数可以是匿名函数，但不是闭包函数，因为闭包有时是有副作用的。
*/

func main() {
	// 闭包副作用
	//s1 := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	//for i, v := range s1 {
	//	go func() {
	//		fmt.Printf("%d %d\n", i, v)
	//	}()
	//}
	//time.Sleep(time.Second)
	/*
		结果也许不是9 9 9….，因为这个goroutine的调度有关
	*/

	// 修改后写法
	s2 := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i, v := range s2 {
		go func(a, b int) {
			fmt.Printf("%d %d\n", a, b)
		}(i, v)
	}
	time.Sleep(time.Second)

}
