# 扩大原子操作的适用范围：atomic.Value

    在 Go 语言标准库中，sync/atomic包将底层硬件提供的原子操作封装成了 Go 的函数。但这些操作只支持几种基本数据类型，因此为了扩大原子操作的适用范围，
    Go 语言在 1.4 版本的时候向sync/atomic包中添加了一个新的类型Value。此类型的值相当于一个容器，可以被用来“原子地"存储（Store）和加载（Load）任意类型的值
    

背景：

    我在golang-dev邮件列表中翻到了14年的这段讨论，有用户报告了encoding/gob包在多核机器上（80-core）上的性能问题，
    认为encoding/gob之所以不能完全利用到多核的特性是因为它里面使用了大量的互斥锁（mutex），
    如果把这些互斥锁换成用atomic.LoadPointer/StorePointer来做并发控制，那性能将能提升20倍。

做法：

    有人提议在已有的atomic包的基础上封装出一个atomic.Value类型，这样用户就可以在不依赖 Go 内部类型unsafe.Pointer的情况下使用到atomic提供的原子操作。
    所以我们现在看到的atomic包中除了atomic.Value外，其余都是早期由汇编写成的，并且atomic.Value类型的底层实现也是建立在已有的atomic包的基础上

问题：

    为什么在上面的场景中，atomic会比mutex性能好很多呢？
原因:

    Mutexes do no scale. Atomic loads do.------ Dmitry Vyukov
    Mutex由操作系统实现，而atomic包中的原子操作则由底层硬件直接提供支持。在 CPU 实现的指令集里，有一些指令被封装进了atomic包，这些指令在执行的过程中是不允许中断（interrupt）的，
    因此原子操作可以在lock-free的情况下保证并发安全，并且它的性能也能做到随 CPU 个数的增多而线性扩展。

原子性
![](./img/write_value_process.png)

    一个或者多个操作在 CPU 执行的过程中不被中断的特性，称为原子性（atomicity） 。这些操作对外表现成一个不可分割的整体，他们要么都执行，要么都不执行，外界不会看到他们只执行到一半的状态。
    而在现实世界中，CPU 不可能不中断的执行一系列操作，但如果我们在执行多个操作时，能让他们的中间状态对外不可见，那我们就可以宣称他们拥有了"不可分割”的原子性。
    有些朋友可能不知道，在 Go（甚至是大部分语言）中，一条普通的赋值语句其实不是一个原子操作。

    例如，在32位机器上写int64类型的变量就会有中间状态，因为它会被拆成两次写操作（MOV）——写低 32 位和写高 32 位.
    如果一个线程刚写完低32位，还没来得及写高32位时，另一个线程读取了这个变量，那它得到的就是一个毫无逻辑的中间变量，这很有可能使我们的程序出现诡异的 Bug。
    
    这还只是一个基础类型，如果我们对一个结构体进行赋值，那它出现并发问题的概率就更高了。很可能写线程刚写完一小半的字段，读线程就来读取这个变量，那么就只能读到仅修改了一部分的值。
    这显然破坏了变量的完整性，读出来的值也是完全错误的。

    面对这种多线程下变量的读写问题，我们的主角——atomic.Value登场了，它使得我们可以不依赖于不保证兼容性的unsafe.Pointer类型，同时又能将任意数据类型的读写操作封装成原子性操作（让中间状态对外不可见）
    atomic.Value类型使用：

操作的数据类型几种操作：

    六种类型：int32, int64, uint32, uint64, uintptr, unsafe.Pinter
    五种操作：Add增减， CompareAndSwap比较并交换， Swap交换， Load载入， Store存储
    函数名以"操作+类型"组合而来。例如AddInt32/AddUint64/LoadInt32/LoadUint64.

```go
//1. LoadXXX(addr): 原子性的获取*addr的值，等价于
return *addr

//2.StoreXXX(addr, val): 原子性的将val的值保存到*addr，等价于：
addr = val

//3.AddXXX(addr, delta): 原子性的将delta的值添加到*addr并返回新值（unsafe.Pointer不支持），等价于
*addr += delta
return *addr

//4.SwapXXX(addr, new) old: 原子性的将new的值保存到*addr并返回旧值，等价于：
old = *addr
*addr = new
return old
//5.CompareAndSwapXXX(addr, old, new) bool: 原子性的比较*addr和old，如果相同则将new赋值给*addr并返回true
if *addr == old {
*addr = new
return true
}
return false
```

## atomic.Value源码分析

```go

type Value struct {
    v interface{}
}
func(v *Value) Load() (x interface{}) //读操作，从线程安全的v中读取上一步存放的内容
func(v *Value) Store(x interface{}) //写操作，将原始的变量x存放在atomic.Value类型的v


// ifaceWords is interface{} internal representation.
type ifaceWords struct {
    typ  unsafe.Pointer
    data unsafe.Pointer
}
```
1.unsafe.Pointer

    Go语言并不支持直接操作内存，但是它的标准库提供一种不保证向后兼容的指针类型unsafe.Pointer， 让程序可以灵活的操作内存，它的特别之处在于：可以绕过Go语言类型系统的检查
    
    也就是说：如果两种类型具有相同的内存结构，我们可以将unsafe.Pointer当作桥梁，让这两种类型的指针相互转换，从而实现同一份内存拥有两种解读方式
    int类型和int32类型内部的存储结构是一致的，但是对于指针类型的转换需要这么做：
```go

var a int32
// 获得a的*int类型指针
(*int)(unsafe.Pointer(&a))
```

![](./img/load_n_store_process.png)

2.实现原子性的读取任意结构操作
```go

func (v *Value) Load() (x interface{}) {
    // 将*Value指针类型转换为*ifaceWords指针类型
	vp := (*ifaceWords)(unsafe.Pointer(v))
	// 原子性的获取到v的类型typ的指针
	typ := LoadPointer(&vp.typ)
	// 如果没有写入或者正在写入，先返回，^uintptr(0)代表过渡状态，见下文
	if typ == nil || uintptr(typ) == ^uintptr(0) {
		return nil
	}
	// 原子性的获取到v的真正的值data的指针，然后返回
	data := LoadPointer(&vp.data)
	xp := (*ifaceWords)(unsafe.Pointer(&x))
	xp.typ = typ
	xp.data = data
	return
}
```
3.实现原子性的存储任意结构操作
```go
func runtime_procPin()  //可以将一个goroutine死死占用当前使用的P ,不允许其他的goroutine抢占
func runtime_procUnpin() //释放
```
```go
func (v *Value) Store(x interface{}) {
	if x == nil {
		panic("sync/atomic: store of nil value into Value")
	}
	// 将现有的值和要写入的值转换为ifaceWords类型，这样下一步就能获取到它们的原始类型和真正的值
	vp := (*ifaceWords)(unsafe.Pointer(v))
	xp := (*ifaceWords)(unsafe.Pointer(&x))
	for {
		// 获取现有的值的type
		typ := LoadPointer(&vp.typ)
		// 如果typ为nil说明这是第一次Store
		if typ == nil {
			// 如果你是第一次，就死死占住当前的processor，不允许其他goroutine再抢
			runtime_procPin()
			// 使用CAS操作，先尝试将typ设置为^uintptr(0)这个中间状态
			// 如果失败，则证明已经有别的线程抢先完成了赋值操作
			// 那它就解除抢占锁，然后重新回到 for 循环第一步
			if !CompareAndSwapPointer(&vp.typ, nil, unsafe.Pointer(^uintptr(0))) {
				runtime_procUnpin()
				continue
			}
			// 如果设置成功，说明当前goroutine中了jackpot
			// 那么就原子性的更新对应的指针，最后解除抢占锁
			StorePointer(&vp.data, xp.data)
			StorePointer(&vp.typ, xp.typ)
			runtime_procUnpin()
			return
		}
		// 如果typ为^uintptr(0)说明第一次写入还没有完成，继续循环等待
		if uintptr(typ) == ^uintptr(0) {
			continue
		}
		// 如果要写入的类型和现有的类型不一致，则panic
		if typ != xp.typ {
			panic("sync/atomic: store of inconsistently typed value into Value")
		}
		// 更新data
		StorePointer(&vp.data, xp.data)
		return
	}
}
```
