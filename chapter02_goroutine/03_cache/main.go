package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// 100000goroutine
func addConcurrent(bignum int) {
	var c int32
	atomic.StoreInt32(&c, 0)

	start := time.Now()

	for i := 0; i < bignum; i++ {
		go atomic.AddInt32(&c, 1)
	}
	for {
		if c == int32(bignum) {
			fmt.Println(time.Since(start))
			break
		}
	}
}

// cpu数量的goroutine
func addCPUNum(bignum int) {
	var c int32
	wg := &sync.WaitGroup{}
	core := runtime.NumCPU()
	start := time.Now()
	wg.Add(core)
	for i := 0; i < core; i++ {
		go func(wg *sync.WaitGroup) {
			for j := 0; j < bignum/core; j++ {
				atomic.AddInt32(&c, 1)
			}
			wg.Done()
		}(wg)

	}
	wg.Wait()
	fmt.Println(time.Since(start))
}

// 1个goroutine
func addOneThread(bignum int) {
	var c int32
	start := time.Now()
	for i := 0; i < bignum; i++ {
		atomic.AddInt32(&c, 1)
	}
	fmt.Println(time.Since(start))
}

func main() {

	bigNum := 100000
	addConcurrent(bigNum) //30.851988ms

	addCPUNum(bigNum)    // 1.472481ms
	addOneThread(bigNum) //558.033µs

}

/*
分析
	显然100000个goroutines处理这种cpu-bound的工作很不利，我之前go调度文章里讲过，线程上下文切换有延迟代价。io-bound处理可以在io wait的时候去切换别的线程做其他事情，但是对于cpu-bound，它会一直处理work，线程切换会损害性能。

	这里还有另外一个重要因素，那就是cache伪共享(false sharing)。每个core都会去共享变量c的相同cache行，频繁操作c会导致内存抖动(cache和主存直接的换页操作)。

	可以看到cpu number个线程并行处理时间是单线程处理时间的三倍，这个延迟代价还是很大的。

	因此，在golang程序中需要避免因为cache伪共享导致的内存抖动，尽量避免多个线程去频繁操作一个相同变量或者是地址相邻变量
*/
