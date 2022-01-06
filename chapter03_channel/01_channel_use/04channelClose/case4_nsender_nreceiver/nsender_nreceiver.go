package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// N 个 sender，N个 receiver
// 如果第 3 种解决方案，由 receiver 直接关闭 stopCh 的话，就会重复关闭一个 channel，导致 panic。
//	因此需要增加一个中间人，M 个 receiver 都向它发送关闭 dataCh 的“请求”，中间人收到第一个请求后，
//	就会直接下达关闭 dataCh 的指令（通过关闭 stopCh，这时就不会发生重复关闭的情况，因为 stopCh 的发送方只有中间人一个）。
//	另外，这里的 N 个 sender 也可以向中间人发送关闭 dataCh 的请求。
func main() {
	rand.Seed(time.Now().UnixNano())
	const Max = 100000
	const NumSenders = 1000
	const NumReceivers = 10

	dataCh := make(chan int, 100)
	stopCh := make(chan struct{})

	// It must be a buffered channel
	toStop := make(chan string, 1)
	var stoppedBy string
	// moderator中间人
	go func() {

		stoppedBy = <-toStop

		close(stopCh)

	}()

	//n个发送者
	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			// 构建数据
			value := rand.Intn(Max)
			if value == 0 {
				//发送结束信号的条件
				select {
				case toStop <- "sender#" + id:
				default:
				}
				return
			}
			select {
			case <-stopCh:
				return
			case dataCh <- value:

			}
		}(strconv.Itoa(i))
	}

	//N个消费者
	for i := 0; i < NumReceivers; i++ {
		go func(id string) {
			for {
				select {
				case <-stopCh:
					return
				case value := <-dataCh:
					if value == Max-1 {
						select {
						case toStop <- "receiver#" + id:
						default:
							//防止阻塞用
						}
						return
					}
					fmt.Println(value)
				}
			}
		}(strconv.Itoa(i))
	}
	select {
	case <-time.After(time.Hour):
	}
}

/*
这里将 toStop 声明成了一个 缓冲型的 channel。假设 toStop 声明的是一个非缓冲型的 channel，那么第一个发送的关闭 dataCh 请求可能会丢失。
因为无论是 sender 还是 receiver 都是通过 select 语句来发送请求，如果中间人所在的 goroutine 没有准备好，那 select 语句就不会选中，直接走 default 选项，
什么也不做。这样，第一个关闭 dataCh 的请求就会丢失。
*/
