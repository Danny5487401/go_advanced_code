# 高级汇编语言--以stack操作为例

Go汇编语言其实是一种高级的汇编语言。在这里高级一词并没有任何褒义或贬义的色彩，而是要强调Go汇编代码和最终真实执行的代码并不完全等价。
Go汇编语言中一个指令在最终的目标代码中可能会被编译为其它等价的机器指令。Go汇编实现的函数或调用函数的指令在最终代码中也会被插入额外的指令。
要彻底理解Go汇编语言就需要彻底了解汇编器到底插入了哪些指令。

##  函数声明

![](.func_images/function_in_plan9.png)
```css
pkgname 包名可以不写，一般都是不写的，可以参考 go 的源码， 另外 add 前的 · 不是 .

代码存储在TEXT段中
                           argsize参数及返回值大小,例如入参是 3 个 int64 类型，返回值是 1 个 int64 类型，那么返回值就是 sizeof(int64) * 4,不过真实世界永远没有我们假设的这么美好，函数参数往往混合了多种类型，还需要考虑内存对齐问题。
                                 |
 TEXT pkgname·add(SB),NOSPLIT,$0-16    -->$framesize-argsize   
         |     |               |
        包名  函数名    framesize栈帧大小(局部变量+如果有对其它函数调用时的话，调用时需要将 callee 的参数、返回值考虑在内。虽然 return address(rip)的值也是存储在 caller 的 stack frame 上的，但是这个过程是由 CALL 指令和 RET 指令完成 PC 寄存器的保存和恢复的，在手写汇编时，同样也是不需要考虑这个 PC 寄存器在栈上所需占用的 8 个字节的)
```
- 为什么要叫 TEXT ？如果对程序数据在文件中和内存中的分段稍有了解的同学应该知道，我们的代码在二进制文件中，是存储在 .text 段中的，这里也就是一种约定俗成的起名方式。实际上在 plan9 中 TEXT 是一个指令，用来定义一个函数。除了 TEXT 之外还有前面变量声明说到的 DATA/GLOBL

- 定义中的 pkgname 部分是可以省略的，非想写也可以写上。不过写上 pkgname 的话，在重命名 package 之后还需要改代码，所以推荐最好还是不要写

- 中点 · 比较特殊，是一个 unicode 的中点，该点在 mac 下的输入方法是 option+shift+9。在程序被链接之后，所有的中点 · 都会被替换为句号 . ，比如你的方法是 runtime·main，在编译之后的程序里的符号则是 runtime.main。

- 以上使用的 RODATA，NOSPLIT flag，还有其他的值，可以参考：https://golang.org/doc/asm#directives

```shell
#include textflag.h

NOPROF = 1
#(For TEXT items.) Don’t profile the marked function. This flag is deprecated.

DUPOK = 2 
# DUPOK表示该变量对应的标识符可能有多个，在链接时只选择其中一个即可（一般用于合并相同的常量字符串，减少重复数据占用的空间）。

NOSPLIT = 4
# 不会生成或包含栈分裂代码，这一般用于没有任何其它函数调用的叶子函数，这样可以适当提高性能。
#(代码段.) Don’t insert the preamble to check if the stack must be split. 
# The frame for the routine, plus anything it calls, must fit in the spare space at the top of the stack segment. 
# Used to protect routines such as the stack splitting code itself.

RODATA = 8
#RODATA标志表示将变量定义在只读内存段，因此后续任何对此变量的修改操作将导致异常（recover也无法捕获）。

NOPTR = 16
#NOPTR则表示此变量的内部不含指针数据，让垃圾回收器忽略对该变量的扫描。如果变量已经在Go代码中声明过的话，Go编译器会自动分析出该变量是否包含指针，这种时候可以不用手写NOPTR标志

WRAPPER = 32
#(代码段.) WRAPPER标志则表示这个是一个包装函数，在panic或runtime.caller等某些处理函数帧的地方不会增加函数帧计数。

NEEDCTXT = 64
#(代码段.) 表示需要一个上下文参数，一般用于闭包函数.
```
当使用这些 flag 的字面量时，需要在汇编文件中 #include "textflag.h"

- framesize:
   - 原则上来说，调用函数时只要不把局部变量覆盖掉就可以了。稍微多分配几个字节的 framesize 也不会死。
   - 在确保逻辑没有问题的前提下，你愿意覆盖局部变量也没有问题。只要保证进入和退出汇编函数时的 caller 和 callee 能正确拿到返回值就可以






## GO版本变化： 函数调用时，传递参数做了修改
- go1.17之前，函数参数是通过栈空间来传递的
- go1.17时做出了改变，在一些平台上（AMD64）可以像C,C++那样使用寄存器传递参数和函数返回值

栈空间：
- 优点：实现简单，不用区分不同的平台，通用性强
- 缺点：效率低
寄存器：
- 优点：速度快
- 缺点：通用性差，不同的平台需要单独处理 当然，这里说的通用性差是对于编译器来说的


## 栈
![](.func_images/stack_memery_info.png)


调用栈call stack，简称栈，是一种栈数据结构，用于存储有关计算机程序的活动 subroutines 信息。在计算机编程中，subroutines 是执行特定任务的一系列程序指令，打包为一个单元。

栈帧stack frame又常被称为帧frame是在调用栈中储存的函数之间的调用关系，每一帧对应了函数调用以及它的参数数据。

有了函数调用自然就要有调用者 caller 和被调用者 callee ，如在 函数 A 里 调用 函数 B，A 是 caller，B 是 callee。


## 栈分裂stack-split
由于 Go 程序中的 goroutine 数目是不可确定的，并且实际场景可能会有百万级别的 goroutine，runtime 必须使用保守的思路来给 goroutine 分配空间以避免吃掉所有的可用内存。
也由于此，每个新的 goroutine 会被 runtime 分配初始为 2KB 大小的栈空间(Go 的栈在底层实际上是分配在堆空间上的)。
随着一个 goroutine 进行自己的工作，可能会超出最初分配的栈空间限制(就是栈溢出的意思)。 

为了防止这种情况发生，runtime 确保 goroutine 在超出栈范围时，会创建一个相当于原来两倍大小的新栈，并将原来栈的上下文拷贝到新栈上。
这个过程被称为 栈分裂(stack-split)，这样使得 goroutine 栈能够动态调整大小.

为了使栈分裂正常工作，编译器会在每一个函数的开头和结束位置插入指令来防止 goroutine 爆栈。
像我们本章早些看到的一样，为了避免不必要的开销，一定不会爆栈的函数会被标记上 NOSPLIT 来提示编译器不要在这些函数的开头和结束部分插入这些检查指令

## 栈结构
典型的函数的栈结构图

![](.func_images/func_call_frame.png)
![](.func_images/layout_of_function_args_n_return_value.png)
![](.func_images/swap.png)

```css
低地址
                       -----------------                                           
                       current func arg0                                           
                       ----------------- <----------- FP(pseudo FP)                
                        caller ret addr                                            
                       +---------------+                                           
                       | caller BP(*)  |                                           
                       ----------------- <----------- SP(pseudo SP，实际上是当前栈帧的 BP 位置)
                       |   Local Var0  |                                           
                       -----------------                                           
                       |   Local Var1  |                                           
                       -----------------                                           
                       |   Local Var2  |                                           
                       -----------------                -                          
                       |   ........    |                                           
                       -----------------                                           
                       |   Local VarN  |                                           
                       -----------------                                           
                       |               |                                           
                       |               |                                           
                       |  temporarily  |                                           
                       |  unused space |                                           
                       |               |                                           
                       |               |                                           
                       -----------------                                           
                       |  call retn    |                                           
                       -----------------                                           
                       |  call ret(n-1)|                                           
                       -----------------                                           
                       |  ..........   |                                           
                       -----------------                                           
                       |  call ret1    |                                           
                       -----------------                                           
                       |  call argn    |                                           
                       -----------------                                           
                       |   .....       |                                           
                       -----------------                                           
                       |  call arg3    |                                           
                       -----------------                                           
                       |  call arg2    |                                           
                       |---------------|                                           
                       |  call arg1    |                                           
                       -----------------   <------------  hardware SP 位置           
                         return addr                                               
                       +---------------+                                           
                                                                                   
高地址
```
从原理上来讲，如果当前函数调用了其它函数，那么 return addr 也是在 caller 的栈上的，不过往栈上插 return addr 的过程是由 CALL 指令完成的，在 RET 时，SP 又会恢复到图上位置。我们在计算 SP 和参数相对位置时，可以认为硬件 SP 指向的就是图上的位置。

图上的 caller BP，指的是 caller 的 BP 寄存器值，有些人把 caller BP 叫作 caller 的 frame pointer，实际上这个习惯是从 x86 架构沿袭来的。Go 的 asm 文档中把伪寄存器 FP 也称为 frame pointer，但是这两个 frame pointer 根本不是一回事。

此外需要注意的是，caller BP 是在编译期由编译器插入的，用户手写代码时，计算 frame size 时是不包括这个 caller BP 部分的。是否插入 caller BP 的主要判断依据是:

1. 函数的栈帧大小大于 0
2. 下述函数返回 true
```go
func Framepointer_enabled(goos, goarch string) bool {
    return framepointer_enabled != 0 && goarch == "amd64" && goos != "nacl"
}
```
如果编译器在最终的汇编结果中没有插入 caller BP(源代码中所称的 frame pointer)的情况下，伪 SP 和伪 FP 之间只有 8 个字节的 caller 的 return address，而插入了 BP 的话，就会多出额外的 8 字节。也就说伪 SP 和伪 FP 的相对位置是不固定的，有可能是间隔 8 个字节，也有可能间隔 16 个字节。并且判断依据会根据平台和 Go 的版本有所不同。

- FP 伪寄存器指向函数的传入参数的开始位置，因为栈是朝低地址方向增长，为了通过寄存器引用参数时方便，所以参数的摆放方向和栈的增长方向是相反的，即：
```css
                              FP
high ----------------------> low
argN, ... arg3, arg2, arg1, arg0
```

假设所有参数均为 8 字节，这样我们就可以用 symname+0(FP) 访问第一个 参数，symname+8(FP) 访问第二个参数，以此类推。

- 用伪 SP 来引用局部变量，原理上来讲差不多，不过因为伪 SP 指向的是局部变量的底部，所以 symname-8(SP) 表示的是第一个局部变量，symname-16(SP)表示第二个，以此类推。当然，这里假设局部变量都占用 8 个字节。

图的最上部的 caller return address 和 current func arg0 都是由 caller 来分配空间的。不算在当前的栈帧内。



## Goroutine 栈操作
![](.func_images/stack_high_n_low.png)
在 Goroutine 中有一个 stack 数据结构，里面有两个属性 lo 与 hi，描述了实际的栈内存地址：

- stack.lo：栈空间的低地址；
- stack.hi：栈空间的高地址

在 Goroutine 中会通过 stackguard0 来判断是否要进行栈增长：

- stackguard0：stack.lo + StackGuard, 用于stack overlow的检测；
- StackGuard：保护区大小，常量Linux上为 928 字节；
- StackSmall：常量大小为 128 字节，用于小函数调用的优化；
- StackBig：常量大小为 4096 字节；

需要注意的是，由于栈是由高地址向低地址增长的，所以对比的时候，都是小于才执行扩容，这里需要大家品品。


## G 的创建
因为栈都是在 Goroutine 上的，所以先从 G 的创建开始看如何创建以及初始化栈空间的。

G 的创建会调用 runtime·newproc进行创建：
```go
func newproc(siz int32, fn *funcval) {
    argp := add(unsafe.Pointer(&fn), sys.PtrSize)
    gp := getg()
    // 获取 caller 的 PC 寄存器
    pc := getcallerpc()
    // 切换到 G0 进行创建
    systemstack(func() {
        newg := newproc1(fn, argp, siz, gp, pc)
        ...
    })
}
```
newproc 方法会切换到 G0 上调用 newproc1 函数进行 G 的创建。
```go
const _StackMin = 2048
func newproc1(fn *funcval, argp unsafe.Pointer, narg int32, callergp *g, callerpc uintptr) *g {
    _g_ := getg()
    ...
    _p_ := _g_.m.p.ptr()
    // 从 P 的空闲链表中获取一个新的 G
    newg := gfget(_p_)
    // 获取不到则调用 malg 进行创建
    if newg == nil {
        newg = malg(_StackMin)
        casgstatus(newg, _Gidle, _Gdead)
        allgadd(newg) // publishes with a g->status of Gdead so GC scanner doesn't look at uninitialized stack.
    }
    ...
    return newg
}
```

newproc1 方法很长，里面主要是获取 G ，然后对获取到的 G 做一些初始化的工作。我们这里只看 malg 函数的调用。

在调用 malg 函数的时候会传入一个最小栈大小的值：_StackMin（2048）。
```go
func malg(stacksize int32) *g {
    // 创建 G 结构体
    newg := new(g)
    if stacksize >= 0 {
        // 这里会在 stacksize 的基础上为每个栈预留系统调用所需的内存大小 _StackSystem
        // 在 Linux/Darwin 上（ _StackSystem == 0 ）本行不改变 stacksize 的大小
        stacksize = round2(_StackSystem + stacksize)
        // 切换到 G0 为 newg 初始化栈内存
        systemstack(func() {
            newg.stack = stackalloc(uint32(stacksize))
        })
        // 设置 stackguard0 ，用来判断是否要进行栈扩容
        newg.stackguard0 = newg.stack.lo + _StackGuard
        newg.stackguard1 = ^uintptr(0) 
        *(*uintptr)(unsafe.Pointer(newg.stack.lo)) = 0
    }
    return newg
}
```
在调用 malg 的时候会将传入的内存大小加上一个 _StackSystem 值预留给系统调用使用，round2 函数会将传入的值舍入为 2 的指数。然后会切换到 G0 执行 stackalloc 函数进行栈内存分配。

分配完毕之后会设置 stackguard0 为 stack.lo + _StackGuard，作为判断是否需要进行栈扩容使用


## 栈的初始化
```go
const (
    // Number of orders that get caching. Order 0 is FixedStack
    // and each successive order is twice as large.
    // We want to cache 2KB, 4KB, 8KB, and 16KB stacks. Larger stacks
    // will be allocated directly.
    // Since FixedStack is different on different systems, we
    // must vary NumStackOrders to keep the same maximum cached size.
    //   OS               | FixedStack | NumStackOrders
    //   -----------------+------------+---------------
    //   linux/darwin/bsd | 2KB        | 4
    //   windows/32       | 4KB        | 3
    //   windows/64       | 8KB        | 2
    //   plan9            | 4KB        | 3
    _NumStackOrders = 4 - sys.PtrSize/4*sys.GoosWindows - 1*sys.GoosPlan9
)
// src/runtime/stack.go
// 全局的栈缓存，分配 32KB以下内存
var stackpool [_NumStackOrders]struct {
    item stackpoolItem
    _    [cpu.CacheLinePadSize - unsafe.Sizeof(stackpoolItem{})%cpu.CacheLinePadSize]byte
}

//go:notinheap
type stackpoolItem struct {
    mu   mutex
    span mSpanList
} 

// 全局的栈缓存，分配 32KB 以上内存
var stackLarge struct {
    lock mutex
    free [heapAddrBits - pageShift]mSpanList // free lists by log_2(s.npages)
}

// 初始化stackpool/stackLarge全局变量
func stackinit() {
    if _StackCacheSize&_PageMask != 0 {
        throw("cache size must be a multiple of page size")
    }
    for i := range stackpool {
        stackpool[i].item.span.init()
        lockInit(&stackpool[i].item.mu, lockRankStackpool)
    }
    for i := range stackLarge.free {
        stackLarge.free[i].init()
        lockInit(&stackLarge.lock, lockRankStackLarge)
    }
}
```

需要注意的是，stackinit 是在调用 runtime·schedinit初始化的，是在调用 runtime·newproc之前进行的。

在执行栈初始化的时候会初始化两个全局变量 stackpool 和 stackLarge。stackpool 可以分配小于 32KB 的内存，stackLarge 用来分配大于 32KB 的栈空间。


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

## 案例
为了便于分析，我们先构造一个禁止栈分裂的printnl函数。printnl函数内部都通过调用runtime.printnl函数输出换行：
```
TEXT ·printnl_nosplit(SB), NOSPLIT, $8
	CALL runtime·printnl(SB)
	RET
```

然后通过`go tool asm -S main_amd64.s`指令查看编译后的目标代码：

```shell
"".printnl_nosplit STEXT nosplit size=29 args=0xffffffff80000000 locals=0x10
0x0000 00000 (main_amd64.s:5) TEXT "".printnl_nosplit(SB), NOSPLIT	$16
0x0000 00000 (main_amd64.s:5) SUBQ $16, SP

0x0004 00004 (main_amd64.s:5) MOVQ BP, 8(SP)
0x0009 00009 (main_amd64.s:5) LEAQ 8(SP), BP

0x000e 00014 (main_amd64.s:6) CALL runtime.printnl(SB)

0x0013 00019 (main_amd64.s:7) MOVQ 8(SP), BP
0x0018 00024 (main_amd64.s:7) ADDQ $16, SP
0x001c 00028 (main_amd64.s:7) RET
```
输出代码中我们删除了非指令的部分。为了便于讲述，我们将上述代码重新排版，并根据缩进表示相关的功能：

```shell
TEXT "".printnl(SB), NOSPLIT, $16
	SUBQ $16, SP
		MOVQ BP, 8(SP)
		LEAQ 8(SP), BP
			CALL runtime.printnl(SB)
		MOVQ 8(SP), BP
	ADDQ $16, SP
RET
```

第一层是TEXT指令表示函数开始，到RET指令表示函数返回。第二层是`SUBQ $16, SP`指令为当前函数帧分配16字节的空间，在函数返回前通过`ADDQ $16, SP`指令回收16字节的栈空间。
我们谨慎猜测在第二层是为函数多分配了8个字节的空间。那么为何要多分配8个字节的空间呢？
再继续查看第三层的指令：开始部分有两个指令`MOVQ BP, 8(SP)`和`LEAQ 8(SP), BP`，首先是将BP寄存器保持到多分配的8字节栈空间，然后将`8(SP)`地址重新保持到了BP寄存器中；
结束部分是`MOVQ 8(SP), BP`指令则是从栈中恢复之前备份的前BP寄存器的值。最里面第四次层才是我们写的代码，调用runtime.printnl函数输出换行。

如果去掉NOSPILT标志，再重新查看生成的目标代码，会发现在函数的开头和结尾的地方又增加了新的指令。下面是经过缩进格式化的结果：

```
TEXT "".printnl_nosplit(SB), $16
L_BEGIN:
	MOVQ (TLS), CX
	CMPQ SP, 16(CX)
	JLS  L_MORE_STK

		SUBQ $16, SP
			MOVQ BP, 8(SP)
			LEAQ 8(SP), BP
				CALL runtime.printnl(SB)
			MOVQ 8(SP), BP
		ADDQ $16, SP

L_MORE_STK:
	CALL runtime.morestack_noctxt(SB)
	JMP  L_BEGIN
RET
```
其中开头有三个新指令，`MOVQ (TLS), CX`用于加载g结构体指针，然后第二个指令`CMPQ SP, 16(CX)`SP栈指针和g结构体中stackguard0成员比较，如果比较的结果小于0则跳转到结尾的L_MORE_STK部分。当获取到更多栈空间之后，通过`JMP L_BEGIN`指令跳转到函数的开始位置重新进行栈空间的检测。

g结构体在`$GOROOT/src/runtime/runtime2.go`文件定义，开头的结构成员如下：
![](.introduction_images/goroutine_stack.png)
```go
type g struct {
	// Stack parameters.
	stack       stack   // offset known to runtime/cgo
	stackguard0 uintptr // offset known to liblink
	stackguard1 uintptr // offset known to liblink

	...
}
```

第一个成员是stack类型，表示当前栈的开始和结束地址。stack的定义如下：
```go
// Stack describes a Go execution stack.
// The bounds of the stack are exactly [lo, hi),
// with no implicit data structures on either side.
type stack struct {
	lo uintptr  //栈空间的低地址；
	hi uintptr  // 栈空间的高地址；
}
```
在g结构体中的stackguard0成员是出现爆栈前的警戒线。

1. 扩容：在 Goroutine 中会通过 stackguard0 来判断是否要进行栈增长：
```go
func newstack() {
    thisg := getg() 

    gp := thisg.m.curg

    // 初始化寄存器相关变量
    morebuf := thisg.m.morebuf
    thisg.m.morebuf.pc = 0
    thisg.m.morebuf.lr = 0
    thisg.m.morebuf.sp = 0
    thisg.m.morebuf.g = 0
    ...
    // 校验是否被抢占
    preempt := atomic.Loaduintptr(&gp.stackguard0) == stackPreempt

    // 如果被抢占
    if preempt {
        // 校验是否可以安全的被抢占
        // 如果 M 上有锁
        // 如果正在进行内存分配
        // 如果明确禁止抢占
        // 如果 P 的状态不是 running
        // 那么就不执行抢占了
        if !canPreemptM(thisg.m) {
            // 到这里表示不能被抢占？
            // Let the goroutine keep running for now.
            // gp->preempt is set, so it will be preempted next time.
            gp.stackguard0 = gp.stack.lo + _StackGuard
            // 触发调度器的调度
            gogo(&gp.sched) // never return
        }
    }

    if gp.stack.lo == 0 {
        throw("missing stack in newstack")
    }
    // 寄存器 sp
    sp := gp.sched.sp
    if sys.ArchFamily == sys.AMD64 || sys.ArchFamily == sys.I386 || sys.ArchFamily == sys.WASM {
        // The call to morestack cost a word.
        sp -= sys.PtrSize
    } 
    ...
    if preempt {
        //需要收缩栈
        if gp.preemptShrink { 
            gp.preemptShrink = false
            shrinkstack(gp)
        }
        // 被 runtime.suspendG 函数挂起
        if gp.preemptStop {
            // 被动让出当前处理器的控制权
            preemptPark(gp) // never returns
        }

        //主动让出当前处理器的控制权
        gopreempt_m(gp) // never return
    }

    // 计算新的栈空间是原来的两倍
    oldsize := gp.stack.hi - gp.stack.lo
    newsize := oldsize * 2 
    ... 
    //将 Goroutine 切换至 _Gcopystack 状态
    casgstatus(gp, _Grunning, _Gcopystack)

    //开始栈拷贝
    copystack(gp, newsize) 
    casgstatus(gp, _Gcopystack, _Grunning)
    gogo(&gp.sched)
}
```
在开始执行栈拷贝之前会先计算新栈的大小是原来的两倍，然后将 Goroutine 状态切换至 _Gcopystack 状态。


2. 收缩

   我们知道Go运行时会定期进行垃圾回收操作，这其中包含栈的回收工作。如果栈使用到比例小于一定到阈值，则分配一个较小到栈空间，然后将栈上面到数据移动到新的栈中，
   栈移动的过程和栈扩容的过程类似.
```go
func scanstack(gp *g, gcw *gcWork) {
    ... 
    // 进行栈收缩
    shrinkstack(gp)
    ...
}
```

```go
func shrinkstack(gp *g) {
    ...
    oldsize := gp.stack.hi - gp.stack.lo
    newsize := oldsize / 2 
    // 当收缩后的大小小于最小的栈的大小时，不再进行收缩
    if newsize < _FixedStack {
        return
    }
    avail := gp.stack.hi - gp.stack.lo
    // 计算当前正在使用的栈数量，如果 gp 使用的当前栈少于四分之一，则对栈进行收缩
    // 当前使用的栈包括到 SP 的所有内容以及栈保护空间，以确保有 nosplit 功能的空间
    if used := gp.stack.hi - gp.sched.sp + _StackLimit; used >= avail/4 {
        return
    }
    // 将旧栈拷贝到新收缩后的栈上
    copystack(gp, newsize)
}
```
新栈的大小会缩小至原来的一半，如果小于 _FixedStack （2KB）那么不再进行收缩。除此之外还会计算一下当前栈的使用情况是否不足 1/4 ，如果使用超过 1/4 那么也不会进行收缩。

最后判断确定要进行收缩则调用 copystack 函数进行栈拷贝的逻辑

