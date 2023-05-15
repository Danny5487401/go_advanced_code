package main

import (
	"fmt"
	"runtime"
	"time"
)

func gen(done chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n: //会阻塞
			case <-done: // 监听关闭channel
				return
			}
		}
	}()
	return out
}

func main() {
	defer func() {
		time.Sleep(time.Second * 2)                      // 保证全部关闭完后再打印
		fmt.Println("关闭时线程数量: ", runtime.NumGoroutine()) //程序关闭时线程数量
	}()

	// 解决方式
	// Set up the pipeline: 接受程序关闭信号
	done := make(chan struct{})
	defer close(done) // 注意defer的顺序   用于发送关闭信号

	out := gen(done, 2, 3, 4, 5)

	for n := range out {
		fmt.Println(n)              // 2
		time.Sleep(5 * time.Second) // done thing, 可能异常中断接收
		if true {                   // if err != nil
			// 突然关闭，只能获取到2
			break
		}
	}
}
