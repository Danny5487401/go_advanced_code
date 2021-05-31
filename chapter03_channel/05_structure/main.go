package main

/* 数据结构 channel_design.png
1. 基于数组的循环队列，有缓冲的channel用它暂存数据
2. 基于链表的单向队列，用于保存阻塞在此channel上的goroutine
3. 锁，用于实现goroutine对channel并发安全，保证某一时刻只有1个goroutine操作channel

*/

/* 看图 channelStructure.png
1. buf是有缓冲的channel所特有的结构，用来存储缓存数据。是个循环链表
2. sendx和recvx用于记录buf这个循环链表中的~发送或者接收的~index
3. lock是个互斥锁。
4. recvq和sendq分别是接收(<-channel)或者发送(channel <- xxx)的goroutine抽象出来的结构体(sudog)的队列。是个双向链表

// 源码在： /runtime/chan.go 结构体是 hchan

type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters
	lock mutex
}

type waitq struct {
	first *sudog
	last  *sudog
}
源码分析PPT: https://speakerdeck.com/kavya719/understanding-channels
	1.channel的创建  make_chan.jpg
		创建channel实际上就是在内存中实例化了一个hchan的结构体，并返回一个ch指针，我们使用过程中channel在函数之间的传递都是用的这个指针，
		这就是为什么函数传递中无需使用channel的指针，而直接用channel就行了，因为channel本身就是一个指针
	2.各种场景的发送和接收
		见  blocked_into_sudog.jpg
	3.goroutine的调度
	4.goroutine的阻塞和唤醒
	5.channel和goroutine在select操作下



send,recv操作：
	注意：缓存链表中以上每一步的操作，都是需要加锁操作的！
	每一步的操作的细节可以细化为：
	• 第一，加锁
	• 第二，把数据从goroutine中copy到“队列”中(或者从队列中copy到goroutine中）。
	• 第三，释放锁
*/

/*  当channel缓存满了之后会发生什么?
Go调度原理web连接：https://i6448038.github.io/2017/12/04/golang-concurrency-principle/  Go的CSP并发模型--->Go线程实现模型MPG

*/
