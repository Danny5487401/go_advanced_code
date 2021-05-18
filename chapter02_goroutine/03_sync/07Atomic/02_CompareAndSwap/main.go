package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var sum uint32 = 100
var wg sync.WaitGroup

func main() {
	for i := uint32(0); i < 100; i++ {
		wg.Add(1)
		go func(t uint32) {
			defer wg.Done()
			// 可以看到sum的值只改变了一次，只有当sum值为100的时候，CAS才将sum的值修改为了sum+1
			atomic.CompareAndSwapUint32(&sum, 100, sum+1)
		}(i)
	}
	wg.Wait()
	fmt.Println(sum)
}
