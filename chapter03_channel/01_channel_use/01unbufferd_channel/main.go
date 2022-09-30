package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan bool) //0xc0000a4000,是引用类型的数据

	go func() {
		time.Sleep(time.Second * 5)
		// 循环结束后，向通道中写数据，表示要结束了。。
		fmt.Println("发送结束信号。。")
		ch1 <- true

	}()
	// 一个通道发送和接收数据，默认是阻塞的
	data := <-ch1 // 从ch1通道中读取数据
	fmt.Println("data-->", data)
	fmt.Println("main。。over。。。。")
}
