/*
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
package main

import "fmt"

func main() {
	var whatever [3]struct{}

	for i := range whatever {
		// 方式一： 调用函数 resource leak
		//defer fmt.Println(i)
		/* 结果
		2  1  0
		*/
		// 方式二 闭包
		defer func() {
			fmt.Println(i)
		}()
		/*结果
		2 2 2
		// 分析：Each time a "defer" statement executes, the function value and parameters to the call are evaluated as usual
		and saved anew but the actual function is not invoked
		每次延迟语句执行时,函数值和调用参数会像以往一样被评估和保存,但是实际函数并不会被调用.
		简单来说:defer后面的语句在执行的时候，函数调用的参数会被保存起来，但是不执行。
		*/
	}
}
