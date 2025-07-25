<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [waitGroup 同步等待组对象](#waitgroup-%E5%90%8C%E6%AD%A5%E7%AD%89%E5%BE%85%E7%BB%84%E5%AF%B9%E8%B1%A1)
  - [使用场景](#%E4%BD%BF%E7%94%A8%E5%9C%BA%E6%99%AF)
  - [源码实现](#%E6%BA%90%E7%A0%81%E5%AE%9E%E7%8E%B0)
    - [结构体](#%E7%BB%93%E6%9E%84%E4%BD%93)
    - [1 Add()](#1-add)
  - [2 Done()](#2-done)
  - [3 Wait()](#3-wait)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# waitGroup 同步等待组对象

## 使用场景

批量发出 RPC 或者 HTTP 请求.

![](.waitGroup_images/waitGroup_info.png)
sync.WaitGroup主要用于等待一组goroutine退出，本质上其实就是一个计数器，我们可以通过Add指定我们需要等待退出的goroutine的数量，然后通过Done来递减，如果为0,则可以退出


## 源码实现

### 结构体

```go
// go1.20/src/sync/waitgroup.go
type WaitGroup struct {
	noCopy noCopy // 避免不小心复制了一个不应该复制的对象

	state atomic.Uint64 // high 32 bits are counter, low 32 bits are waiter count.
	sema  uint32
}
```
- noCopy: 
- state 代表两个字段：counter，waiter 
  - counter ：可以理解为一个计数器，计算经过 wg.Add(N), wg.Done() 后的值。
  - waiter ：当前等待 WaitGroup 任务结束的等待者数量。其实就是调用 wg.Wait() 的次数，所以通常这个值是 1 。
- sema ：信号量，用来唤醒 Wait() 函数。


### 1 Add()

![](.waitGroup_images/waitGroup_add.png)

Add()不止是简单的将信号量增 delta，还需要考虑很多因素

- 内部运行计数不能为负
- Add 必须与 Wait 属于 happens before 关系
  - 毕竟 Wait 是同步屏障，没有 Add，Wait 就没有了意义
- 通过信号量通知所有正在等待的 goroutine


为什么要将 counter 和 waiter 放在一起 ？

当同时发现 wg.counter <= 0 && wg.waiter != 0 时，才会去唤醒等待的 waiters，让等待的协程继续运行。
但是使用 WaitGroup 的调用方一般都是并发操作，如果不同时获取的 counter 和 waiter 的话，就会造成获取到的 counter 和 waiter 可能不匹配，造成程序 deadlock 或者程序提前结束等待

```go
// sync/waitgroup.go
func (wg *WaitGroup) Add(delta int) {
	// 获取当前计数
	statep, semap := wg.state()
	if race.Enabled {
		_ = *statep // trigger nil deref early
		if delta < 0 {
			// Synchronize decrements with Wait.
			race.ReleaseMerge(unsafe.Pointer(wg))
		}
		race.Disable()
		defer race.Enable()
	}
	// 使用高32位进行counter计数
	state := atomic.AddUint64(statep, uint64(delta)<<32)
    // 协程数量运行计数
	v := int32(state >> 32)
	// 获取低32位即waiter等待计数
	w := uint32(state)
	if race.Enabled && delta > 0 && v == int32(delta) {
		// The first increment must be synchronized with Wait.
		// Need to model this as a read, because there can be
		// several concurrent wg.counter transitions from 0.
		race.Read(unsafe.Pointer(semap))
	}
	if v < 0 {
		// 任务计数器不能为负数
		panic("sync: negative WaitGroup counter")
	}

	if w != 0 && delta > 0 && v == int32(delta) {
      // 添加与等待同时调用（应该是happens before的关系）
      // 已经执行了Wait，不容许再执行Add
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}

	if v > 0 || w == 0 {
        // 如果当前v>0,协程数量大于0，说明还有协程在运行，不需要唤醒等待者；
        // 如果w==0,等待者数量为0，说明没有等待者，也不需要唤醒。
        // 以上两种情况直接返回即可
		return
	}
	// This goroutine has set counter to 0 when waiters > 0.
	// Now there can't be concurrent mutations of state:
	// - Adds must not happen concurrently with Wait,
	// - Wait does not increment waiters if it sees counter == 0.
	// Still do a cheap sanity check to detect WaitGroup misuse.
	
	// 当waiters > 0 的时候，并且当前v==0，这个时候如果检查发现state状态前后发生改变，则
	// 证明当前有人修改过，则删除
	// 如果走到这个地方则证明经过之前的操作后，当前的v==0,w!=0,就证明之前一轮的Done已经全部完成，现在需要唤醒所有在wait的goroutine
	// 此时如果发现当前的*statep值又发生了改变，则证明有有人进行了Add操作
	// 也就是这里的WaitGroup滥用
	if *statep != state { // 确保Add方法和Wait方法没有并发调用。如果有并发调用，panic。
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	// Reset waiters count to 0.
	// 将当前state的状态设置为0，就可以进行下次的重用了
	*statep = 0
	for ; w != 0; w-- {
		// 释放所有排队的waiter
		runtime_Semrelease(semap, false, 0)
	}
}
```
对于 wg.state 的状态变更，WaitGroup 的 Add()，Wait() 是使用 atomic 来做原子计算的(为了避免锁竞争)。但是由于 atomic 需要使用者保证其 64 位对齐，
所以将 counter 和 waiter 都设置成 uint32，同时作为一个变量，即满足了 atomic 的要求，同时也保证了获取 waiter 和 counter 的状态完整性。但这也就导致了 32位，64位机器上获取 state 的方式并不相同。

![](.waitGroup_images/waitgroup_in_32bit_n_64bit.png)

```go
// state returns pointers to the state and sema fields stored within wg.state1.
func (wg *WaitGroup) state() (statep *uint64, semap *uint32) {
	if uintptr(unsafe.Pointer(&wg.state1))%8 == 0 {
		return (*uint64)(unsafe.Pointer(&wg.state1)), &wg.state1[2]
	} else {
		return (*uint64)(unsafe.Pointer(&wg.state1[1])), &wg.state1[0]
	}
}
```

因为 64 位机器上本身就能保证 64 位对齐，所以按照 64 位对齐来取数据，拿到 state1[0], state1[1] 本身就是64 位对齐的。
但是 32 位机器上并不能保证 64 位对齐，因为 32 位机器是 4 字节对齐，如果也按照 64 位机器取 state[0]，state[1] 就有可能会造成 atomic 的使用错误。

于是 32 位机器上空出第一个 32 位，也就使后面 64 位天然满足 64 位对齐，第一个 32 位放入 sema 刚好合适。
早期 WaitGroup 的实现 sema 是和 state1 分开的，也就造成了使用 WaitGroup 就会造成 4 个字节浪费，不过 go1.11 之后就是现在的结构了。




## 2 Done()
```go

func (wg *WaitGroup) Done() {
    // 减去一个-1
    wg.Add(-1)
}
```


## 3 Wait()
Wait()主要用于阻塞 g，直到 WaitGroup 的计数为 0。

![](.waitGroup_images/waitGroup_wait.png)
```go
func (wg *WaitGroup) Wait() {
    statep, semap := wg.state()
    if race.Enabled {
        _ = *statep // trigger nil deref early
        race.Disable()
    }
    for {
        // 获取state的状态
        state := atomic.LoadUint64(statep)
        v := int32(state >> 32) // 获取高32位的count
        w := uint32(state) // 获取当前正在Wait的数量
        if v == 0 { // 如果当前v==0就直接return， 表示当前不需要等待
            // Counter is 0, no need to wait.
            if race.Enabled {
                race.Enable()
                race.Acquire(unsafe.Pointer(wg))
            }
            return
        }
        // CAS操作，将等待者数量加1。如果没成功，则下一次循环再次尝试，这也是有个无限for循环的原因
        if atomic.CompareAndSwapUint64(statep, state, state+1) {
            if race.Enabled && w == 0 {
                // Wait must be synchronized with the first Add.
                // Need to model this is as a write to race with the read in Add.
                // As a consequence, can do the write only for the first waiter,
                // otherwise concurrent Waits will race with each other.
                race.Write(unsafe.Pointer(semap))
            }
            // 如果成功则进行排队休眠等待唤醒
            runtime_Semacquire(semap)
            // 如果唤醒后发现state的状态不为0，则证明在唤醒的过程中WaitGroup又被重用，则panic
            if *statep != 0 {
                panic("sync: WaitGroup is reused before previous Wait has returned")
            }
            if race.Enabled {
                race.Enable()
                race.Acquire(unsafe.Pointer(wg))
            }
            return
        }
    }
}
```

## 参考
