package main

import (
	"fmt"
	"sync"
)

/*
FAN-IN,FAN-OUT介绍:
	以汽车组装为例，汽车生产线上有个阶段是给小汽车装4个轮子，可以把这个阶段任务交给4个人同时去做，这4个人把轮子都装完后，再把汽车移动到生产线下一个阶段。
	这个过程中，就有任务的分发，和任务结果的收集。其中任务分发是FAN-OUT，任务收集是FAN-IN
	FAN-OUT模式：
		多个goroutine从同一个通道读取数据，直到该通道关闭。OUT是一种张开的模式，所以又被称为扇出，可以用来分发任务。
	FAN-IN模式：
		1个goroutine从多个通道读取数据，直到这些通道关闭。IN是一种收敛的模式，所以又被称为扇入，用来收集处理的结果。

代码实现：
	producer()保持不变，负责生产数据。
	square()也不变，负责计算平方值。
	修改main()，启动3个square，这3个square从producer生成的通道读数据，这是FAN-OUT。
	增加merge()，入参是3个square各自写数据的通道，给这3个通道分别启动1个协程，把数据写入到自己创建的通道，并返回该通道，这是FAN-IN。
*/

func producer(nums ...int) <-chan int {
	out := make(chan int)
	go func() {

		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func square(inCh <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range inCh {
			out <- n * n
		}
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	collect := func(in <-chan int) {
		defer wg.Done()
		for n := range in {
			out <- n
		}
	}
	wg.Add(len(cs))
	//FAN_IN
	for _, c := range cs {
		go collect(c)
	}
	// 错误方式：直接等待是bug，死锁，因为merge写了out，main却没有读
	// wg.Wait()
	// close(out)

	// 正确方式
	go func() {
		wg.Wait()
		close(out)
	}()

	return out

}

func main() {
	in := producer(1, 2, 3, 4, 5)
	// FAN-OUT
	c1 := square(in)
	c2 := square(in)
	c3 := square(in)

	// consumer
	for ret := range merge(c1, c2, c3) {
		fmt.Printf("%3d ", ret)
	}
}
