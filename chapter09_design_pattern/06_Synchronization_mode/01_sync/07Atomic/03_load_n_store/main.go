package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*载入操作

背景：
	当我们对某个变量进行读取操作时，可能该变量正在被其他操作改变，或许我们读取的是被修改了一半的数据。
做法
	所以我们通过Load这类函数来确保我们正确的读取，函数名以Load为前缀，加具体类型名
// LoadInt32 atomically loads *addr.
1.func LoadInt32(addr *int32) (val int32)
	Load函数参数为需要读取的变量地址，返回值为读取的值
	Load和Store操作对应与变量的原子性读写，许多变量的读写无法在一个时钟周期内完成，而此时执行可能会被调度到其他线程，无法保证并发安全。Load只保证读取的不是正在写入的值，Store只保证写入是原子操作。所以在使用的时候要注意

// StoreInt32 atomically stores val into *addr.
2.func StoreInt32(addr *int32, val int32)
	Store函数参数为需要存储的变量地址及需要写入的值，存储某个值时，任何CPU都不会对该值进行读或写操作，存储操作总会成功，它不关心旧值是什么，与CompareAndSwap不同
*/

var c int32

func main() {
	wg := sync.WaitGroup{}
	//我们启100个goroutine
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tmp := atomic.LoadInt32(&c)
			if !atomic.CompareAndSwapInt32(&c, tmp, tmp+1) {
				fmt.Println("c 修改失败")
			}
		}()
	}
	wg.Wait()
	//c的值有可能不等于100，频繁修改变量值情况下，CAS操作有可能不成功。
	fmt.Println("c : ", c)

}
