# Timer定时器

我们不管用 NewTimer, timer.After，还是 timer.AfterFun 来初始化一个 timer, 这个 timer 最终都会加入到一个全局 timer 堆中， 由 Go runtime 统一管理。

全局的 timer 堆也经历过三个阶段的重要升级。

- Go 1.9 版本之前，所有的计时器由全局唯一的四叉堆维护，协程间竞争激烈。
- Go 1.10 - 1.13，全局使用 64 个四叉堆维护全部的计时器，没有本质解决 1.9 版本之前的问题
- Go 1.14 版本之后，每个 P 单独维护一个四叉堆

## 源码分析

### 四叉堆原理
四叉堆其实就是四叉树，Go timer 是如何维护四叉堆的呢？
- Go runtime 调度 timer 时，触发时间更早的 timer，要减少其查询次数，尽快被触发。所以四叉树的父节点的触发时间是一定小于子节点的。
- 四叉树顾名思义最多有四个子节点，为了兼顾四叉树插、删除、重排速度，所以四个兄弟节点间并不要求其按触发早晚排序。

两张动图简单演示下 timer 的插入和删除
1. 把 timer 插入堆
![](.timer_images/insert_timer.gif)
插入时，兄弟节点无顺序，会与父节点进行比较
2. 把 timer 从堆中删除
![](.timer_images/del_timer.gif)
删除节点20，用子节点48替代，比较兄弟节点，与最小节点10交换位置

### timer 是如何被调度的
![](.timer_images/timer_schedule.png)
- 调用 NewTimer，timer.After, timer.AfterFunc 生产 timer, 加入对应的 P 的堆上。
- 调用 timer.Stop, timer.Reset 改变对应的 timer 的状态。
- GMP 在调度周期内中会调用 checkTimers ，遍历该 P 的 timer 堆上的元素，根据对应 timer 的状态执行真的操作。

### timer 是如何加入到 timer 堆上的？
- 通过 NewTimer, time.After, timer.AfterFunc 初始化 timer 后，相关 timer 就会被放入到对应 p 的 timer 堆上。

- timer 已经被标记为 timerRemoved，调用了 timer.Reset(d)，这个 timer 也会重新被加入到 p 的 timer 堆上

- timer 还没到需要被执行的时间，被调用了 timer.Reset(d)，这个 timer 会被 GMP 调度探测到，先将该 timer 从 timer 堆上删除，然后重新加入到 timer 堆上

- STW 时，runtime 会释放不再使用的 p 的资源，p.destroy()->timer.moveTimers，将不再被使用的 p 的 timers 上有效的 timer
(状态是：timerWaiting，timerModifiedEarlier，timerModifiedLater) 都重新加入到一个新的 p 的 timer 上

### timer的执行流程
![](.timer_images/runtime_timer_process.png)
timer 的真正执行者是 GMP。GMP 会在每个调度周期内，通过 runtime.checkTimers 调用 timer.runtimer().
timer.runtimer 会检查该 p 的 timer 堆上的所有 timer，判断这些 timer 是否能被触发。

如果该 timer 能够被触发，会通过回调函数 sendTime 给 Timer 的 channel C 发一个当前时间，告诉我们这个 timer 已经被触发了。

如果是 ticker 的话，被触发后，会计算下一次要触发的时间，重新将 timer 加入 timer 堆中。


## Timer 使用中的坑
确实 timer 是我们开发中比较常用的工具，但是 timer 也是最容易导致内存泄露，CPU 狂飙的杀手之一。

不过仔细分析可以发现，其实能够造成问题就两个方面：

- 错误创建很多的 timer，导致资源浪费
- 由于 Stop 时不会主动关闭 C，导致程序阻塞

### 1 错误创建很多 timer，导致资源浪费
```go
func main() {
    for {
        // xxx 一些操作
        timeout := time.After(30 * time.Second)
        select {
        case <- someDone:
            // do something
        case <-timeout:
            return
        }
    }
}
```

上面这段代码是造成 timer 异常的最常见的写法，也是我们最容易忽略的写法。

造成问题的原因其实也很简单，因为 timer.After 底层是调用的 timer.NewTimer，NewTimer 生成 timer 后，会将 timer 放入到全局的 timer 堆中。

for 会创建出来数以万计的 timer 放入到 timer 堆中，导致机器内存暴涨，同时不管 GMP 周期 checkTimers，还是插入新的 timer 都会疯狂遍历 timer 堆，导致 CPU 异常。

要注意的是，不只 time.After 会生成 timer, NewTimer，time.AfterFunc 同样也会生成 timer 加入到 timer 中，也都要防止循环调用

解决办法: 使用 time.Reset 重置 timer，重复利用 timer。

```go
func main() {
    timer := time.NewTimer(time.Second * 5)    
    for {
        timer.Reset(time.Second * 5)

        select {
        case <- someDone:
            // do something
        case <-timer.C:
            return
        }
    }
}
```
我们已经知道 time.Reset 会重新设置 timer 的触发时间，然后将 timer 重新加入到 timer 堆中，等待被触发调用。

### 2 程序阻塞，造成内存或者 goroutine 泄露
```go
func main() {
    timer1 := time.NewTimer(2 * time.Second)
    <-timer1.C
    println("done")
}
```
上面的代码可以看出来，只有等待 timer 超时 "done" 才会输出，原理很简单：程序阻塞在 <-timer1.C 上，一直等待 timer 被触发时，回调函数 time.sendTime 才会发送一个当前时间到 timer1.C 上，程序才能继续往下执行。

不过使用 timer.Stop 的时候就要特别注意了，比如：
```go
func main() {
    timer1 := time.NewTimer(2 * time.Second)
    go func() {
        timer1.Stop()
    }()
    <-timer1.C

    println("done")
}
```
程序就会一直死锁了，因为 timer1.Stop 并不会关闭 channel C，使程序一直阻塞在 timer1.C 上。

上面这个例子过于简单了，试想下如果 <- timer1.C 是阻塞在子协程中，timer 被的 Stop 方法被调用，那么子协程可能就会被永远的阻塞在那里，造成 goroutine 泄露，内存泄露。

stop正确方式
```go
func main() {
    timer1 := time.NewTimer(2 * time.Second)
    go func() {
        if !timer1.Stop() {
            <-timer1.C
        }
    }()

    select {
    case <-timer1.C:
        fmt.Println("expired")
    default:
    }
    println("done")
}
```