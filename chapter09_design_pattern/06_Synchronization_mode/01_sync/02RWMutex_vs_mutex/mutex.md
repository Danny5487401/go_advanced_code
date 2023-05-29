<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [mutex 互斥锁 和 RWMutex读写锁](#mutex-%E4%BA%92%E6%96%A5%E9%94%81-%E5%92%8C-rwmutex%E8%AF%BB%E5%86%99%E9%94%81)
  - [什么时候需要用到锁？](#%E4%BB%80%E4%B9%88%E6%97%B6%E5%80%99%E9%9C%80%E8%A6%81%E7%94%A8%E5%88%B0%E9%94%81)
  - [死锁](#%E6%AD%BB%E9%94%81)
    - [死锁产生的原因](#%E6%AD%BB%E9%94%81%E4%BA%A7%E7%94%9F%E7%9A%84%E5%8E%9F%E5%9B%A0)
  - [Mutex结构体](#mutex%E7%BB%93%E6%9E%84%E4%BD%93)
    - [互斥锁的状态](#%E4%BA%92%E6%96%A5%E9%94%81%E7%9A%84%E7%8A%B6%E6%80%81)
    - [lock加锁过程](#lock%E5%8A%A0%E9%94%81%E8%BF%87%E7%A8%8B)
      - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [解锁过程](#%E8%A7%A3%E9%94%81%E8%BF%87%E7%A8%8B)
    - [案例](#%E6%A1%88%E4%BE%8B)
      - [1. 一个goroutine](#1-%E4%B8%80%E4%B8%AAgoroutine)
      - [2. 两个goroutine](#2-%E4%B8%A4%E4%B8%AAgoroutine)
      - [3. 三个goroutine](#3-%E4%B8%89%E4%B8%AAgoroutine)
      - [4. 没有加锁，直接解锁问题-异常](#4-%E6%B2%A1%E6%9C%89%E5%8A%A0%E9%94%81%E7%9B%B4%E6%8E%A5%E8%A7%A3%E9%94%81%E9%97%AE%E9%A2%98-%E5%BC%82%E5%B8%B8)
  - [RWMutex结构体](#rwmutex%E7%BB%93%E6%9E%84%E4%BD%93)
    - [方法](#%E6%96%B9%E6%B3%95)
    - [读和写锁关系](#%E8%AF%BB%E5%92%8C%E5%86%99%E9%94%81%E5%85%B3%E7%B3%BB)
    - [写锁饥饿问题](#%E5%86%99%E9%94%81%E9%A5%A5%E9%A5%BF%E9%97%AE%E9%A2%98)
    - [写锁计数](#%E5%86%99%E9%94%81%E8%AE%A1%E6%95%B0)
    - [读锁加锁实现](#%E8%AF%BB%E9%94%81%E5%8A%A0%E9%94%81%E5%AE%9E%E7%8E%B0)
    - [读锁释放实现](#%E8%AF%BB%E9%94%81%E9%87%8A%E6%94%BE%E5%AE%9E%E7%8E%B0)
    - [写锁加锁实现](#%E5%86%99%E9%94%81%E5%8A%A0%E9%94%81%E5%AE%9E%E7%8E%B0)
    - [写锁释放实现](#%E5%86%99%E9%94%81%E9%87%8A%E6%94%BE%E5%AE%9E%E7%8E%B0)
    - [写锁与读锁的公平性](#%E5%86%99%E9%94%81%E4%B8%8E%E8%AF%BB%E9%94%81%E7%9A%84%E5%85%AC%E5%B9%B3%E6%80%A7)
    - [总结 读写互斥锁的实现](#%E6%80%BB%E7%BB%93-%E8%AF%BB%E5%86%99%E4%BA%92%E6%96%A5%E9%94%81%E7%9A%84%E5%AE%9E%E7%8E%B0)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# mutex 互斥锁 和 RWMutex读写锁

## 什么时候需要用到锁？

当程序中就一个线程的时候，是不需要加锁的，但是通常实际的代码不会只是单线程，所以这个时候就需要用到锁了，那么关于锁的使用场景主要涉及到哪些呢？
1. 多个线程在读相同的数据时
2. 多个线程在写相同的数据时
3. 同一个资源，有读又有写

## 死锁
多线程以及多进程改善了系统资源的利用率并提高了系统 的处理能力。然而，并发执行也带来了新的问题——死锁。
死锁是指两个或两个以上的进程（线程）在运行过程中因争夺资源而造成的一种僵局（Deadly-Embrace) ) ，若无外力作用，这些进程（线程）都将无法向前推进。

### 死锁产生的原因

- 竞争不可抢占资源引起死锁
通常系统中拥有的不可抢占资源，其数量不足以满足多个进程运行的需要，使得进程在运行过程中，会因争夺资源而陷入僵局，如磁带机、打印机等。只有对不可抢占资源的竞争 才可能产生死锁，对可抢占资源的竞争是不会引起死锁的。
- 竞争可消耗资源引起死锁
- 进程推进顺序不当引起死锁
进程在运行过程中，请求和释放资源的顺序不当，也同样会导致死锁。例如，并发进程 P1、P2分别保持了资源R1、R2，而进程P1申请资源R2，进程P2申请资源R1时，两者都会因为所需资源被占用而阻塞。
信号量使用不当也会造成死锁。进程间彼此相互等待对方发来的消息，结果也会使得这 些进程间无法继续向前推进。例如，进程A等待进程B发的消息，进程B又在等待进程A 发的消息，可以看出进程A和B不是因为竞争同一资源，而是在等待对方的资源导致死锁。



## Mutex结构体

```go
type Mutex struct {
    state int32  // 当前互斥锁的状态
    sema  uint32 // 信号
}
```
Mutex 的实现主要借助了 CAS 指令 + 自旋 + 信号量来实现

### 互斥锁的状态
在默认情况下，互斥锁的所有状态位都是 0，int32 中的不同位分别表示了不同的状态：
```go
const (
	mutexLocked = 1 << iota // 表示互斥锁的锁定状态；
	mutexWoken  // 表示从正常模式被从唤醒；
	mutexStarving  // 当前的互斥锁进入饥饿状态；
	mutexWaiterShift = iota

	// Mutex fairness.
	//
	// Mutex can be in 2 modes of operations: 正常和饥饿
	// 正常情况下，waiters 是先进先出FIFO, but a woken up waiter
	// does not own the mutex and competes with new arriving goroutines over
	// the ownership.新请求的 Goroutine 进入自旋时是仍然拥有 CPU 的, 所以比等待信号量唤醒的 Goroutine 更容易获取锁. 
	// 用官方话说就是，新请求锁的 Goroutine具有优势，它正在CPU上执行，而且可能有好几个，所以刚刚唤醒的 Goroutine 有很大可能在锁竞争中失败. 
	//  In such case it is queued at front of the wait queue. 
	// 那些等待超过 1 ms 的 Goroutine 还没有获取到锁，该 Goroutine 就会进入饥饿状态。
	//
	// 饥饿模式下，直接将唤醒信号发给第一个等待的 Goroutine.
	// New arriving goroutines don't try to acquire the mutex even if it appears
	// to be unlocked, and don't try to spin. Instead they queue themselves at
	// the tail of the wait queue.
	//
	// If a waiter receives ownership of the mutex and sees that either
	// (1) it is the last waiter in the queue, or (2) it waited for less than 1 ms,
	// it switches mutex back to normal operation mode.
	//
	// Normal mode has considerably better performance as a goroutine can acquire
	// a mutex several times in a row even if there are blocked waiters.
	// Starvation mode is important to prevent pathological cases of tail latency.
	starvationThresholdNs = 1e6
)

```

- waitersCount — 当前互斥锁上等待的 Goroutine 个数

Note：注意Mutex 状态（mutexLocked，mutexWoken，mutexStarving，mutexWaiterShift） 与 Goroutine 之间的状态（starving，awoke）改变


### lock加锁过程
![](.mutex_images/mutex_lock.png)

如果互斥锁的状态不是 0 时就会调用 sync.Mutex.lockSlow 尝试通过自旋（Spinning）等方式等待锁的释放，该方法的主体是一个非常大 for 循环，这里将它分成几个部分介绍获取锁的过程：
1. 判断当前 Goroutine 能否进入自旋；
2. 通过自旋等待互斥锁的释放；
3. 计算互斥锁的最新状态；
4. 更新互斥锁的状态并获取锁


#### 源码分析
```go
func (m *Mutex) Lock() {
	// Fast path: grab unlocked mutex.
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		if race.Enabled {
			race.Acquire(unsafe.Pointer(m))
		}
		return
	}
	// Slow path (outlined so that the fast path can be inlined)
	m.lockSlow()
}
```

```go
func (m *Mutex) lockSlow() {
	// 。。。
    queueLifo := waitStartTime != 0
    if waitStartTime == 0 {
        waitStartTime = runtime_nanotime()
    }
    // 入队
    runtime_SemacquireMutex(&m.sema, queueLifo, 1)
    
    // 。。。 
}
```

1. Goroutine 第一次被阻塞：
   由于 waitStartTime 等于 0，runtime_SemacquireMutex 的 queueLifo 等于 false, 于是该 Goroutine 放入到队列的尾部.

2. goroutine 被唤醒过，但是没加锁成功，再次被阻塞：由于 Goroutine 被唤醒过，waitStartTime 不等于 0，runtime_SemacquireMutex 的 queueLifo 等于 true, 于是该 Goroutine 放入到队列的头部

```go

func (m *Mutex) unlockSlow(new int32) {
	if (new+mutexLocked)&mutexLocked == 0 {
		throw("sync: unlock of unlocked mutex")
	}
	if new&mutexStarving == 0 {
		// 当前 mutex 不是饥饿状态：
		old := new
		for {
			// If there are no waiters or a goroutine has already
			// been woken or grabbed the lock, no need to wake anyone.
			// In starvation mode ownership is directly handed off from unlocking
			// goroutine to the next waiter. We are not part of this chain,
			// since we did not observe mutexStarving when we unlocked the mutex above.
			// So get off the way.
			if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
				return
			}
			// Grab the right to wake someone.
			new = (old - 1<<mutexWaiterShift) | mutexWoken
			if atomic.CompareAndSwapInt32(&m.state, old, new) {
				// 设置 runtime_Semrelease 的 handoff 参数是 false, 于是唤醒其中一个 Goroutine
				runtime_Semrelease(&m.sema, false, 1)
				return
			}
			old = m.state
		}
	} else {
		// 当前 mutex 已经是饥饿状态
		// Starving mode: handoff mutex ownership to the next waiter, and yield
		// our time slice so that the next waiter can start to run immediately.
		// Note: mutexLocked is not set, the waiter will set it after wakeup.
		// But mutex is still considered locked if mutexStarving is set,
		// so new coming goroutines won't acquire it.
		// 设置 runtime_Semrelease 的 handoff 参数是 true, 于是让等待队列头部的第一个 Goroutine 获得锁
		runtime_Semrelease(&m.sema, true, 1)
	}
}
```


### 解锁过程
![](.mutex_images/mutex_unlock.png)

1. 当互斥锁已经被解锁时，调用 sync.Mutex.Unlock 会直接抛出异常；
2. 当互斥锁处于饥饿模式时，将锁的所有权交给队列中的下一个等待者，等待者会负责设置 mutexLocked 标志位；
3. 当互斥锁处于普通模式时，如果没有 Goroutine 等待锁的释放或者已经有被唤醒的 Goroutine 获得了锁，会直接返回；在其他情况下会通过 sync.runtime_Semrelease 唤醒对应的 Goroutine

![](.mutex_images/lock_member.png)


### 案例
#### 1. 一个goroutine
![](.mutex_images/one_goroutine_lock.png)
#### 2. 两个goroutine
![](.mutex_images/two_gorountine_lock.png)
#### 3. 三个goroutine
![](.mutex_images/three_goroutine_lock.png)

#### 4. 没有加锁，直接解锁问题-异常
![](.mutex_images/unlock_again.png)

## RWMutex结构体
```go

type RWMutex struct {
    w           Mutex  // held if there are pending writers
    writerSem   uint32 // 用于writer等待读完成排队的信号量
    readerSem   uint32 // 用于reader等待写完成排队的信号量
    readerCount int32  // 读锁的计数器
    readerWait  int32  // 获取写锁时需要等待的写者的数量，用于防止写者饿死
}
```

在读多写少的环境中，可以优先使用读写互斥锁（sync.RWMutex），它比互斥锁更加高效。sync 包中的 RWMutex 提供了读写互斥锁的封装

分类:读锁和写锁
- 如果设置了一个写锁，那么其它读的线程以及写的线程都拿不到锁，这个时候，与互斥锁的功能相同
- 如果设置了一个读锁，那么其它写的线程是拿不到锁的，但是其它读的线程是可以拿到锁

```go
const rwmutexMaxReaders = 1 << 30 // 支持最多2^30个读锁
```
### 方法

写操作使用 sync.RWMutex.Lock 和 sync.RWMutex.Unlock 方法；


读操作使用 sync.RWMutex.RLock 和 sync.RWMutex.RUnlock 方法；

### 读和写锁关系
调用 sync.RWMutex.Lock 尝试获取写锁时；
1. 每次 sync.RWMutex.RUnlock 都会将 readerCount 其减一，当它归零时该 Goroutine 会获得写锁；
2. 将 readerCount 减少 rwmutexMaxReaders 个数以阻塞后续的读操作；

调用 sync.RWMutex.Unlock 释放写锁时，会先通知所有的读操作，然后才会释放持有的互斥锁

### 写锁饥饿问题
因为读锁是共享的，所以如果当前已经有读锁，那后续goroutine继续加读锁正常情况下是可以加锁成功，
但是如果一直有读锁进行加锁，那尝试加写锁的goroutine则可能会长期获取不到锁，这就是因为读锁而导致的写锁饥饿问题

### 写锁计数

读写锁中允许加读锁的最大数量是4294967296，在go里面对写锁的计数采用了负值进行，通过递减最大允许加读锁的数量从而进行写锁对读锁的抢占
```go
const rwmutexMaxReaders = 1 << 30
```

### 读锁加锁实现
![](.mutex_images/readerMutex_lock.png)

```go
func (rw *RWMutex) RLock() {
	// // 竞争检测代码，不看
    if race.Enabled {
        _ = rw.w.state
        race.Disable()
    }
    // 累加reader计数器，如果小于0则表明有writer正在等待
    if atomic.AddInt32(&rw.readerCount, 1) < 0 {
        // 当前有writer正在等待读锁，读锁就加入排队
        runtime_SemacquireMutex(&rw.readerSem, false)
    }
    if race.Enabled {
        race.Enable()
        race.Acquire(unsafe.Pointer(&rw.readerSem))
    }
}
```

### 读锁释放实现
![](.mutex_images/readerMutex_release.png)

```go
func (rw *RWMutex) RUnlock() {
    if race.Enabled {
        _ = rw.w.state
        race.ReleaseMerge(unsafe.Pointer(&rw.writerSem))
        race.Disable()
    }
    // 如果小于0，则表明当前有writer正在等待
    if r := atomic.AddInt32(&rw.readerCount, -1); r < 0 {
        if r+1 == 0 || r+1 == -rwmutexMaxReaders {
            race.Enable()
            throw("sync: RUnlock of unlocked RWMutex")
        }
        // 将等待reader的计数减1，证明当前是已经有一个读的，如果值==0，则进行唤醒等待的
        if atomic.AddInt32(&rw.readerWait, -1) == 0 {
            //当检查到有加写锁的情况下，就递减readerWait，并由最后一个释放reader lock的goroutine来实现唤醒写锁
            // The last reader unblocks the writer.
            runtime_Semrelease(&rw.writerSem, false)
        }
    }
    if race.Enabled {
        race.Enable()
    }
}
```

### 写锁加锁实现
![](.mutex_images/writerMutex_lock.png)
```go
func (rw *RWMutex) Lock() {
    if race.Enabled {
        _ = rw.w.state
        race.Disable()
    }
    // 首先获取mutex锁，同时多个goroutine只有一个可以进入到下面的逻辑
    rw.w.Lock()
    // 对readerCounter进行进行抢占，通过递减rwmutexMaxReaders允许最大读的数量
    // 来实现写锁对读锁的抢占
    r := atomic.AddInt32(&rw.readerCount, -rwmutexMaxReaders) + rwmutexMaxReaders
    // 记录需要等待多少个reader完成,如果发现不为0，则表明当前有reader正在读取，当前goroutine
    // 需要进行排队等待
    if r != 0 && atomic.AddInt32(&rw.readerWait, r) != 0 {
    	// // 写锁发现需要等待的读锁释放的数量不为0，就自己自己去休眠了
        runtime_SemacquireMutex(&rw.writerSem, false)
    }
    if race.Enabled {
        race.Enable()
        race.Acquire(unsafe.Pointer(&rw.readerSem))
        race.Acquire(unsafe.Pointer(&rw.writerSem))
    }
}

```

### 写锁释放实现
![](.mutex_images/writer_mutex_release.png)
```go
func (rw *RWMutex) Unlock() {
    if race.Enabled {
        _ = rw.w.state
        race.Release(unsafe.Pointer(&rw.readerSem))
        race.Disable()
    }

    // 将reader计数器复位，上面减去了一个rwmutexMaxReaders现在再重新加回去即可复位
    r := atomic.AddInt32(&rw.readerCount, rwmutexMaxReaders)
    if r >= rwmutexMaxReaders {
        race.Enable()
        throw("sync: Unlock of unlocked RWMutex")
    }
    // 唤醒所有的读锁
    for i := 0; i < int(r); i++ {
        runtime_Semrelease(&rw.readerSem, false)
    }
    // 释放mutex
    rw.w.Unlock()
    if race.Enabled {
        race.Enable()
    }
}
```

### 写锁与读锁的公平性

在加读锁和写锁的工程中都使用atomic.AddInt32来进行递增，而该指令在底层是会通过LOCK来进行CPU总线加锁的，

因此多个CPU同时执行readerCount其实只会有一个成功，从这上面看其实是写锁与读锁之间是相对公平的，
谁先达到谁先被CPU调度执行，进行LOCK锁cache line成功，谁就加成功锁


### 总结 读写互斥锁的实现

1. 读锁不能阻塞读锁，引入readerCount实现
2. 读锁需要阻塞写锁，直到所有读锁都释放，引入readerSem实现
3. 写锁需要阻塞读锁，直到所有写锁都释放，引入wirterSem实现
4. 写锁需要阻塞写锁，引入Metux实现