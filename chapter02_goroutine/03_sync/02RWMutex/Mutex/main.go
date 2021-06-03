package main

import (
	"fmt"
	"sync"
	"time"
)

// demo：制作一个读多写少的例子，分别开启 3 个 goroutine 进行读和写，输出最终的读写次数

// 使用互斥锁：

var (
	count int
	//互斥锁
	countGuard sync.Mutex
)

func read(mapA map[string]string) {
	for {
		countGuard.Lock()
		var _ string = mapA["name"]
		count += 1
		countGuard.Unlock()
	}
}

func write(mapA map[string]string) {
	for {
		countGuard.Lock()
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

// 最终读写次数：4885
