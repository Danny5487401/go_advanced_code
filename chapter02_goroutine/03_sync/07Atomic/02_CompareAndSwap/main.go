package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
// CompareAndSwapInt32 executes the compare-and-swap operation for an int32 value.
func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
函数会先判断参数addr指向的值与参数old是否相等，如果相等，则用参数new替换参数addr的值。最后返回swapped是否替换成功。
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
