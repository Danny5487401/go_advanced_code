
package main

import (
	"fmt"
	"time"
)

func main()  {
	a := 1
	go func() {
		a = 2
		fmt.Println("子goroutine。。",a) //子goroutine。。 2
	}()
	a = 3
	time.Sleep(1)
	fmt.Println("main goroutine。。",a)  // main goroutine。。 2 ---------noted:不是3
}

// 终端运行： go run --race main_test.go  可以看到WARNING: DATA RACE
