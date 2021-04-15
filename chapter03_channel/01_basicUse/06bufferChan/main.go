package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	/*
		非缓存通道：make(chan T)
		缓存通道：make(chan T ,size)
			缓存通道，理解为是队列：

		非缓存，发送还是接受，都是阻塞的
		缓存通道,缓存区的数据满了，才会阻塞状态。。

	*/
	ch1 := make(chan int)           //非缓存的通道
	fmt.Println(len(ch1), cap(ch1)) //0 0
	//ch1 <- 100//阻塞的，需要其他的goroutine解除阻塞，否则deadlock

	ch2 := make(chan int, 5)        //缓存的通道，缓存区大小是5
	fmt.Println(len(ch2), cap(ch2)) //0 5
	ch2 <- 100                      //
	fmt.Println(len(ch2), cap(ch2)) //1 5

	//ch2 <- 200
	//ch2 <- 300
	//ch2 <- 400
	//ch2 <- 500
	//ch2 <- 600
	fmt.Println("--------------")
	ch3 := make(chan string, 4)
	go sendData3(ch3)
	for {
		time.Sleep(1*time.Second)
		v, ok := <-ch3
		if !ok {
			fmt.Println("读完了，，", ok)
			break
		}
		fmt.Println("\t读取的数据是：", v)
	}

	fmt.Println("main...over...")
}

func sendData3(ch3 chan string) {
	for i := 0; i < 10; i++ {
		ch3 <- "数据" + strconv.Itoa(i)
		fmt.Println("子goroutine，写出第", i, "个数据")
	}
	close(ch3)
}
