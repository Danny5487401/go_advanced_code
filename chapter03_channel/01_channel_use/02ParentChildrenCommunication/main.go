package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)

	done := make(chan bool) // 结束通道
	go func() {
		fmt.Println("子goroutine执行。。。")
		time.Sleep(3 * time.Second)
		data := <-ch1 // 从通道中读取数据
		fmt.Println("data：", data)
		done <- true
	}()
	// 向子协程写数据。。
	time.Sleep(5 * time.Second)
	ch1 <- 100

	<-done
	fmt.Println("main。。over")

}
