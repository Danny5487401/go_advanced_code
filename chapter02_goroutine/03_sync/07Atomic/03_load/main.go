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
