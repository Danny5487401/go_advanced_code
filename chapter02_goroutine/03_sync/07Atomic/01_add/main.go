package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 加减操作
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
