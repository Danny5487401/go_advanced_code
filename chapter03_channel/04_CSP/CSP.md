#CSP(Communicating Sequential Process )
中文可以叫做通信顺序进程，是一种并发编程模型，是一个很强大的并发数据模型.

    CSP 也是一门自定义的编程语言，作者定义了输入输出语句，用于 processes 间的通信（communicatiton）。
    processes 被认为是需要输入驱动，并且产生输出，供其他 processes 消费，processes 可以是进程、线程、甚至是代码块。
    输入命令是：!，用来向 processes 写入；输出是：?，用来从 processes 读出。

##并发编程模型

    大多数的编程语言的并发编程模型是基于线程和内存同步访问控制，Go 的并发编程的模型则用 goroutine 和 channel 来替代。
    Goroutine 和线程类似，channel 和 mutex (用于内存同步访问控制)类似


Golang，其实只用到了 CSP 的很小一部分，即理论中的 Process/Channel（对应到语言中的 goroutine/channel）：
这两个并发原语之间没有从属关系， Process 可以订阅任意个 Channel，Channel 也并不关心是哪个 Process 在利用它进行通信；
Process 围绕 Channel 进行读写，形成一套有序阻塞和可预测的并发模型。
Go语言的CSP模型是由协程Goroutine与通道Channel实现:
1. Go协程goroutine: 是一种轻量线程，它不是操作系统的线程，而是将一个操作系统线程分段使用，通过调度器实现协作式调度。
   是一种绿色线程，微线程，它与Coroutine协程也有区别，能够在发现堵塞后启动新的微线程。
2. 通道channel: 类似Unix的Pipe，用于协程之间通讯和同步。协程之间虽然解耦，但是它们和Channel有着耦合
   Goroutine 用于执行并发任务，channel 用于 goroutine 之间的同步、通信。
   Channel 分为两种：带缓冲、不带缓冲。对不带缓冲的 channel 进行的操作实际上可以看作 “同步模式”，带缓冲的则称为 “异步模式”

