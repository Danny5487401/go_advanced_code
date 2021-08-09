package main

import (
	"fmt"
	"sync"
)

/*
使用场景
	使用场景是批量发出 RPC 或者 HTTP 请求
*/
var wg sync.WaitGroup // 创建同步等待组对象
func main() {
	/*
		WaitGroup：同步等待组
			可以使用Add(),设置等待组中要 执行的子goroutine的数量，

			在main 函数中，使用wait(),让主程序处于等待状态。直到等待组中子程序执行完毕。解除阻塞

			子goroutine对应的函数中。wg.Done()，用于让等待组中的子程序的数量减1
	*/
	//设置等待组中，要执行的goroutine的数量
	wg.Add(2)
	go fun1()
	go fun2()
	fmt.Println("main进入阻塞状态。。。等待wg中的子goroutine结束。。")
	wg.Wait() //表示main goroutine进入等待，意味着阻塞
	fmt.Println("main，解除阻塞。。")

}
func fun1() {
	for i := 1; i <= 10; i++ {
		fmt.Printf("fun1.。。i:%v\n", i)
	}
	wg.Done() //给wg等待中的执行的goroutine数量减1.同Add(-1)
}
func fun2() {
	defer wg.Done()
	for j := 1; j <= 10; j++ {
		fmt.Printf("fun2..j：%v\n", j)
	}
}

/*
结论：
	sync.WaitGroup 必须在 sync.WaitGroup.Wait 方法返回之后才能被重新使用；
	sync.WaitGroup.Done 只是对 sync.WaitGroup.Add 方法的简单封装，我们可以向 sync.WaitGroup.Add 方法传入任意负数（需要保证计数器非负）快速将计数器归零以唤醒等待的 Goroutine；
	可以同时有多个 Goroutine 等待当前 sync.WaitGroup 计数器的归零，这些 Goroutine 会被同时唤醒
*/
