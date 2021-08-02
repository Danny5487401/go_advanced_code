package main

import (
	"fmt"
	"sync"
	"time"
)

/*
Cond条件变量
	即等待或宣布事件发生的 goroutines 的会合点，它会保存一个通知列表。基本思想是当某中状态达成，goroutine 将会等待（Wait）在那里，
	当某个时刻状态改变时通过通知的方式（Broadcast，Signal）的方式通知等待的 goroutine。
	这样，不满足条件的 goroutine 唤醒继续向下执行，满足条件的重新进入等待序列。
与channel对比：
	提供了 Broadcast 方法，可以通知所有的等待者。
*/

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
一。NoCopy机制
	强调no copy的原因是为了安全，因为结构体对象中包含指针对象的话，直接赋值拷贝是浅拷贝，是极不安全的
	1.源码体现
		type Cond struct {
				noCopy noCopy  不允许copy

				// L is held while observing or changing the condition
				L Locker

				notify  notifyList
				checker copyChecker
			}
		// copyChecker holds back pointer to itself to detect object copying.
		type copyChecker uintptr

		func (c *copyChecker) check() {
			if uintptr(*c) != uintptr(unsafe.Pointer(c)) &&
				!atomic.CompareAndSwapUintptr((*uintptr)(c), 0, uintptr(unsafe.Pointer(c))) &&
				uintptr(*c) != uintptr(unsafe.Pointer(c)) {
				panic("sync.Cond is copied")
			}
		}
	2.工具
		go vet工具来检查，那么这个对象必须实现sync.Locker
		// A Locker represents an object that can be locked and unlocked.
		type Locker interface {
			Lock()
			Unlock()
		}
		type noCopy struct{}

		// Lock is a no-op used by -copylocks checker from `go vet`.
		func (*noCopy) Lock()   {}
		func (*noCopy) Unlock() {}
二。注意事项
	1. 不能不加锁直接调用 cond.Wait
		我们看到 Wait 内部会先调用 c.L.Unlock()，来先释放锁。如果调用方不先加锁的话，会触发“fatal error: sync: unlock of unlocked mutex”。
	2. 为什么不能 sync.Cond 不能复制 ？
		sync.Cond 不能被复制的原因，并不是因为 sync.Cond 内部嵌套了 Locker。因为 NewCond 时传入的 Mutex/RWMutex 指针，对于 Mutex 指针复制是没有问题的。
		主要原因是 sync.Cond 内部是维护着一个 notifyList。如果这个队列被复制的话，那么就在并发场景下导致不同 goroutine 之间操作的 notifyList.wait、notifyList.notify 并不是同一个，这会导致出现有些 goroutine 会一直堵塞
*/
