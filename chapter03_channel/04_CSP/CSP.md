# CSP(Communicating Sequential Process )
中文可以叫做通信顺序进程，是一种并发编程模型，是一个很强大的并发数据模型.

有2个支持高并发的模型：CSP和Actor。鉴于Occam和Erlang都选用了CSP(来自Go FAQ)，并且效果不错，Go也选了CSP，但与前两者不同的是，Go把channel作为头等公民。

CSP 也是一门自定义的编程语言，作者定义了输入输出语句，用于 processes 间的通信（communication）。
processes 被认为是需要输入驱动，并且产生输出，供其他 processes 消费，processes 可以是进程、线程、甚至是代码块。
输入命令是：!，用来向 processes 写入；输出是：?，用来从 processes 读出。

## 并发编程模型

大多数的编程语言的并发编程模型是基于线程和内存同步访问控制，Go 的并发编程的模型则用 goroutine 和 channel 来替代。
goroutine来自协程的概念，让一组可复用的函数运行在一组线程之上，即使有协程阻塞，该线程的其他协程也可以被runtime调度，转移到其他可运行的线程上。
最关键的是，程序员看不到这些底层的细节，这就降低了编程的难度，提供了更容易的并发。


Golang，其实只用到了 CSP 的很小一部分，即理论中的 Process/Channel（对应到语言中的 goroutine/channel）：
这两个并发原语之间没有从属关系， Process 可以订阅任意个 Channel，Channel 也并不关心是哪个 Process 在利用它进行通信；
Process 围绕 Channel 进行读写，形成一套有序阻塞和可预测的并发模型。
Go语言的CSP模型是由协程Goroutine与通道Channel实现:

1. Go协程goroutine: 是一种轻量线程，它不是操作系统的线程，而是将一个操作系统线程分段使用，通过调度器实现协作式调度。
   是一种绿色线程，微线程，它与Coroutine协程也有区别，能够在发现堵塞后启动新的微线程。
2. 通道channel: 类似Unix的Pipe，用于协程之间通讯和同步。协程之间虽然解耦，但是它们和Channel有着耦合
   Goroutine 用于执行并发任务，channel 用于 goroutine 之间的同步、通信。
   Channel 分为两种：带缓冲、不带缓冲。对不带缓冲的 channel 进行的操作实际上可以看作 “同步模式”，带缓冲的则称为 “异步模式”

