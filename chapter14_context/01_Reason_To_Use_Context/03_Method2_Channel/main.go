package main

// 通道方式

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func worker(exitChan chan struct{}) {
Loop:
	for {
		fmt.Println("worker")
		select {
		case <-exitChan: // 等待接收上级通知
			break Loop
		default:
			time.Sleep(time.Second)
		}
	}
	// 如何接收外部命令实现退出
	wg.Done()
}

func main() {
	var exitChan = make(chan struct{})
	wg.Add(1)
	go worker(exitChan)
	// 如何优雅的实现结束子goroutine
	time.Sleep(time.Second * 3) // sleep3秒以免程序过快退出
	exitChan <- struct{}{}      // 给子goroutine发送退出信号
	close(exitChan)

	wg.Wait()
	fmt.Println("over")
}

//  管道方式存在的问题：
// 1. 使用全局变量在跨包调用时不容易实现规范和统一，需要维护一个共用的channel
