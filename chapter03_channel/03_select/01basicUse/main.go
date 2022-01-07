package main

import (
	"fmt"
	"time"
)

/*
	分支语句：if，switch，select
	select 语句类似于 switch 语句，
		但是select会随机执行一个可运行的case。
		如果没有case可运行，它将阻塞，直到有case可运行。
*/

func main() {

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		time.Sleep(2 * time.Second)
		ch1 <- 100
	}()
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- 200
	}()
	// 如果有多个case都可以运行，select会随机公平地选出一个执行。其他不会执行
	select {
	case num1 := <-ch1:
		fmt.Println("ch1中取数据。。", num1)
	case num2, ok := <-ch2:
		if ok {
			fmt.Println("ch2中取数据。。", num2)
		} else {
			fmt.Println("ch2通道已经关闭。。")
		}
	}
}

/*
运行结果：可能执行第一个case，打印100，也可能执行第二个case，打印200。(多运行几次，结果就不同了)
*/
