package main

/*  channel发送不接收问题
原因：当接收者停止工作，发送者并不知道，还在傻傻地向下游发送数据。故而，我们需要一种机制去通知发送者。
做法：Go 可以通过 channel 的关闭向所有的接收者发送广播信息。
*/

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
			case out <- n:
			case <-done: // 监听关闭channel
				return
			}
		}
	}()
	return out
}

func main() {
	defer func() {
		time.Sleep(time.Second)
		fmt.Println("关闭时线程数量: ", runtime.NumGoroutine()) //程序关闭时线程数量
	}()

	//解决方式
	// Set up the pipeline.
	done := make(chan struct{})
	defer close(done) // 注意defer的顺序   用于发送关闭信号

	out := gen(done, 2, 3, 4, 5)

	for n := range out {
		fmt.Println(n)              // 2
		time.Sleep(5 * time.Second) // done thing, 可能异常中断接收
		if true {                   // if err != nil
			break
		}
	}
}
