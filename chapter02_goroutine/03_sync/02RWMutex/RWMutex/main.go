package main

import (
	"fmt"
	"sync"
	"time"
)

/*
问题
	什么时候需要用到锁？
解决：
	当程序中就一个线程的时候，是不需要加锁的，但是通常实际的代码不会只是单线程，所以这个时候就需要用到锁了，那么关于锁的使用场景主要涉及到哪些呢？
	1.多个线程在读相同的数据时
	2.多个线程在写相同的数据时
	3.同一个资源，有读又有写
读写锁（sync.RWMutex）
	在读多写少的环境中，可以优先使用读写互斥锁（sync.RWMutex），它比互斥锁更加高效。sync 包中的 RWMutex 提供了读写互斥锁的封装
	分类:读锁和写锁
		如果设置了一个写锁，那么其它读的线程以及写的线程都拿不到锁，这个时候，与互斥锁的功能相同
		如果设置了一个读锁，那么其它写的线程是拿不到锁的，但是其它读的线程是可以拿到锁
*/

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
