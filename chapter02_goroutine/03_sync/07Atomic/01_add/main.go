package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// sync/atomic 库提供了原子操作的支持，原子操作直接有底层CPU硬件支持，因而一般要比基于操作系统API的锁方式效率高些。

var sum uint32 = 100
var wg sync.WaitGroup

func main() {
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 方式一： 结果sum的值是不定的，取决于sum的同步访问情况
			//sum += 1
			// 方式二： 结果是确定而且正确的，同一时间只有一个goroutine修改sum
			atomic.AddUint32(&sum, 1)
		}()
	}
	wg.Wait()
	fmt.Println(sum)
}
