# waitGroup同步等待组对象

## 使用场景
批量发出 RPC 或者 HTTP 请求

## 结构体

```go
type WaitGroup struct {
    noCopy noCopy
    state1 [3]uint32
}
```
WaitGroup.state1 其实代表三个字段：counter，waiter，sema。
- counter ：可以理解为一个计数器，计算经过 wg.Add(N), wg.Done() 后的值。
- waiter ：当前等待 WaitGroup 任务结束的等待者数量。其实就是调用 wg.Wait() 的次数，所以通常这个值是 1 。
- sema ：信号量，用来唤醒 Wait() 函数。


## 1. add()
![](.waitGroup_images/waitGroup_add.png)

### 为什么要将 counter 和 waiter 放在一起 ？

当同时发现 wg.counter <= 0 && wg.waiter != 0 时，才会去唤醒等待的 waiters，让等待的协程继续运行。
但是使用 WaitGroup 的调用方一般都是并发操作，如果不同时获取的 counter 和 waiter 的话，就会造成获取到的 counter 和 waiter 可能不匹配，造成程序 deadlock 或者程序提前结束等待

```go
// sync/waitgroup.go
func (wg *WaitGroup) Add(delta int) {
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
	state := atomic.AddUint64(statep, uint64(delta)<<32)
	v := int32(state >> 32)
	w := uint32(state)
	if race.Enabled && delta > 0 && v == int32(delta) {
		// The first increment must be synchronized with Wait.
		// Need to model this as a read, because there can be
		// several concurrent wg.counter transitions from 0.
		race.Read(unsafe.Pointer(semap))
	}
	if v < 0 {
		panic("sync: negative WaitGroup counter")
	}
	if w != 0 && delta > 0 && v == int32(delta) {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	if v > 0 || w == 0 {
		return
	}
	// This goroutine has set counter to 0 when waiters > 0.
	// Now there can't be concurrent mutations of state:
	// - Adds must not happen concurrently with Wait,
	// - Wait does not increment waiters if it sees counter == 0.
	// Still do a cheap sanity check to detect WaitGroup misuse.
	if *statep != state {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	// Reset waiters count to 0.
	*statep = 0
	for ; w != 0; w-- {
		runtime_Semrelease(semap, false, 0)
	}
}
```
对于 wg.state 的状态变更，WaitGroup 的 Add()，Wait() 是使用 atomic 来做原子计算的(为了避免锁竞争)。但是由于 atomic 需要使用者保证其 64 位对齐，
所以将 counter 和 waiter 都设置成 uint32，同时作为一个变量，即满足了 atomic 的要求，同时也保证了获取 waiter 和 counter 的状态完整性。但这也就导致了 32位，64位机器上获取 state 的方式并不相同。
![](.waitGroup_images/waitgroup_in_32bit_n_64bit.png)


因为 64 位机器上本身就能保证 64 位对齐，所以按照 64 位对齐来取数据，拿到 state1[0], state1[1] 本身就是64 位对齐的。但是 32 位机器上并不能保证 64 位对齐，因为 32 位机器是 4 字节对齐，如果也按照 64 位机器取 state[0]，state[1] 就有可能会造成 atmoic 的使用错误。

于是 32 位机器上空出第一个 32 位，也就使后面 64 位天然满足 64 位对齐，第一个 32 位放入 sema 刚好合适。早期 WaitGroup 的实现 sema 是和 state1 分开的，也就造成了使用 WaitGroup 就会造成 4 个字节浪费，不过 go1.11 之后就是现在的结构了。

## 2. wait()
![](.waitGroup_images/waitGroup_wait.png)