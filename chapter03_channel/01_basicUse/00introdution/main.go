package main

import "fmt"

/*
channel的使用场景:把channel用在数据流动的地方
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
