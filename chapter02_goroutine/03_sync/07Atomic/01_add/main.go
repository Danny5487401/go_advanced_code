package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
以Int32为例：
	// AddInt32 atomically adds delta to *addr and returns the new value.
	func AddInt32(addr *int32, delta int32) (new int32)
	AddInt32可以实现对元素的原子增加或减少，函数会直接在传递的地址上进行修改操作。
	addr需要修改的变量的地址，delta修改的差值[正或负数]，返回new修改之后的新值。
*/

// 加减操作
var sum int32 = 100
var wg sync.WaitGroup
var sub uint32 = 10000

func main() {
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 方式一： 结果sum的值是不定的，取决于sum的同步访问情况
			//sum += 1
			// 方式二： 结果是确定而且正确的，同一时间只有一个goroutine修改sum
			atomic.AddInt32(&sum, 2)
			//需要注意的是如果是uint32,unint64时,不能直接传负数，所以需要利用二进制补码机制
			// // atomic.Adduint32(&b, ^uint32(N-1)) //N为需要减少的正整数值
			atomic.AddUint32(&sub, ^uint32(2-1)) // 等价于 b -= 2
		}()
	}
	wg.Wait()
	fmt.Println("和为", sum)
	fmt.Println("减为", sub)
}
