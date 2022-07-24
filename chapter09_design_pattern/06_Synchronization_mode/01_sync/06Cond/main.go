package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var (
		locker sync.Mutex
		cond   = sync.NewCond(&locker)
		wg     sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(number int) {
			// wait()方法内部是先释放锁 然后在加锁 所以这里需要先 Lock()
			cond.L.Lock()
			defer cond.L.Unlock()
			cond.Wait() // 等待通知,阻塞当前 goroutine
			fmt.Printf("g %v ok~ \n", number)
			wg.Done()
		}(i)
	}
	for i := 0; i < 5; i++ {
		// 每过 50毫秒 唤醒一个 goroutine
		cond.Signal()
		time.Sleep(time.Millisecond * 50)
	}
	time.Sleep(time.Millisecond * 50)
	// 剩下5个 goroutine 一起唤醒
	cond.Broadcast()
	fmt.Println("Broadcast...")
	wg.Wait()
}
