package main

import (
	"fmt"
	"time"
)

func main()  {
	ch1 := make(chan int)
	go sendData(ch1)
	/*
		子goroutine，写出数据10个
				每写一个，阻塞一次，主程序读取一次，解除阻塞

		主goroutine：循环读
				每次读取一个，堵塞一次，子程序，写出一个，解除阻塞

		发送发，关闭通道的--->接收方，接收到的数据是该类型的零值，以及false
	*/
	//主程序中获取通道的数据
	for{
		time.Sleep(1*time.Second)
		v, ok := <- ch1 //其他goroutine，显示的调用close方法关闭通道。
		if !ok{
			fmt.Println("已经读取了所有的数据，", ok)
			break
		}
		fmt.Println("取出数据：",v, ok)
	}

	fmt.Println("main...over....")
}
func sendData(ch1 chan int)  {
	// 发送方：10条数据
	for i:=0;i<10 ;i++  {
		ch1 <- i//将i写入通道中
	}
	close(ch1) //将ch1通道关闭了。
}