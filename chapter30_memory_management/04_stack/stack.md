# stack



在Go应用程序运行时，每个goroutine都维护着一个自己的栈区，这个栈区只能自己使用不能被其他goroutine使用。
栈区的初始大小是2KB（比x86_64架构下线程的默认栈2M要小很多），在goroutine运行的时候栈区会按照需要增长和收缩，占用的内存最大限制的默认值在64位系统上是1GB。

栈大小的初始值和上限这部分的设置都可以在Go的源码runtime/stack.go里找到：
```go
//Users/python/go/go1.18/src/runtime/stack.go

const (
    // StackSystem is a number of additional bytes to add
    // to each stack below the usual guard area for OS-specific
    // purposes like signal handling. Used on Windows, Plan 9,
    // and iOS because they do not use a separate stack.
    _StackSystem = goos.IsWindows*512*goarch.PtrSize + goos.IsPlan9*512 + goos.IsIos*goarch.IsArm64*1024
    
    // The minimum size of stack used by Go code
    _StackMin = 2048 // 2KB


)
var maxstacksize uintptr = 1 << 20 // 1GB,enough until runtime.main sets it for real

```

## 基本概念

### 栈分裂 stack-split
由于 Go 程序中的 goroutine 数目是不可确定的，并且实际场景可能会有百万级别的 goroutine，runtime 必须使用保守的思路来给 goroutine 分配空间以避免吃掉所有的可用内存。
也由于此，每个新的 goroutine 会被 runtime 分配初始为 2KB 大小的栈空间(Go 的栈在底层实际上是分配在堆空间上的)。
随着一个 goroutine 进行自己的工作，可能会超出最初分配的栈空间限制(就是栈溢出的意思)。

为了防止这种情况发生，runtime 确保 goroutine 在超出栈范围时，会创建一个相当于原来两倍大小的新栈，并将原来栈的上下文拷贝到新栈上。
这个过程被称为 栈分裂(stack-split)，这样使得 goroutine 栈能够动态调整大小.

为了使栈分裂正常工作，编译器会在每一个函数的开头和结束位置插入指令来防止 goroutine 爆栈。
有时候为了避免不必要的开销，一定不会爆栈的函数会被标记上 NOSPLIT 来提示编译器不要在这些函数的开头和结束部分插入这些检查指令

### 栈拷贝(stack copying)

栈拷贝初始阶段与分段栈类似。goroutine在栈上运行着，当用光栈空间，它遇到与旧方案中相同的栈溢出检查。但是与旧方案采用的保留一个返 回前一段栈的link不同，新方案创建一个两倍于原stack大小的新stack，并将旧栈拷贝到其中。

这意味着当栈实际使用的空间缩小为原先的 大小时，go运行时不用做任何事情。栈缩小是一个无任何代价的操作。此外，当栈再次增长时，运行时也无需做任何事情，我们只需要重用之前分配的空 闲空间即可


## 分段栈和连续栈

### 分段栈(Segmented Stacks)
在 Go 1.3 版本之前 ，使用的栈结构是分段栈，每个go函数在函数入口处都会有一小段代码(called prologue)，这段代码会检查是否用光了已分配的栈空间,随着goroutine 调用的函数层级的深入或者局部变量需要的越来越多时，
运行时会调用 runtime.morestack 和 runtime.newstack创建一个新的栈空间，这些栈空间是不连续的，但是当前 goroutine 的多个栈空间会以双向链表的形式串联起来，运行时会通过指针找到连续的栈片段。

分段栈虽然能够按需为当前 goroutine 分配内存并且及时减少内存的占用，但是它也存在一个比较大的问题：

如果当前 goroutine 的栈几乎充满，那么任意的函数调用都会触发栈的扩容，当函数返回后又会触发栈的收缩，如果在一个循环中调用函数，栈的分配和释放就会造成巨大的额外开销，这被称为热分裂问题(Hot split)。

为了解决这个问题，Go 在 1.2 版本的时候不得不将栈的初始化内存从 4KB 增大到了 8KB。后来把采用连续栈结构后，又把初始栈大小减小到了 2KB。


### 连续栈(continuous stacks）

## 栈的分配
源码路径：src/runtime/stack.go

### 小栈内存分配
```go
func stackalloc(n uint32) stack { 
    // 这里的 G 是 G0
    thisg := getg()
    ...
    var v unsafe.Pointer
    // 在 Linux 上，_FixedStack = 2048、_NumStackOrders = 4、_StackCacheSize = 32768
    // 如果申请的栈空间小于 32KB
    if n < _FixedStack<<_NumStackOrders && n < _StackCacheSize {
        order := uint8(0)
        n2 := n
        // 大于 2048 ,那么 for 循环 将 n2 除 2,直到 n 小于等于 2048
        for n2 > _FixedStack {
            // order 表示除了多少次
            order++
            n2 >>= 1
        }
        var x gclinkptr
        //preemptoff != "", 在 GC 的时候会进行设置,表示如果在 GC 那么从 stackpool 分配
        // thisg.m.p = 0 会在系统调用和 改变 P 的个数的时候调用,如果发生,那么也从 stackpool 分配
        if stackNoCache != 0 || thisg.m.p == 0 || thisg.m.preemptoff != "" { 
            lock(&stackpool[order].item.mu)
            // 从 stackpool 分配
            x = stackpoolalloc(order)
            unlock(&stackpool[order].item.mu)
        } else {
            // 从 P 的 mcache 分配内存
            c := thisg.m.p.ptr().mcache
            x = c.stackcache[order].list
            if x.ptr() == nil {
                // 从堆上申请一片内存空间填充到stackcache中
                stackcacherefill(c, order)
                x = c.stackcache[order].list
            }
            // 移除链表的头节点
            c.stackcache[order].list = x.ptr().next
            c.stackcache[order].size -= uintptr(n)
        }
        // 获取到分配的span内存块
        v = unsafe.Pointer(x)
    } else {
        ...
    }
    ...
    return stack{uintptr(v), uintptr(v) + uintptr(n)}
}
```

stackalloc 会根据传入的参数 n 的大小进行分配，在 Linux 上如果 n 小于 32768 bytes，也就是 32KB ，那么会进入到小栈的分配逻辑中。

小栈指大小为 2K/4K/8K/16K 的栈，在分配的时候，会根据大小计算不同的 order 值，如果栈大小是 2K，那么 order 就是 0，4K 对应 order 就是 1，以此类推。这样一方面可以减少不同 Goroutine 获取不同栈大小的锁冲突，另一方面可以预先缓存对应大小的 span ，以便快速获取。

thisg.m.p == 0可能发生在系统调用 exitsyscall 或改变 P 的个数 procresize 时，thisg.m.preemptoff != ""会发生在 GC 的时候。也就是说在发生在系统调用 exitsyscall 或改变 P 的个数在变动，亦或是在 GC 的时候，会从 stackpool 分配栈空间，否则从 mcache 中获取。

如果 mcache 对应的 stackcache 获取不到，那么调用 stackcacherefill 从堆上申请一片内存空间填充到stackcache中。


### 大栈内存分配

```go
func stackalloc(n uint32) stack { 
    thisg := getg() 
    var v unsafe.Pointer

    if n < _FixedStack<<_NumStackOrders && n < _StackCacheSize {
        // 32k...
    } else {
        // 申请的内存空间过大，从 runtime.stackLarge 中检查是否有剩余的空间
        var s *mspan
        // 计算需要分配多少个 span 页， 8KB 为一页
        npage := uintptr(n) >> _PageShift
        // 计算 npage 能够被2整除几次，用来作为不同大小内存的块的索引
        log2npage := stacklog2(npage)

        lock(&stackLarge.lock)
        // 如果 stackLarge 对应的链表不为空
        if !stackLarge.free[log2npage].isEmpty() {
            //获取链表的头节点，并将其从链表中移除
            s = stackLarge.free[log2npage].first
            stackLarge.free[log2npage].remove(s)
        }
        unlock(&stackLarge.lock)

        lockWithRankMayAcquire(&mheap_.lock, lockRankMheap)
        //这里是stackLarge为空的情况
        if s == nil {
            // 从堆上申请新的内存 span
            s = mheap_.allocManual(npage, &memstats.stacks_inuse)
            if s == nil {
                throw("out of memory")
            }
            // OpenBSD 6.4+ 系统需要做额外处理
            osStackAlloc(s)
            s.elemsize = uintptr(n)
        }
        v = unsafe.Pointer(s.base())
    }
    ...
    return stack{uintptr(v), uintptr(v) + uintptr(n)}
}
```



## 参考资料
1. [解密Go协程的栈内存管理](https://mp.weixin.qq.com/s/ErnQDHeL5K8MPDYUPwjSYA)


