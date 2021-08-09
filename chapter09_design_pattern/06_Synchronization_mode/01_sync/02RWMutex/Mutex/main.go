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
/*
源码：
	type Mutex struct {
		state int32  // 当前互斥锁的状态
		sema  uint32 // 当前互斥锁的状态
	}
互斥锁的状态
	在默认情况下，互斥锁的所有状态位都是 0，int32 中的不同位分别表示了不同的状态：
	mutexLocked — 表示互斥锁的锁定状态；
	mutexWoken — 表示从正常模式被从唤醒；
	mutexStarving — 当前的互斥锁进入饥饿状态；
	waitersCount — 当前互斥锁上等待的 Goroutine 个数
加锁过程
	如果互斥锁的状态不是 0 时就会调用 sync.Mutex.lockSlow 尝试通过自旋（Spinnig）等方式等待锁的释放，该方法的主体是一个非常大 for 循环，这里将它分成几个部分介绍获取锁的过程：
	1。判断当前 Goroutine 能否进入自旋；
	2。通过自旋等待互斥锁的释放；
	3。计算互斥锁的最新状态；
	4。更新互斥锁的状态并获取锁
解锁过程
	1。当互斥锁已经被解锁时，调用 sync.Mutex.Unlock 会直接抛出异常；
	2。当互斥锁处于饥饿模式时，将锁的所有权交给队列中的下一个等待者，等待者会负责设置 mutexLocked 标志位；
	3。当互斥锁处于普通模式时，如果没有 Goroutine 等待锁的释放或者已经有被唤醒的 Goroutine 获得了锁，会直接返回；在其他情况下会通过 sync.runtime_Semrelease 唤醒对应的 Goroutine

*/
