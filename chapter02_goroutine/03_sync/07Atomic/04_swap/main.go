package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
func SwapInt32(addr *int32, new int32) (old int32)
Swap会直接执行赋值操作，并将原值作为返回值返回
*/

func main() {
	var e int32
	wg2 := sync.WaitGroup{}
	//开启10个goroutine
	for i := 0; i < 10; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			tmp := atomic.LoadInt32(&e)
			old := atomic.SwapInt32(&e, tmp+1)
			fmt.Println("e old : ", old)
		}()
	}
	wg2.Wait()
	fmt.Println("e : ", e)
}
