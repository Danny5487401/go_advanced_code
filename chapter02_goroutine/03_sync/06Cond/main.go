package main

import (
	"fmt"
	"sync"
	"time"
)

/*
Cond 实现一个条件变量，即等待或宣布事件发生的 goroutines 的会合点，它会保存一个通知列表。
	基本思想是当某中状态达成，goroutine 将会等待（Wait）在那里，当某个时刻状态改变时通过通知的方式（Broadcast，Signal）的方式通知等待的 goroutine。
	这样，不满足条件的 goroutine 唤醒继续向下执行，满足条件的重新进入等待序列。
*/

func main() {
	locker := new(sync.Mutex)
	cond := sync.NewCond(locker)

	for i := 0; i < 30; i++ {
		go func(x int) {
			cond.L.Lock()
			fmt.Println(x, " 获取锁")
			defer cond.L.Unlock()
			cond.Wait()
			fmt.Println(x, " 被唤醒")
			time.Sleep(time.Second)
		}(i)
	}

	time.Sleep(time.Second)
	fmt.Println("Signal...")
	cond.Signal()
	time.Sleep(time.Second)
	cond.Signal()
	time.Sleep(time.Second * 3)
	cond.Broadcast()

	fmt.Println("Broadcast...")
	time.Sleep(time.Minute)
}
