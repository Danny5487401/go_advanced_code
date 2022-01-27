package main

import (
	"fmt"
	"sync"
	"time"

	"context"
)

var wg sync.WaitGroup

func worker1(ctx context.Context) {
Loop:
	for {
		fmt.Println("worker1")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			fmt.Println("worker1 结束")
			break Loop
		default:

		}
	}

	// 如何接收外部命令实现退出
	wg.Done()
}

func worker2(ctx context.Context) {
Loop:
	for {
		fmt.Println("worker2")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			fmt.Println("worker2 结束")
			break Loop
		default:

		}
	}
	// 如何接收外部命令实现退出
	wg.Done()
}

func main() {
	wg.Add(2)
	ctx, cancel := context.WithCancel(context.Background())
	go worker1(ctx)
	go worker2(ctx)
	// 如何优雅的实现结束子goroutine
	time.Sleep(time.Second * 3)
	cancel() // 通知子goroutine结束

	wg.Wait()
	fmt.Println("over")
}
