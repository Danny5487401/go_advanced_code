package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	locker := new(sync.Mutex)
	cond := sync.NewCond(locker)

	// 多个协程
	for i := 0; i < 30; i++ {
		go func(x int) {
			// 加锁才能使用wait()
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
