<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [sync.Pool](#syncpool)
  - [背景](#%E8%83%8C%E6%99%AF)
    - [缺点](#%E7%BC%BA%E7%82%B9)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [Pool结构体](#pool%E7%BB%93%E6%9E%84%E4%BD%93)
    - [poolChain](#poolchain)
    - [poolDequeue](#pooldequeue)
    - [对象的清理注册](#%E5%AF%B9%E8%B1%A1%E7%9A%84%E6%B8%85%E7%90%86%E6%B3%A8%E5%86%8C)
    - [sync.Pool的 Get 函数](#syncpool%E7%9A%84-get-%E5%87%BD%E6%95%B0)
    - [sync.Pool的 Put 函数](#syncpool%E7%9A%84-put-%E5%87%BD%E6%95%B0)
  - [常见问题](#%E5%B8%B8%E8%A7%81%E9%97%AE%E9%A2%98)
    - [1. Pool 的内容会清理？清理会造成数据丢失吗？](#1-pool-%E7%9A%84%E5%86%85%E5%AE%B9%E4%BC%9A%E6%B8%85%E7%90%86%E6%B8%85%E7%90%86%E4%BC%9A%E9%80%A0%E6%88%90%E6%95%B0%E6%8D%AE%E4%B8%A2%E5%A4%B1%E5%90%97)
    - [2. runtime.GOMAXPROCS 与 pool 之间的关系？](#2-runtimegomaxprocs-%E4%B8%8E-pool-%E4%B9%8B%E9%97%B4%E7%9A%84%E5%85%B3%E7%B3%BB)
    - [3. New() 的作用？假如没有 New 会出现什么情况？](#3-new-%E7%9A%84%E4%BD%9C%E7%94%A8%E5%81%87%E5%A6%82%E6%B2%A1%E6%9C%89-new-%E4%BC%9A%E5%87%BA%E7%8E%B0%E4%BB%80%E4%B9%88%E6%83%85%E5%86%B5)
    - [4. 先 Put，再 Get 会出现什么情况？](#4-%E5%85%88-put%E5%86%8D-get-%E4%BC%9A%E5%87%BA%E7%8E%B0%E4%BB%80%E4%B9%88%E6%83%85%E5%86%B5)
    - [5. 只 Get 不 Put 会内存泄露吗？](#5-%E5%8F%AA-get-%E4%B8%8D-put-%E4%BC%9A%E5%86%85%E5%AD%98%E6%B3%84%E9%9C%B2%E5%90%97)
    - [6. 为什么要禁止 copy sync.Pool 实例？](#6-%E4%B8%BA%E4%BB%80%E4%B9%88%E8%A6%81%E7%A6%81%E6%AD%A2-copy-syncpool-%E5%AE%9E%E4%BE%8B)
  - [优秀应用实践](#%E4%BC%98%E7%A7%80%E5%BA%94%E7%94%A8%E5%AE%9E%E8%B7%B5)
    - [1. 官方包fmt源码分析](#1-%E5%AE%98%E6%96%B9%E5%8C%85fmt%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [2. 第三方库应用（gin)](#2-%E7%AC%AC%E4%B8%89%E6%96%B9%E5%BA%93%E5%BA%94%E7%94%A8gin)
  - [参考资料](#%E5%8F%82%E8%80%83%E8%B5%84%E6%96%99)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# sync.Pool

定位不是做类似连接池的东西，它的用途仅仅是增加对象重用的几率，减少gc的负担.

sync.Pool 是一个内存池。通常内存池是用来防止内存泄露的（例如C/C++)。sync.Pool 这个内存池却不是干这个的，
带 GC 功能的语言都存在垃圾回收 STW 问题，需要回收的内存块越多，STW 持续时间就越长。如果能让 new 出来的变量，一直不被回收，得到重复利用，是不是就减轻了 GC 的压力


sync.Pool中就是使用的PoolChain来实现的，它是一个单生产者多消费者的队列，可以同时有多个消费者消费数据，但是只有一个生产者生产数据

## 背景

Go是自动垃圾回收的(garbage collector)，这大大减少了程序编程负担。但gc是一把双刃剑，带来了编程的方便但同时也增加了运行时开销，
使用不当甚至会严重影响程序的性能。因此性能要求高的场景不能任意产生太多的垃圾（有gc但又不能完全依赖它挺恶心的），如何解决呢？

那就是要重用对象了，我们可以简单的使用一个chan把这些可重用的对象缓存起来，但如果很多goroutine竞争一个chan性能肯定是问题.

### 缺点

上面我们可以看到pool创建的时候是不能指定大小的，所有sync.Pool的缓存对象数量是没有限制的（只受限于内存），
因此使用sync.pool是没办法做到控制缓存对象数量的个数的。另外sync.pool缓存对象的期限是很诡异的，这是很多人错误理解的地方，

## 源码分析

![](.pool_images/sync_pool_structure.png)

### Pool结构体

![](.pool_images/pool_structure.png)

![](.pool_images/sync_pool_structure.png)
一个goroutine固定在一个局部调度器P上，从当前 P 对应的 poolLocal 取值， 若取不到，则从对应的 shared 数组上取，若还是取不到；
则尝试从其他 P 的 shared 中偷。 若偷不到，则调用 New 创建一个新的对象。池中所有临时对象在一次 GC 后会被全部清空。

```go
type Pool struct {
	noCopy noCopy

	// 每个 P 的本地队列，实际类型为 [P]poolLocal
	local     unsafe.Pointer // local fixed-size per-P pool, actual type is [P]poolLocal

	// [P]poolLocal的大小
	localSize uintptr        // size of the local array

	// victim 和 victimSize 在 (GC)poolCleanup 流程里赋值为 local 和 localSize
	victim     unsafe.Pointer // local from previous cycle
	victimSize uintptr        // size of victims array

    // 自定义的对象创建回调函数，当 pool 中无可用对象时会调用此函数.
	New func() interface{}
}
```

local 字段存储指向 [P]poolLocal 数组（严格来说，它是一个切片）的指针，localSize 则表示 local 数组的大小。
访问时，P 的 id 对应 [P]poolLocal 下标索引。通过这样的设计，多个 goroutine 使用同一个 Pool 时，减少了竞争，提升了性能。


victim 的机制用于减少 GC 后冷启动导致的性能抖动，让分配对象更平滑.

Victim Cache 本来是计算机架构里面的一个概念，是 CPU 硬件处理缓存的一种技术，sync.Pool 引入的意图在于降低 GC 压力的同时提高命中率.
![](.pool_images/poolDequeueTrait.png)
```go
// Local per-P Pool appendix.
type poolLocalInternal struct {
    // P 的私有缓存区，使用时无需要加锁
	private interface{} // Can be used only by the respective P.

    // 公共缓存区, 本地 P 可以 pushHead/popHead；其他 P popTail
	shared  poolChain   // Local P can pushHead/popHead; any P can popTail.
}

type poolLocal struct {
	poolLocalInternal

	// 将 poolLocal 补齐至两个缓存行的倍数，防止 false sharing,
	// 每个缓存行具有 64 bytes，即 512 bit
	// 目前我们的处理器一般拥有 32 * 1024 / 64 = 512 条缓存行
	// 伪共享，仅占位用，防止在 cache line 上分配多个 poolLocalInternal
	pad [128 - unsafe.Sizeof(poolLocalInternal{})%128]byte
}
```

sync.Pool 采用了一种类似 Go 运行时调度的机制，针对每个 p 有一个 private 的数据，同时还有一个 shared 的数据，如果在本地 private、shared 中没有数据，就去其他 P 对应的 shared 去偷取。
可能有多个 P 偷取同一个 shared, 这是多消费者



pad 涉及伪共享

> 现代 cpu 中，cache 都划分成以 cache line (cache block) 为单位，在 x86_64 体系下一般都是 64 字节，cache line 是操作的最小单元。

> 程序即使只想读内存中的 1 个字节数据，也要同时把附近 63 节字加载到 cache 中，如果读取超个 64 字节，那么就要加载到多个 cache line 中。

> 简单来说，如果没有 pad 字段，那么当需要访问 0 号索引的 poolLocal 时，CPU 同时会把 0 号和 1 号索引同时加载到 cpu cache。
> 在只修改 0 号索引的情况下，会让 1 号索引的 poolLocal 失效。这样，当其他线程想要读取 1 号索引时，发生 cache miss，还得重新再加载，对性能有损。
> 增加一个 pad，补齐缓存行，让相关的字段能独立地加载到缓存行就不会出现 false sharding 了。

### poolChain 

PoolChain 是在 PoolDequeue 的基础上实现的一个动态尺寸的队列，它的实现和 PoolDequeue 类似，只是增加了一个 headTail 的链表，用于存储多个 PoolDequeue
```go
// poolChain 是一个双端队列的实现

type poolChain struct {
    // 只有生产者会 push to，不用加锁
	head *poolChainElt

	// tail 是消费者用来pop的 poolDequeue。消费者访问，所以需要原子操作
	tail *poolChainElt
}

type poolChainElt struct {
	poolDequeue

	// next 被 producer 写，consumer 读。所以只会从 nil 变成 non-nil
	// prev 被 consumer 写，producer 读。所以只会从 non-nil 变成 nil
	next, prev *poolChainElt
}
```

整体的思想就是将多个poolDequeue串联起来，生产者在head处增加数据，消费者在tail处消费数据，当tail的poolDequeue为空时，就从head处获取一个poolDequeue。 当head满了的时候，就增加一个新的poolDequeue。 这样就实现了动态尺寸的队





### poolDequeue
poolDequeue 被实现为单生产者、多消费者的固定大小的无锁（atomic 实现） Ring 式队列（底层存储使用数组，使用两个指针标记 head、tail）。
生产者可以从 head 插入、head 删除，而消费者仅可从 tail 删除。

```go
// poolDequeue 是一个固定尺寸，使用 ringbuffer (环形队列) 方式实现的队列
type poolDequeue struct {
    // headTail 包含一个 32 位的 head 和一个 32 位的 tail 指针。这两个值都和 len(vals)-1 取模过。
    // tail 是队列中最老的数据，head 指向下一个将要填充的 slot
    // slots 的有效范围是 [tail, head)，由 consumers 持有。
	headTail uint64

	// vals 是一个存储 interface{} 的环形队列，它的 size 必须是 2 的幂
	// 如果 slot 为空，则 vals[i].typ 为空；否则，非空。
	// 一个 slot 在这时宣告无效：tail 不指向它了，vals[i].typ 为 nil
	// 由 consumer 设置成 nil，由 producer 读
	vals []eface
}
```

为什么headTail 变量将 head 和 tail 打包在了一起？

是为了实现lock free。对于一个 poolDequeue 来说，可能会被多个 P 同时访问就会出现并发问题。




两个重要的字段：
```go
const dequeueBits = 32

// 实现了pack和unpack方法，用于将 head 和 tail 打包到一个 uint64 中，或者从 uint64 中解包出 head 和 tail
func (d *poolDequeue) unpack(ptrs uint64) (head, tail uint32) {
	const mask = 1<<dequeueBits - 1
	head = uint32((ptrs >> dequeueBits) & mask)
	tail = uint32(ptrs & mask)
	return
}
```
- headTail： 一个 atomic.Uint64 类型的字段，它的高 32 位是 head，低 32 位是 tail。head 是下一个要填充的位置，tail 是最老的数据的位置
- vals： 一个 eface 类型的切片，它是一个环形队列，大小必须是 2 的幂次方。

我们看到 Pool 并没有直接使用 poolDequeue，原因是它的大小是固定的，而 Pool 的大小是没有限制的。因此，在 poolDequeue 之上包装了一下，变成了一个 poolChainElt 的双向链表，可以动态增长


生产者可以使用下面的方法：

- pushHead: 在队列头部新增加一个数据。如果队列满了，增加失败
- popHead： 在队列头部弹出一个数据。生产者总是弹出新增加的数据，除非队列为空


消费者可以使用下面的一个方法：

- popTail: 从队尾处弹出一个数据，除非队列为空。所以消费者总是消费最老的数据，这也正好符合大部分的场景

生产者增加数据
```go
const dequeueBits = 32

func (d *poolDequeue) pushHead(val any) bool {
	ptrs := d.headTail.Load()
	head, tail := d.unpack(ptrs)
	if (tail+uint32(len(d.vals)))&(1<<dequeueBits-1) == head {
		// 队列满
		return false
	}
	slot := &d.vals[head&uint32(len(d.vals)-1)]

	// 检查 head slot 是否被 popTail 释放
	typ := atomic.LoadPointer(&slot.typ)
	if typ != nil {
		// 另一个 goroutine 正在清理 tail，所以队列还是满的
		return false
	}

	// 如果值为空，那么设置一个特殊值
	if val == nil {
		val = dequeueNil(nil)
	}
	*(*any)(unsafe.Pointer(slot)) = val

	// Increment head. This passes ownership of slot to popTail
	// and acts as a store barrier for writing the slot.
	d.headTail.Add(1 << dequeueBits)
	return true
}
```




消费者消费数据的逻辑
```go
// 出队 popTail（从队尾获取元素）
func (d *poolDequeue) popTail() (any, bool) {
	var slot *eface
	for {
		ptrs := d.headTail.Load()
		head, tail := d.unpack(ptrs)
		if tail == head {
			// Queue is empty.
			return nil, false
		}

		// Confirm head and tail (for our speculative check
		// above) and increment tail. If this succeeds, then
		// we own the slot at tail.
		ptrs2 := d.pack(head, tail+1)
		if d.headTail.CompareAndSwap(ptrs, ptrs2) {
            // 成功读取了一个 slot
			slot = &d.vals[tail&uint32(len(d.vals)-1)]
			break
		}
	}

	// We now own slot.
	val := *(*any)(unsafe.Pointer(slot))
	if val == dequeueNil(nil) { //如果本身就存储的nil
		val = nil
	}

	// 释放 slot，这样 pushHead 就可以继续写入这个 slot 了
	slot.val = nil
	atomic.StorePointer(&slot.typ, nil)
	// At this point pushHead owns the slot.

	return val, true
}
```


```go
// 出队 popHead（从头部获取元素）
func (d *poolDequeue) popHead() (any, bool) {
	var slot *eface
	for {
		ptrs := d.headTail.Load()
		head, tail := d.unpack(ptrs)
		if tail == head {
			// Queue is empty.
			return nil, false
		}

		// Confirm tail and decrement head. We do this before
		// reading the value to take back ownership of this
		// slot.
		head--
		ptrs2 := d.pack(head, tail)
		if d.headTail.CompareAndSwap(ptrs, ptrs2) {
			// We successfully took back slot.
			slot = &d.vals[head&uint32(len(d.vals)-1)]
			break
		}
	}

	val := *(*any)(unsafe.Pointer(slot))
	if val == dequeueNil(nil) {
		val = nil
	}
	// Zero the slot. Unlike popTail, this isn't racing with
	// pushHead, so we don't need to be careful here.
	*slot = eface{}
	return val, true
}
```


### 对象的清理注册

对于 Pool 而言，并不能无限扩展，否则对象占用内存太多了，会引起内存溢出。

```go
// go1.23.0/src/sync/pool.go
func init() {
	runtime_registerPoolCleanup(poolCleanup)
}

func poolCleanup() {
	// This function is called with the world stopped, at the beginning of a garbage collection.
	// It must not allocate and probably should not call any runtime functions.

	// Because the world is stopped, no pool user can be in a
	// pinned section (in effect, this has all Ps pinned).

	// 清空oldPools中 victim 的对象
	for _, p := range oldPools {
		p.victim = nil
		p.victimSize = 0
	}

	// 将allPools对象池中，local对象迁移到 victim上。
	for _, p := range allPools {
		p.victim = p.local
		p.victimSize = p.localSize
		p.local = nil
		p.localSize = 0
	}

	// 将allPools迁移到oldPools，并清空allPools
	oldPools, allPools = allPools, nil
}
```

```go
// go1.23.0/src/runtime/mgc.go

//go:linkname sync_runtime_registerPoolCleanup sync.runtime_registerPoolCleanup
func sync_runtime_registerPoolCleanup(f func()) {
	poolcleanup = f
}

func clearpools() {
	// clear sync.Pools
	if poolcleanup != nil {
		poolcleanup()
	}
	// ...
}

func gcStart(trigger gcTrigger) {
	// ...
	
    // clearpools before we start the GC. If we wait the memory will not be
    // reclaimed until the next GC cycle.
    clearpools()
	
	// ..
}
```
在一轮 GC 到来时，victim 和 victimSize 会分别“接管” local 和 localSize。

可以看到pool包在init的时候注册了一个poolCleanup函数，它会清除所有的pool里面的所有缓存的对象，该函数注册进去之后会在每次gc之前都会调用，
因此sync.Pool缓存的期限只是两次gc之间这段时间.

正因为这样，我们是不可以使用sync.Pool去实现一个socket连接池的。



### sync.Pool的 Get 函数

![](.pool_images/pool_get.png)

```go
func (p *Pool) Get() any {
    // ..
	l, pid := p.pin()
	x := l.private
	l.private = nil
	if x == nil {
		// Try to pop the head of the local shard. We prefer
		// the head over the tail for temporal locality of
		// reuse.
		x, _ = l.shared.popHead()
		if x == nil {
			x = p.getSlow(pid)
		}
	}
	runtime_procUnpin()
    // ...
	if x == nil && p.New != nil {
		x = p.New()
	}
	return x
}
```

```go
func (c *poolChain) popHead() (any, bool) {
	d := c.head
	for d != nil {
		if val, ok := d.popHead(); ok {
			// 从 head 位置获取对象，如果该环形队列中还有数据则会返回 true；
			return val, ok
		}
		
		// 如果 head 位置的环形队列空了，会定位到 prev 节点继续尝试获取对象

		d = d.prev.Load()
	}
	return nil, false
}

```

### sync.Pool的 Put 函数

![](.pool_images/pool_put.png)

```go
// go1.23.0/src/sync/pool.go

// Put adds x to the pool.
func (p *Pool) Put(x any) {
    // ...
	l, _ := p.pin()
	if l.private == nil {
		l.private = x
	} else {
		l.shared.pushHead(x)
	}
	runtime_procUnpin()
    // ..
}
```

```go

func (c *poolChain) pushHead(val any) {
	d := c.head
	if d == nil {
		// 如果c.head为空，初始化链表.
		const initSize = 8 // Must be a power of 2
		d = new(poolChainElt)
		d.vals = make([]eface, initSize)
		c.head = d
		c.tail.Store(d)
	}

	// 将对象放入head中的环形队列poolDequeue
	if d.pushHead(val) { // 调用 poolDequeue 的 pushHead
		return
	}

	// 当poolDequeue满了，则新建一个双倍容量的链表节点，环形队列最大容量为 (1<<32)/4 =1073741824。
	newSize := len(d.vals) * 2
	if newSize >= dequeueLimit {
		// Can't make it any bigger.
		newSize = dequeueLimit
	}

	d2 := &poolChainElt{}
	d2.prev.Store(d)
	d2.vals = make([]eface, newSize)
	c.head = d2 //新建的节点放入head位置
	d.next.Store(d2)
	// 调用 poolDequeue 的 pushHead
	d2.pushHead(val)
}

```



## 常见问题

### 1. Pool 的内容会清理？清理会造成数据丢失吗？

Go 会在每个 GC 周期内定期清理 sync.Pool 内的数据。

要分几个方面来说这个问题。

已经从 sync.Pool Get 的值，在 poolClean 时虽说将 pool.local 置成了nil，Get 到的值依然是有效的，是被 GC 标记为黑色的，不会被 GC回收，当 Put 后又重新加入到 sync.Pool 中
在第一个 GC 周期内 Put 到 sync.Pool 的数值，在第二个 GC 周期没有被 Get 使用，就会被放在 local.victim 中。如果在 第三个 GC 周期仍然没有被使用就会被 GC 回收。

### 2. runtime.GOMAXPROCS 与 pool 之间的关系？

```go
func (p *Pool) pinSlow() (*poolLocal, int) {
	// Retry under the mutex.
	// Can not lock the mutex while pinned.
	runtime_procUnpin()
	allPoolsMu.Lock()
	defer allPoolsMu.Unlock()
	pid := runtime_procPin()
	// poolCleanup won't be called while we are pinned.
	s := p.localSize
	l := p.local
	if uintptr(pid) < s {
		return indexLocal(l, pid), pid
	}
	if p.local == nil {
		allPools = append(allPools, p)
	}
	// If GOMAXPROCS changes between GCs, we re-allocate the array and lose the old one.
	// 获取当前最大的 p 的数量
	size := runtime.GOMAXPROCS(0)
	local := make([]poolLocal, size)
	atomic.StorePointer(&p.local, unsafe.Pointer(&local[0])) // store-release
	runtime_StoreReluintptr(&p.localSize, uintptr(size))     // store-release
	return &local[pid], pid
}
```

### 3. New() 的作用？假如没有 New 会出现什么情况？

从上面的 pool.Get 流程图可以看出来，从 sync.Pool 获取一个内存会尝试从当前 private，shared，其他的 p 的 shared 获取或者 victim 获取，如果实在获取不到时，才会调用 New 函数来获取。也就是 New() 函数才是真正开辟内存空间的。New() 开辟出来的的内存空间使用完毕后，调用 pool.Put 函数放入到 sync.Pool 中被重复利用。

如果 New 函数没有被初始化会怎样呢？很明显，sync.Pool 就废掉了，因为没有了初始化内存的地方了。

### 4. 先 Put，再 Get 会出现什么情况？

```go
func main(){
    pool:= sync.Pool{
        New: func() interface{} {
            return item{}
        },
    }
    pool.Put(item{value:1})
    data := pool.Get()
    fmt.Println(data)
}
```

如果你直接跑这个例子，能得到你想像的结果，但是在某些情况下就不是这个结果了。

在 Pool.Get 注释里面有这么一句话：“Callers should not assume any relation between values passed to Put and the values returned by Get.”，告诉我们不能把值 Pool.Put 到 sync.Pool 中，再使用 Pool.Get 取出来，因为 sync.Pool 不是 map 或者 slice，放入的值是有可能拿不到的，sync.Pool 的数据结构就不支持做这个事情。

前面说使用 sync.Pool 容易被错误示例误导，就是上面这个写法。为什么 Put 的值 再 Get 会出现问题？

- 情况1：sync.Pool 的 poolCleanup 函数在系统 GC 时会被调用，Put 到 sync.Pool 的值，由于有可能一直得不到利用，被在某个 GC 周期内就有可能被释放掉了。
- 情况2：不同的 goroutine 绑定的 p 有可能是不一样的，当前 p 对应的 goroutine 放入到 sync.Pool 的值有可能被其他的 p 对应的 goroutine 取到，导致当前 goroutine 再也取不到这个值。
- 情况3：使用 runtime.GOMAXPROCS（N) 来改变 p 的数量，会使 sync.Pool 的 pool.poolLocal 释放重新开辟新的空间，导致 sync.Pool 被释放掉。
- 情况4：还有很多情况

### 5. 只 Get 不 Put 会内存泄露吗？

使用其他的池，如连接池，如果取连接使用后不放回连接池，就会出现连接池泄露，「是不是 sync.Pool 也有这个问题呢？」

通过上面的流程图，可以看出来 Pool.Get 的时候会尝试从当前 private，shared，其他的 p 的 shared 获取或者 victim 获取，如果实在获取不到时，才会调用 New 函数来获取，New 出来的内容本身还是受系统 GC 来控制的。所以如果我们提供的 New 实现不存在内存泄露的话，那么 sync.Pool 是不会内存泄露的。当 New 出来的变量如果不再被使用，就会被系统 GC 给回收掉。

如果不 Put 回 sync.Pool，会造成 Get 的时候每次都调用的 New 来从堆栈申请空间，达不到减轻 GC 压力。



### 6. 为什么要禁止 copy sync.Pool 实例？
因为 copy 后，对于同一个 Pool 实例中的 cache 对象，就有了两个指向来源。
原 Pool 清空之后，copy 的 Pool 没有清理掉，那么里面的对象就全都泄露了。并且 Pool 的无锁设计的基础是多个 Goroutine 不会操作到同一个数据结构，Pool 拷贝之后则不能保证这点（因为存储的成员都是指针）。
## 优秀应用实践

### 1. 官方包fmt源码分析

```go
func Printf(format string, a ...interface{}) (n int, err error) {
	return Fprintf(os.Stdout, format, a...)
}

func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
    p := newPrinter()
    p.doPrintf(format, a)
    n, err = w.Write(p.buf)
    p.free()
    return
}
```

获取对象

```go
// newPrinter allocates a new pp struct or grabs a cached one.
func newPrinter() *pp {
    p := ppFree.Get().(*pp)
    p.panicking = false
    p.erroring = false
    p.wrapErrs = false
    p.fmt.init(&p.buf)
    return p
}

var ppFree = sync.Pool{
    New: func() interface{} { return new(pp) },
}
```

归还

```go
func (p *pp) free() {
  if cap(p.buf) > 64<<10 {
    return
  }

  p.buf = p.buf[:0]
  p.arg = nil
  p.value = reflect.Value{}
  p.wrappedErr = nil
  ppFree.Put(p)
}
```

### 2. 第三方库应用（gin)

/Users/python/go/pkg/mod/github.com/gin-gonic/gin@v1.7.7/gin.go

```go
// ServeHTTP conforms to the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := engine.pool.Get().(*Context)
	c.writermem.reset(w)
	c.Request = req
	c.reset()

	engine.handleHTTPRequest(c)

	engine.pool.Put(c)
}
```

## 参考资料

- [lock-free、高性能的单生产者多消费者的队列：PoolDequeue 和 PoolChain](https://mp.weixin.qq.com/s/fj87oGZPkFKQiGZxhrYRVQ)
- [深入理解Golang的sync.Pool原理](https://cloud.tencent.com/developer/article/2217768)

