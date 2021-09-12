package main

import "fmt"

/*
定义：
	Each time a “defer” statement executes, the function value and parameters to the call are evaluated as usual and saved anew but the actual function is not invoked.
	Instead, deferred functions are invoked immediately before the surrounding function returns, in the reverse order they were deferred.
	If a deferred function value evaluates to nil, execution panics when the function is invoked, not when the “defer” statement is executed.

	每次defer语句执行的时候，会把函数“压栈”，函数参数会被拷贝下来；当外层函数（非代码块，如一个for循环）退出时，defer函数按照定义的逆序执行；如果defer执行的函数为nil, 那么会在最终调用函数的产生panic.
特性：
    1. 关键字 defer 用于注册延迟调用。
    2. 这些调用直到 return 前才被执。因此，可以用来做资源清理。
    3. 多个defer语句，按先进后出的方式执行。
    4. defer语句中的变量，在defer声明时就决定了
用途：
    1. 关闭文件句柄
    2. 锁资源释放
    3. 数据库连接释放
*/

type number int

func (n number) print()   { fmt.Println(n) }
func (n *number) pprint() { fmt.Println(*n) }

func main() {
	var n number

	// 对n直接求值，开始的时候n=0,
	defer n.print() //0

	// n是引用
	defer n.pprint() //3

	// 闭包
	defer func() { n.print() }()  //3
	defer func() { n.pprint() }() //3

	n = 3
}
