package main

import (
	"fmt"
	"math/rand"
	"time"
)

// N 个 sender，一个 receiver
// 增加一个传递关闭信号的 channel，receiver 通过信号 channel 下达关闭数据 channel 指令。senders 监听到关闭信号后，停止发送数据。
func main() {
	rand.Seed(time.Now().UnixNano())
	const Max = 100000
	const NumSenders = 1000

	dataCh := make(chan int, 100)
	stopCh := make(chan struct{})

	//n个发送者
	for i := 0; i < NumSenders; i++ {
		go func() {
			select {
			case <-stopCh:
				return
			case dataCh <- rand.Intn(Max):

			}
		}()
	}

	//一个消费者
	go func() {
		for value := range dataCh {
			if value == Max-1 {
				fmt.Println("发送停止信号给发送着")
				close(stopCh)
				return
			}
			//fmt.Println(value)
		}
	}()

	select {
	case <-time.After(time.Hour):
	}
}

/*
需要说明的是，上面的代码并没有明确关闭 dataCh。在 Go 语言中，对于一个 channel，如果最终没有任何 goroutine 引用它，
不管 channel 有没有被关闭，最终都会被 gc 回收。所以，在这种情形下，所谓的优雅地关闭 channel 就是不关闭 channel，让 gc 代劳
*/
