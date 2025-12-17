package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {

	// context.Background()  -->  cancelA --> cancelB

	now := time.Now()
	// 创建基础上下文ctxA
	ctxA, cancelA := context.WithCancel(context.Background())
	defer cancelA() // 确保调用基础上下文的取消函数

	// 创建具有超时机制的派生上下文ctxB
	ctxB, cancelB := context.WithTimeout(ctxA, 5*time.Second)
	defer cancelB() // 确保调用派生上下文的取消函数，即使ctxA被取消

	// 模拟在ctxB超时前，主动取消ctxA
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("ctxA is being canceled")
		cancelA()
	}()

	// 等待ctxB的结束
	<-ctxB.Done()
	fmt.Println("duration of ctxB is ", time.Since(now))

	// 检查ctxB的结束原因
	if err := ctxB.Err(); errors.Is(err, context.Canceled) {
		fmt.Println("ctxB was canceled due to ctxA's cancellation")
	} else if errors.Is(err, context.DeadlineExceeded) {
		fmt.Println("ctxB was canceled due to its own timeout")
	}
}

/*
ctxA is being canceled
duration of ctxB is  1.001365042s
ctxB was canceled due to ctxA's cancellation
*/
