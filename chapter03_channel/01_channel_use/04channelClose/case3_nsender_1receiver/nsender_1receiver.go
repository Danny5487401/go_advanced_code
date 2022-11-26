package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	// ...
	const MaxRandomNumber = 100000
	const NumSenders = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(1)

	// ...
	dataCh := make(chan int, 100)
	stopCh := make(chan struct{})
	// stopCh 是一个信号通道。
	// 它的发送者是 dataCh 的接收者。
	// 它的接收者是 dataCh 的发送者。

	// 发送者
	for i := 0; i < NumSenders; i++ {
		go func() {
			for {
				// 第一个 select 是为了尽可能早的尝试退出 goroutine。
				//  事实上，在这个特殊的例子中，这不是必要的，所以它能省略。
				select {
				case <-stopCh:
					return
				default:
				}

				// 即使 stopCh 已经关闭，如果发送给 dataCh 没有阻塞，那么在第二个 select 中第一个分支可能会在一些循环中不会执行。
				// 但是在这里例子中是可接受的， 所以上面的第一个 select 代码块可以被省略。
				select {
				case <-stopCh:
					return
				case dataCh <- rand.Intn(MaxRandomNumber):
				}
			}
		}()
	}

	// 接收者
	go func() {
		defer wgReceivers.Done()

		for value := range dataCh {
			if value == MaxRandomNumber-1 {
				//  dataCh 通道的接收者也是 stopCh 通道的发送者。
				// 在这里关闭停止通道是安全的。.
				close(stopCh)
				return
			}

			log.Println(value)
		}
	}()

	// ...
	wgReceivers.Wait()
}

/*
Note： 如果没有明确关闭 dataCh，在 Go 语言中，对于一个 channel，如果最终没有任何 goroutine 引用它，
不管 channel 有没有被关闭，最终都会被 gc 回收。所以，在这种情形下，所谓的优雅地关闭 channel 就是不关闭 channel，让 gc 代劳
*/
