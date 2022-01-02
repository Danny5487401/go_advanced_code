package main

// Goroutines泄漏

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
			out <- n //发送1个就阻塞了，下面close()调用不到
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
	out := gen(2, 3, 4, 5)

	for n := range out {
		fmt.Println(n)              // 2
		time.Sleep(5 * time.Second) // done thing, 可能异常中断接收
		// 模拟接收者没接受完就异常关闭
		if true { // if err != nil
			break
		}
	}
}
