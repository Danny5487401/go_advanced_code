# Cond条件变量

    即等待或宣布事件发生的 goroutines 的会合点，它会保存一个通知列表。基本思想是当某中状态达成，goroutine 将会等待（Wait）在那里，
    当某个时刻状态改变时通过通知的方式（Broadcast，Signal）的方式通知等待的 goroutine。
    这样，不满足条件的 goroutine 唤醒继续向下执行，满足条件的重新进入等待序列。
与channel对比：

    提供了 Broadcast 方法，可以通知所有的等待者。

## 源码体现
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
