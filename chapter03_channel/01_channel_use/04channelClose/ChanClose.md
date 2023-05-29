<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Channel关闭](#channel%E5%85%B3%E9%97%AD)
  - [关闭 channel 的原则](#%E5%85%B3%E9%97%AD-channel-%E7%9A%84%E5%8E%9F%E5%88%99)
  - [如何关闭](#%E5%A6%82%E4%BD%95%E5%85%B3%E9%97%AD)
    - [优雅的关闭](#%E4%BC%98%E9%9B%85%E7%9A%84%E5%85%B3%E9%97%AD)
  - [场景：](#%E5%9C%BA%E6%99%AF)
  - [原理](#%E5%8E%9F%E7%90%86)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Channel关闭

关于 channel 的使用，有几点不方便的地方：

1. 在不改变 channel 自身状态的情况下，无法获知一个 channel 是否关闭。

2. 关闭一个 closed channel 会导致 panic。所以，如果关闭 channel 的一方在不知道 channel 是否处于关闭状态时就去贸然关闭 channel 是很危险的事情。

3. 向一个 closed channel 发送数据会导致 panic。所以，如果向 channel 发送数据的一方不知道 channel 是否处于关闭状态时就去贸然向 channel 发送数据是很危险的事情。



## 关闭 channel 的原则
> don't close a channel from the receiver side and don't close a channel if the channel has multiple concurrent senders.

不要从一个 receiver 侧关闭 channel，也不要在有多个 sender 时，关闭 channel。


向 channel 发送元素的就是 sender，因此 sender 可以决定何时不发送数据，并且关闭 channel。
但是如果有多个 sender，某个 sender 同样没法确定其他 sender 的情况，这时也不能贸然关闭 channel。

但是上面所说的并不是最本质的，最本质的原则就只有一条：

> don't close (or send values to) closed channels.

## 如何关闭
有两个不那么优雅地关闭 channel 的方法：

1. 使用 defer-recover 机制，放心大胆地关闭 channel 或者向 channel 发送数据。即使发生了 panic，有 defer-recover 在兜底。

生产者关闭潜在的关闭通道
```go
func SafeClose(ch chan T) (justClosed bool) {
    defer func() {
        if recover() != nil {
            // 返回值可以被修改
            // 在一个延时函数的调用中。
            justClosed = false
        }
    }()

    // 假设这里 ch != nil 。
    close(ch)   // 如果 ch 已经被关闭将会引发 panic
    return true // <=> justClosed = true; return
}
```

生产者将数据发送到一个潜在的关闭通道
```go
func SafeSend(ch chan T, value T) (closed bool) {
    defer func() {
        if recover() != nil {
            closed = true
        }
    }()

    ch <- value  // 如果 ch 已经被关闭将会引发 panic
    return false // <=> closed = false; return
}
```

2. 使用 sync.Once 来保证只关闭一次
```go
type MyChannel struct {
    C    chan T
    once sync.Once
}

func NewMyChannel() *MyChannel {
    return &MyChannel{C: make(chan T)}
}

func (mc *MyChannel) SafeClose() {
    mc.once.Do(func() {
        close(mc.C)
    })
}
```
3. 使用 sync.Mutex 去避免多次关闭一个通道：
```go
type MyChannel struct {
    C      chan T
    closed bool
    mutex  sync.Mutex
}

func NewMyChannel() *MyChannel {
    return &MyChannel{C: make(chan T)}
}

func (mc *MyChannel) SafeClose() {
    mc.mutex.Lock()
    defer mc.mutex.Unlock()
    if !mc.closed {
        close(mc.C)
        mc.closed = true
    }
}

func (mc *MyChannel) IsClosed() bool {
    mc.mutex.Lock()
    defer mc.mutex.Unlock()
    return mc.closed
}
```

根据 sender 和 receiver 的个数，分下面几种情况：

- 1 一个 sender，一个 receiver
- 2 一个 sender， M 个 receiver
- 3 N 个 sender，一个 reciver
- 4 N 个 sender， M 个 receiver

### 优雅的关闭

对于 1，2，只有一个 sender 的情况就不用说了，直接从 sender 端关闭就好了，没有问题

第 3 种情形下，优雅关闭 channel 的方法是：the only receiver says "please stop sending more" by closing an additional signal channel。
解决方案就是增加一个传递关闭信号的 channel，receiver 通过信号 channel 下达关闭数据 channel 指令。senders 监听到关闭信号后，停止发送数据。

第 4 种情形下,如果采取第 3 种解决方案，由 receiver 直接关闭 stopCh 的话，就会重复关闭一个 channel，导致 panic。
因此需要增加一个中间人，M 个 receiver 都向它发送关闭 dataCh 的“请求”，中间人收到第一个请求后，就会直接下达关闭 dataCh 的指令（通过关闭 stopCh，这时就不会发生重复关闭的情况，因为 stopCh 的发送方只有中间人一个）。另外，这里的 N 个 sender 也可以向中间人发送关闭 dataCh 的请求。

## 场景：

退出时，显示通知所有协程退出

## 原理

所有读ch的协程都会收到close(ch)的信号
