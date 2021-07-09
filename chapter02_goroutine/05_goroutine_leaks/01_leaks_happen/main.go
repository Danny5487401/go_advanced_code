package main

/*
Goroutines泄漏：
	谈到内存管理，go语言为我们处理好了大量的细节。Go语言编译器使用逃逸分析（escape analysis）来决变量的存储。
	运行时通过使用垃圾回收器跟踪和管理堆分配。虽然在应用程序中不产生内存泄漏并不可能，但几率大大降低。
	一种常见的内存泄漏就是goroutines泄漏。如果你启动了一个goroutine，你期望它最终会终止，但是它并没有，这就是goroutine泄漏了。
	它在应用程序的运行期内存在, 分配给该goroutine的内存无法释放。这是“永远不要在不知道它将如何停止的情况下启动一个goroutine”的建议背后的部分理由
*/

// 以下代码发生泄漏

import (
	"fmt"
	"runtime"
	"time"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)
	// 启动一个子线程
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out) // 未关闭
	}()
	return out
}

func main() {
	defer func() {
		fmt.Println("关闭时线程数量: ", runtime.NumGoroutine()) //程序关闭时线程数量
	}()

	// Set up the pipeline.开始调用
	out := gen(2, 3)

	for n := range out {
		fmt.Println(n)              // 2
		time.Sleep(5 * time.Second) // done thing, 可能异常中断接收
		// 模拟接收者异常关闭
		if true { // if err != nil
			break
		}
	}
}

/*
场景
	如果这段代码在常驻服务中执行，比如 http server，每接收到一个请求，便会启动一次 sayHello，时间流逝，每次启动的 goroutine 都得不到释放，
	你的服务将会离奔溃越来越近
泄露情况分类
	前面介绍的例子由于在 goroutine 运行死循环导致的泄露。接下来，我会按照并发的数据同步方式对泄露的各种情况进行分析。简单可归于两类，即：
	1.channel 导致的泄露
	2.传统同步机制导致的泄露
	传统同步机制主要指面向共享内存的同步机制，比如排它锁、共享锁等。这两种情况导致的泄露还是比较常见的。但是go 由于 defer 的存在，第二类情况，一般情况下还是比较容易避免的。

channel 引起的泄露分析:
1. 发送不接收
	理想情况下，我们希望接收者总能接收完所有发送的数据，这样就不会有任何问题。但现实是，一旦接收者发生异常退出，停止继续接收上游数据，发送者就会被阻塞
*/
