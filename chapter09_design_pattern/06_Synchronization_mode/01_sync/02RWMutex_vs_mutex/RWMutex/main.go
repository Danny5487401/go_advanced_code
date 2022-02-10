package main

import (
	"fmt"
	"sync"
	"time"
)

// demo：制作一个读多写少的例子，分别开启 3 个 goroutine 进行读和写，输出最终的读写次数
// 使用独写锁
var (
	count int
	//读写锁
	countGuard sync.RWMutex
)

func read(mapA map[string]string) {
	for {
		countGuard.RLock() //这里定义了一个读锁
		var _ string = mapA["name"]
		count += 1
		countGuard.RUnlock()
	}
}

func write(mapA map[string]string) {
	for {
		countGuard.Lock() //这里定义了一个写锁
		mapA["name"] = "johny"
		count += 1
		time.Sleep(time.Millisecond * 3)
		countGuard.Unlock()
	}
}

func main() {
	var num = 3
	var mapA = map[string]string{"name": ""}

	for i := 0; i < num; i++ {
		go read(mapA)
	}
	for i := 0; i < num; i++ {
		go write(mapA)
	}

	time.Sleep(time.Second * 3)
	fmt.Printf("最终读写次数：%d\n", count)
}

// 最终读写次数：11823, 与Mutex结果差距大概在 2 倍左右，读锁的效率要快很多
