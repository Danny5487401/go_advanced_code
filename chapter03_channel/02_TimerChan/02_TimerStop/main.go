package main

import (
	"fmt"
	"time"
)

func main() {

	fmt.Println("-------------------------------")

	//新建计时器，5秒后触发
	timer2 := time.NewTimer(5 * time.Second)

	//新开启一个线程来处理触发后的事件
	go func() {

		//等触发时的信号
		<-timer2.C

		fmt.Println("Timer 2 结束。。")

	}()

	//由于上面的等待信号是在新线程中，所以代码会继续往下执行，停掉计时器
	time.Sleep(3 * time.Second)
	stop := timer2.Stop()

	if stop {
		fmt.Println("Timer 2 停止。。")
	}

}

/*
Timer 2 停止。。
*/
