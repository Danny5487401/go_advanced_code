package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	locker := new(sync.Mutex)
	cond := sync.NewCond(locker)

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

/*

二。注意事项
	1. 不能不加锁直接调用 cond.Wait
		我们看到 Wait 内部会先调用 c.L.Unlock()，来先释放锁。如果调用方不先加锁的话，会触发“fatal error: sync: unlock of unlocked mutex”。
	2. 为什么不能 sync.Cond 不能复制 ？
		sync.Cond 不能被复制的原因，并不是因为 sync.Cond 内部嵌套了 Locker。因为 NewCond 时传入的 Mutex/RWMutex 指针，对于 Mutex 指针复制是没有问题的。
		主要原因是 sync.Cond 内部是维护着一个 notifyList。如果这个队列被复制的话，那么就在并发场景下导致不同 goroutine 之间操作的 notifyList.wait、notifyList.notify 并不是同一个，这会导致出现有些 goroutine 会一直堵塞
*/
