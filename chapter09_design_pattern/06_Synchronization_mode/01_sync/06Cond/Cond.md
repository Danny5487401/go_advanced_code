# Cond条件变量

    即等待或宣布事件发生的 goroutines 的会合点，它会保存一个通知列表。基本思想是当某中状态达成，goroutine 将会等待（Wait）在那里，
    当某个时刻状态改变时通过通知的方式（Broadcast，Signal）的方式通知等待的 goroutine。
    这样，不满足条件的 goroutine 唤醒继续向下执行，满足条件的重新进入等待序列。
与channel对比：

    提供了 Broadcast 方法，可以通知所有的等待者。

## 背景
怎么去通知阻塞协程继续运行：
如果一个协程走到某个逻辑后，需要在某种条件达成后才能继续往下走，这怎么实现呢。对于单个协程，golang中的channel完全可以实现协程的通信。
```go

func main() {
	s := &signal{
		reRUn: make(chan  struct{}, 1),
	}
	go func() {
		time.Sleep(2 * time.Second)
		s.signalReRun()
	}()

	go func(s *signal) {
		//阻塞
		<- s.reRUn
		fmt.Println("goroutine start!")
	}(s)

	time.Sleep(5 * time.Second)

}

//通知协程继续执行
func (s *signal) signalReRun() {
	select {
	case s.reRUn <- struct{}{}:
		fmt.Println("notify wait goroutine!")
	default:
		fmt.Println("default")
		// The Channel is already full, so a reconnection attempt will occur.
	}
}


```

但如果是多个协程同时等待某个条件呢，怎么实现呢？golang中提供了sync.cond即条件变量来完成这件事。

## 源码体现
结构体
```go
type Cond struct {
    noCopy noCopy  //不允许copy

    // L is held while observing or changing the condition
    L Locker

    notify  notifyList
    checker copyChecker
}
// copyChecker holds back pointer to itself to detect object copying.
type copyChecker uintptr

func (c *copyChecker) check() {
    if uintptr(*c) != uintptr(unsafe.Pointer(c)) &&
        !atomic.CompareAndSwapUintptr((*uintptr)(c), 0, uintptr(unsafe.Pointer(c))) &&
        uintptr(*c) != uintptr(unsafe.Pointer(c)) {
        panic("sync.Cond is copied")
    }
}
```
## 流程
cond.Signal()
![](signal_process.png)

cond.Broadcast()
![](broadcast_process.png)

cond.Wait()
![](wait_process.png)


## 注意事项
1. 不能不加锁直接调用 cond.Wait
我们看到 Wait 内部会先调用 c.L.Unlock()，来先释放锁。如果调用方不先加锁的话，会触发“fatal error: sync: unlock of unlocked mutex”。

2. 为什么不能 sync.Cond 不能复制 ？
sync.Cond 不能被复制的原因，并不是因为 sync.Cond 内部嵌套了 Locker。因为 NewCond 时传入的 Mutex/RWMutex 指针，对于 Mutex 指针复制是没有问题的。
主要原因是 sync.Cond 内部是维护着一个 notifyList。如果这个队列被复制的话，那么就在并发场景下导致不同 goroutine 之间操作的 notifyList.wait、notifyList.notify 并不是同一个，这会导致出现有些 goroutine 会一直堵塞。

## NoCopy机制
    noCopy 是 go1.7 开始引入的一个静态检查机制。它不仅仅工作在运行时或标准库，同时也对用户代码有效
    强调no copy的原因是为了安全，因为结构体对象中包含指针对象的话，直接赋值拷贝是浅拷贝，是极不安全的

工具
    go vet工具来检查，那么这个对象必须实现sync.Locker
```go
// A Locker represents an object that can be locked and unlocked.
type Locker interface {
    Lock()
    Unlock()
}

// noCopy 用于嵌入一个结构体中来保证其第一次使用后不会被复制
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
```

## 第三方库实现-->熔断框架hystrix-go
因为我们有一个条件是最大并发控制，采用的是令牌的方式进行流量控制，每一个请求都要获取一个令牌，使用完毕要把令牌还回去
```go
{
    ticketCond := sync.NewCond(cmd)
    ticketChecked := false
    // When the caller extracts error from returned errChan, it's assumed that
    // the ticket's been returned to executorPool. Therefore, returnTicket() can
    // not run after cmd.errorWithFallback().
    returnTicket := func () {
        cmd.Lock()
        // Avoid releasing before a ticket is acquired.
        for !ticketChecked {
            ticketCond.Wait()
        }
        cmd.circuit.executorPool.Return(cmd.ticket)
        cmd.Unlock()
    }
}
```
使用sync.NewCond创建一个条件变量，用来协调通知你可以归还令牌了.