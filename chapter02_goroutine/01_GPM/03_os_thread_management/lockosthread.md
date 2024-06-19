# 线程管理

Go 语言的运行时会通过调度器改变线程的所有权，它也提供了 runtime.LockOSThread 和 runtime.UnlockOSThread 让我们有能力绑定 Goroutine 和线程完成一些比较特殊的操作。Goroutine 应该在调用操作系统服务或者依赖线程状态的非 Go 语言库时调用 runtime.LockOSThread 函数11，例如：C 语言图形库等

## 背景

golang的scheduler可以理解为公平协作调度和抢占的综合体，他不支持优先级调度。当你开了几十万个goroutine，并且大多数协程已经在runq等待调度了, 那么如果你有一个重要的周期性的协程需要优先执行该怎么办？


可以借助runtime.LockOSThread()方法来绑定线程，绑定线程M后的好处在于，他可以由system kernel内核来调度，因为他本质是线程了


runtime.LockOSThread会锁定当前协程只跑在一个系统线程上，这个线程里也只跑该协程。


## 


