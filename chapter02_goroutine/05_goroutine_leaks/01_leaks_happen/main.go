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
		close(out) // 没有调用到导致未关闭
	}()
	return out
}

func main() {
	defer func() {
		time.Sleep(5 * time.Second)
		fmt.Println("关闭时线程数量: ", runtime.NumGoroutine()) //程序关闭时线程数量
	}()

	// Set up the pipeline.开始调用
	out := gen(2, 3, 4, 5)

	for n := range out {
		fmt.Println("接受到数据:", n) // 2
		//time.Sleep(5 * time.Second) // done thing, 可能异常中断接收
		//模拟接收者没接受完就异常关闭
		if true { // if err != nil
			break
		}
	}

}

/*
泄漏的原因是 goroutine 操作 channel 后，处于发送或接收阻塞状态，而 channel 处于满或空的状态，一直得不到改变。
同时，垃圾回收器也不会回收此类资源，进而导致 goroutine 会一直处于等待队列中，不见天日。
*/
