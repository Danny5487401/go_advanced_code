package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
一。cas算法
	AS算法是一种有名的无锁算法。无锁编程，即不使用锁的情况下实现多线程之间的变量同步，也就是在没有线程被阻塞的情况下实现变量的同步，
	所以也叫非阻塞同步（Non-blocking Synchronization）。CAS算法涉及到三个操作数
		1。需要读写的内存值V
		2。进行比较的值A
		3。拟写入的新值B
		当且仅当 V 的值等于 A时，CAS通过原子方式用新值B来更新V的值，否则不会执行任何操作（比较和替换是一个原子操作）。一般情况下是一个自旋操作，即不断的重试
二。案例
	// CompareAndSwapInt32 executes the compare-and-swap operation for an int32 value.
	func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
	函数会先判断参数addr指向的值与参数old是否相等，如果相等，则用参数new替换参数addr的值。最后返回swapped是否替换成功。
三。自旋锁
	自旋锁是指当一个线程在获取锁的时候，如果锁已经被其他线程获取，那么该线程将循环等待，然后不断地判断是否能够被成功获取，直到获取到锁才会退出循环。
	获取锁的线程一直处于活跃状态，但是并没有执行任何有效的任务，使用这种锁会造成busy-waiting。
	它是为实现保护共享资源而提出的一种锁机制。其实，自旋锁与互斥锁比较类似，它们都是为了解决某项资源的互斥使用。无论是互斥锁，还是自旋锁，
	在任何时刻，最多只能由一个保持者，也就说，在任何时刻最多只能有一个执行单元获得锁。但是两者在调度机制上略有不同。对于互斥锁，如果资源已经被占用，
	资源申请者只能进入睡眠状态。但是自旋锁不会引起调用者睡眠，如果自旋锁已经被别的执行单元保持，调用者就一直循环在那里看是否该自旋锁的保持者已经释放了锁
四。自旋锁源码
	type spinLock uint32
	func (sl *spinLock) Lock() {
		for !atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
			runtime.Gosched()
		}
	}
	func (sl *spinLock) Unlock() {
		atomic.StoreUint32((*uint32)(sl), 0)
	}
	func NewSpinLock() sync.Locker {
		var lock spinLock
		return &lock
	}
*/
var sum uint32 = 100
var wg sync.WaitGroup

func main() {
	for i := uint32(0); i < 100; i++ {
		wg.Add(1)
		go func(t uint32) {
			defer wg.Done()
			// 可以看到sum的值只改变了一次，只有当sum值为100的时候，CAS才将sum的值修改为了sum+1
			ok := atomic.CompareAndSwapUint32(&sum, 100, sum+1)
			if ok {
				fmt.Printf("第{%v}更新成功\n", t)
			}
		}(i)
	}
	wg.Wait()
	fmt.Println(sum)
}
