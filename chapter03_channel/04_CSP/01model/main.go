// CSP 是 Communicating Sequential Process 的简称，中文可以叫做通信顺序进程，是一种并发编程模型，是一个很强大的并发数据模型
package main

/*
Golang，其实只用到了 CSP 的很小一部分，即理论中的 Process/Channel（对应到语言中的 goroutine/channel）：
	这两个并发原语之间没有从属关系， Process 可以订阅任意个 Channel，Channel 也并不关心是哪个 Process 在利用它进行通信；
	Process 围绕 Channel 进行读写，形成一套有序阻塞和可预测的并发模型。
Go语言的CSP模型是由协程Goroutine与通道Channel实现:
1. Go协程goroutine: 是一种轻量线程，它不是操作系统的线程，而是将一个操作系统线程分段使用，通过调度器实现协作式调度。
	是一种绿色线程，微线程，它与Coroutine协程也有区别，能够在发现堵塞后启动新的微线程。
2. 通道channel: 类似Unix的Pipe，用于协程之间通讯和同步。协程之间虽然解耦，但是它们和Channel有着耦合
Goroutine 用于执行并发任务，channel 用于 goroutine 之间的同步、通信。
Channel 分为两种：带缓冲、不带缓冲。对不带缓冲的 channel 进行的操作实际上可以看作 “同步模式”，带缓冲的则称为 “异步模式”
 */
func main()  {
	
}
/*
Go线程实现 两级线程模型 MPG:
1. M指的是Machine，一个M直接关联了一个内核线程。
2. P指的是”processor”，代表了M所需的上下文环境，也是处理用户级代码逻辑的处理器。
3. G指的是Goroutine，其实本质上也是一种轻量级的线程。
 */
