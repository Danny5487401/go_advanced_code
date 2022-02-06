package main

import (
	"fmt"
	"runtime"
	"time"
)

// Goexit()使用
func fun() {
	// 调度器确保所有已注册 defer延迟调用被执行。
	defer fmt.Println("defer。。。")

	//return           //终止此函数
	runtime.Goexit() //终止所在的协程

	fmt.Println("fun函数。。。")
}

func main() {
	//创建新建的协程
	go func() {
		fmt.Println("goroutine开始。。。")

		//调用了别的函数
		fun()

		fmt.Println("goroutine结束。。") //Goexit()之后 运行不到
	}() //别忘了()

	//睡一会儿，不让主协程结束
	time.Sleep(3 * time.Second)
}

/*结果：
goroutine开始。。。
defer。。。
*/
