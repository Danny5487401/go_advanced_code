package main

import "fmt"

/*
Golang的并发核心思路
	Golang并发核心思路是关注数据流动。数据流动的过程交给channel，数据处理的每个环节都交给goroutine，把这些流程画起来，有始有终形成一条线，那就能构成流水线模型。
*/

/*
需求：
	计算一个整数切片中元素的平方值并把它打印出来。
方式一：
	非并发的方式是使用for遍历整个切片，然后计算平方，打印结果。
方式二：
	使用流水线模型实现这个简单的功能，从流水线的角度，可以分为3个阶段：

		1.遍历切片，这是生产者。
		2.计算平方值。
		3.打印结果，这是消费者。
代码实现：
	producer()负责生产数据，它会把数据写入通道，并把它写数据的通道返回。
	square()负责从某个通道读数字，然后计算平方，将结果写入通道，并把它的输出通道返回。
	main()负责启动producer和square，并且还是消费者，读取square的结果，并打印出来。
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

func main() {
	intChan := producer(1, 2, 3)
	outChan := square(intChan)

	for ret := range outChan {
		fmt.Printf("%3d", ret)
	}

}
