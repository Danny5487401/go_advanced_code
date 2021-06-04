package main

import "fmt"

// Golang使用goroutine和channel简单、高效的解决并发问题，channel解决的是goroutine之间的通信。

/*
channel存在3种状态:
1. nil，未初始化的状态，只进行了声明，或者手动赋值为nil
2. active，正常的channel，可读或者可写
3. closed，已关闭，千万不要误认为关闭channel后，channel的值是nil

操作					nil的channel	正常channel					已关闭channel
<- ch				阻塞				成功或阻塞					读到零值
ch <-				阻塞				成功或阻塞					panic
close(ch)			panic			成功							panic

对于nil通道的情况，也并非完全遵循上表，有1个特殊场景：当nil的通道在select的某个case中时，这个case会阻塞，但不会造成死锁


channel的使用场景:
把channel用在数据流动的地方：

1. 消息传递、消息过滤
2. 信号广播
3. 事件订阅与广播
4. 请求、响应转发
5. 任务分发
6. 结果汇总
7. 并发控制
8. 同步与异步
*/

func main() {
	//未初始化
	var a chan int
	if a == nil {
		fmt.Println("channel 是 nil 的, 不能使用，需要先创建通道。。")
		a = make(chan int)
		fmt.Printf("数据类型是： %T\n", a)
	}
	// 初始化
	ch1 := make(chan int)
	fmt.Printf("%T,%p\n", ch1, ch1)

	test1(ch1)
}
func test1(ch chan int) {
	// channel是引用类型的数据，在作为参数传递的时候，传递的是内存地址
	fmt.Printf("%T,%p\n", ch, ch)
}
