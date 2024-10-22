<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [信号量(semaphore）](#%E4%BF%A1%E5%8F%B7%E9%87%8Fsemaphore)
  - [工作原理](#%E5%B7%A5%E4%BD%9C%E5%8E%9F%E7%90%86)
  - [适用场景](#%E9%80%82%E7%94%A8%E5%9C%BA%E6%99%AF)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [Acquire请求信号量资源](#acquire%E8%AF%B7%E6%B1%82%E4%BF%A1%E5%8F%B7%E9%87%8F%E8%B5%84%E6%BA%90)
    - [Release归还信号量资源](#release%E5%BD%92%E8%BF%98%E4%BF%A1%E5%8F%B7%E9%87%8F%E8%B5%84%E6%BA%90)
  - [参考资料](#%E5%8F%82%E8%80%83%E8%B5%84%E6%96%99)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 信号量(semaphore）

信号量(Semaphore)，有时被称为信号灯，是[多线程环境下使用的一种设施，是可以用来保证两个或多个关键代码段不被并发调用。
在进入一个关键代码段之前，线程必须获取一个信号量；一旦该关键代码段完成了，那么该线程必须释放信号量。
其它想进入该关键代码段的线程必须等待直到第一个线程释放信号量。为了完成这个过程，需要创建一个信号量VI，
然后将Acquire Semaphore VI以及Release Semaphore VI分别放置在每个关键代码段的首末端。确认这些信号量VI引用的是初始创建的信号量

是一个同步对象，用于保持在0至指定最大值之间的一个计数值。
当线程完成一次对该semaphore对象的等待（wait）时，该计数值减一；
当线程完成一次对semaphore对象的释放（release）时，计数值加一。
当计数值为0，则线程等待该semaphore对象不再能成功直至该semaphore对象变成signaled状态。
semaphore对象的计数值大于0，为signaled状态；计数值等于0，为nonsignaled状态.

## 工作原理

信号量是由操作系统来维护的，信号量只能进行两种操作等待和发送信号，操作总结来说，核心就是PV操作：
	
- P原语：P是荷兰语Proberen(测试)的首字母。为阻塞原语，负责把当前进程由运行状态转换为阻塞状态，直到另外一个进程唤醒它。
		操作为：申请一个空闲资源(把信号量减1)，若成功，则退出；若失败，则该进程被阻塞；

- V原语：V是荷兰语Verhogen(增加)的首字母。为唤醒原语，负责把一个被阻塞的进程唤醒，它有一个参数表，存放着等待被唤醒的进程信息。
		操作为：释放一个被占用的资源(把信号量加1)，如果发现有被阻塞的进程，则选择一个唤醒之。

在信号量进行PV操作时都为原子操作，并且在PV原语执行期间不允许有中断的发生。

PV原语对信号量的操作可以分为三种情况：

1. 把信号量视为某种类型的共享资源的剩余个数，实现对一类共享资源的访问
2. 量用作进程间的同步
3. 视信号量为一个加锁标志，实现对一个共享变量的访问


## 适用场景

semaphore对象适用于控制一个仅支持有限个用户的共享资源，是一种不需要使用忙碌等待（busy waiting）的方法。

如果信号量是一个任意的整数，通常被称为计数信号量（Counting semaphore），或一般信号量（general semaphore）；
如果信号量只有二进制的0或1，称为二进制信号量（binary semaphore）。在linux系统中，二进制信号量（binary semaphore）又称互斥锁（Mutex）

我们一般用信号量保护一组资源，比如数据库连接池、一组客户端的连接等等。
**每次获取资源时都会将信号量中的计数器减去对应的数值，在释放资源时重新加回来。当信号量没资源时尝试获取信号量的线程就会进入休眠，等待其他线程释放信号量。
如果信号量是只有0和1的二进位信号量，那么，它的 P/V 就和互斥锁的 Lock/Unlock 就一样了



## 源码分析
Go 内部使用信号量来控制goroutine的阻塞和唤醒，比如互斥锁sync.Mutex结构体定义的第二个字段就是一个信号量
```go
type Mutex struct {
    state int32
    sema  uint32
}
```
信号量的 PV 操作在Go内部是通过下面这几个底层函数实现的
```go
func runtime_Semacquire(s *uint32)
func runtime_SemacquireMutex(s *uint32, lifo bool, skipframes int)
func runtime_Semrelease(s *uint32, handoff bool, skipframes int)

```
这几个函数就是信号量的PV操作，不过他们都是给Go内部使用的.
通过 go:linkname 链接的对应实现的函数：
```go
// go1.21.5/src/runtime/sema.go

//go:linkname poll_runtime_Semacquire internal/poll.runtime_Semacquire
func poll_runtime_Semacquire(addr *uint32) {
	semacquire1(addr, false, semaBlockProfile, 0, waitReasonSemacquire)
}

//go:linkname sync_runtime_Semrelease sync.runtime_Semrelease
func sync_runtime_Semrelease(addr *uint32, handoff bool, skipframes int) {
	semrelease1(addr, handoff, skipframes)
}

//go:linkname sync_runtime_SemacquireMutex sync.runtime_SemacquireMutex
func sync_runtime_SemacquireMutex(addr *uint32, lifo bool, skipframes int) {
	semacquire1(addr, lifo, semaBlockProfile|semaMutexProfile, skipframes, waitReasonSyncMutexLock)
}

```

可以看到实际实现的函数为 semacquire1 与 semrelease1


如果想使用信号量，那就可以使用官方的扩展包golang.org/x/sync：Semaphore，这是一个带权重的信号量

扩展包源码：
```go
type Weighted struct {
	size    int64 // 总容量大小
	cur     int64 // 已经占用容量大小
	mu      sync.Mutex  // 提供临界区保护
	waiters list.List // 阻塞等待的调用者列表
}

// NewWeighted为并发访问创建一个新的加权信号量，该信号量具有给定的最大权值。
func NewWeighted(n int64) *Weighted {
	w := &Weighted{size: n}
	return w
}
//主要方法
func (s *Weighted) Acquire(ctx context.Context, n int64) error
func (s *Weighted) Release(n int64)

//等待者
type waiter struct {
	n     int64 /// 等待调用者权重值
	ready chan<- struct{} // 利用channel的close机制实现唤醒
}
```

semaphore库核心结构就是Weighted，主要有4个字段：
- size：这个代表的是最大权值，在创建Weighted对象指定
- cur：相当于一个游标，来记录当前已使用的权值
- mu：互斥锁，并发情况下做临界区保护
- waiters：阻塞等待的调用者列表，使用链表数据结构保证先进先出的顺序，存储的数据是waiter对象

### Acquire请求信号量资源
Acquire方法会监控资源是否可用，而且还要检测传递进来的context.Context对象是否发送了超时过期或者取消的信号
```go
func (s *Weighted) Acquire(ctx context.Context, n int64) error {
	s.mu.Lock() // 加锁保护临界区
	// 有资源可用并且没有等待获取权值的goroutine
	if s.size-s.cur >= n && s.waiters.Len() == 0 {
		s.cur += n // 加权
		s.mu.Unlock() // 释放锁
		return nil
	}
	// 要获取的权值n大于最大的权值了
	if n > s.size {
		// 先释放锁，确保其他goroutine调用Acquire的地方不被阻塞
		s.mu.Unlock()
		// 阻塞等待context的返回
		<-ctx.Done()
		return ctx.Err()
	}
	// 走到这里就说明现在没有资源可用了
	// 创建一个channel用来做通知唤醒
	ready := make(chan struct{})
	// 创建waiter对象
	w := waiter{n: n, ready: ready}
	// waiter按顺序入队
	elem := s.waiters.PushBack(w)
	// 释放锁，等待唤醒，别阻塞其他goroutine
	s.mu.Unlock()

	// 阻塞等待唤醒
	select {
	// context关闭
	case <-ctx.Done():
		err := ctx.Err() // 先获取context的错误信息
		s.mu.Lock()
		select {
		case <-ready:
			// 在context被关闭后被唤醒了，那么试图修复队列，假装我们没有取消
			err = nil
		default:
			// 判断是否是第一个元素
			isFront := s.waiters.Front() == elem
			// 移除第一个元素
			s.waiters.Remove(elem)
			// 如果是第一个元素且有资源可用通知其他waiter
			if isFront && s.size > s.cur {
				s.notifyWaiters()
			}
		}
		s.mu.Unlock()
		return err
	// 被唤醒了
	case <-ready:
		return nil
	}
}

```

### Release归还信号量资源  

```func (s *Weighted) Release(n int64) {
    s.mu.Lock()
    s.cur -= n
    if s.cur < 0 {
      s.mu.Unlock()
      panic("semaphore: released more than held")
    }
    s.notifyWaiters()
    s.mu.Unlock()
}

```

## 参考资料
1. [信号量的使用方法和其实现原理](https://juejin.cn/post/6906677772479889422#heading-5)



