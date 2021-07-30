package main

import (
	"fmt"
	"sync"
)

/*
使用close(ch)关闭所有下游协程
场景：
	退出时，显示通知所有协程退出
原理：
	所有读ch的协程都会收到close(ch)的信号
*/

func dataProducer(ch1 chan int, wg *sync.WaitGroup) {
	// 发送方：10条数据
	go func() {
		for i := 0; i < 5; i++ {
			ch1 <- i
		}
		// 关闭chan
		close(ch1)
		wg.Done()
	}()
}

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)
	dataProducer(ch, &wg)
	wg.Add(1)
	dataReceiver(ch, &wg)
	wg.Wait()

}

func dataReceiver(ch chan int, wg *sync.WaitGroup) {
	go func() {
		for {
			if dada, ok := <-ch; ok {
				fmt.Println(dada)
			} else {
				fmt.Println("chan 关闭")
				break
			}
		}
		wg.Done()
	}()
}

func isCanceled(cancelChan chan struct{}) bool {
	select {
	case <-cancelChan:
		return true
	default:
		return false
	}
}

// 关闭chan: 广播
func cancel(cancelChan chan struct{}) {
	close(cancelChan)
}
