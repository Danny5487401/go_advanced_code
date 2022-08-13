package main

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	// ...
	const MaxRandomNumber = 100000
	const NumReceivers = 10
	const NumSenders = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	// ...
	dataCh := make(chan int, 100)
	stopCh := make(chan struct{})
	// stopCh 是一个信号通道。
	// 它的发送者是下面的主持人 goroutine。
	// 它的接收者是 dataCh的所有发送者和接收者。
	toStop := make(chan string, 1) // 请注意，toStop 通道的缓存大小（容量）是 1。 这是为了避免第一个通知在主持人 goroutine 准备接收从 toStop 来的通知之前发送的情况下丢失
	// toStop 通道通常用来通知主持人去关闭信号通道( stopCh )。
	// 它的发送者是 dataCh的任意发送者和接收者。
	// 它的接收者是下面的主持人 goroutine

	var stoppedBy string

	// 主持人
	go func() {
		stoppedBy = <-toStop
		close(stopCh)
	}()

	// 发送者
	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			for {
				value := rand.Intn(MaxRandomNumber)
				if value == 0 {
					// 这里，一个用于通知主持人关闭信号通道的技巧。
					select {
					case toStop <- "sender#" + id:
					default:
					}
					return
				}

				// 这里的第一个 select 是为了尽可能早的尝试退出 goroutine。
				//这个 select 代码块有 1 个接受行为 的 case 和 一个将作为 Go 编译器的 try-receive 操作进行特别优化的默认分支。
				select {
				case <-stopCh:
					return
				default:
				}

				// 即使 stopCh 被关闭， 如果 发送到 dataCh 的操作没有阻塞，那么第二个 select 的第一个分支可能会在一些循环（和在理论上永远） 不会执行。
				// 这就是为什么上面的第一个 select 代码块是必须的。
				select {
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	// 接收者
	for i := 0; i < NumReceivers; i++ {
		go func(id string) {
			defer wgReceivers.Done()

			for {
				// 和发送者 goroutine 一样， 这里的第一个 select 是为了尽可能早的尝试退出 这个 goroutine。
				// 这个 select 代码块有 1 个发送行为 的 case 和 一个将作为 Go 编译器的 try-send 操作进行特别优化的默认分支。
				select {
				case <-stopCh:
					return
				default:
				}

				//  即使 stopCh 被关闭， 如果从 dataCh 接受数据不会阻塞，那么第二个 select 的第一分支可能会在一些循环（和理论上永远）不会被执行到。
				// 这就是为什么第一个 select 代码块是必要的。
				select {
				case <-stopCh:
					return
				case value := <-dataCh:
					if value == MaxRandomNumber-1 {
						// 同样的技巧用于通知主持人去关闭信号通道。
						select {
						case toStop <- "receiver#" + id:
						default:
						}
						return
					}

					log.Println(value)
				}
			}
		}(strconv.Itoa(i))
	}

	// ...
	wgReceivers.Wait()
	log.Println("stopped by", stoppedBy)
}
