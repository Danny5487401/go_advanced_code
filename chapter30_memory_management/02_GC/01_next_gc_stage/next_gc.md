<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [gc 模式](#gc-%E6%A8%A1%E5%BC%8F)
- [gcController 和 work 这两个数据结构](#gccontroller-%E5%92%8C-work-%E8%BF%99%E4%B8%A4%E4%B8%AA%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84)
- [启动 gc](#%E5%90%AF%E5%8A%A8-gc)
- [GC 阶段](#gc-%E9%98%B6%E6%AE%B5)
- [问题](#%E9%97%AE%E9%A2%98)
  - [1. 生产速度（申请内存） 大于消费速度（GC）？](#1-%E7%94%9F%E4%BA%A7%E9%80%9F%E5%BA%A6%E7%94%B3%E8%AF%B7%E5%86%85%E5%AD%98-%E5%A4%A7%E4%BA%8E%E6%B6%88%E8%B4%B9%E9%80%9F%E5%BA%A6gc)
- [观察GC方式](#%E8%A7%82%E5%AF%9Fgc%E6%96%B9%E5%BC%8F)
- [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


## gc 模式

```go
// gcMode indicates how concurrent a GC cycle should be.
type gcMode int

const (
	gcBackgroundMode gcMode = iota // 默认模式，标记与清扫过程都是并发执行的
	gcForceMode                    // stop-the-world GC now, 只在清扫阶段支持并发
	gcForceBlockMode               // stop-the-world GC now and STW sweep (forced by user) GC全程需要STW
)
```

##  gcController 和 work 这两个数据结构
gcController GC 控制器

```go
var gcController gcControllerState

type gcControllerState struct {
    // ...

	// globalsScan is the total amount of global variable space
	// that is scannable.
	globalsScan atomic.Uint64
	
	heapMarked uint64
	
	heapScanWork    atomic.Int64
	stackScanWork   atomic.Int64
	globalsScanWork atomic.Int64

	// 是并发的后台扫描总信用累计值，此值通过 mutator assists 在后台扫描并窃取。它是原子更新操作
	bgScanCredit atomic.Int64

	// 一个周期内辅助助手(mutator assists)所花费的时间，单位 ns
	assistTime atomic.Int64

	// dedicatedMarkTime is the nanoseconds spent in dedicated mark workers
	// during this cycle. This is updated at the end of the concurrent mark
	// phase.
	dedicatedMarkTime atomic.Int64

	// fractionalMarkTime is the nanoseconds spent in the fractional mark
	// worker during this cycle. This is updated throughout the cycle and
	// will be up-to-date if the fractional mark worker is not currently
	// running.
	fractionalMarkTime atomic.Int64

	// 一个周期内空闲标记的时间(ms)。它在整个周期内都会自动更新
	idleMarkTime atomic.Int64

	// 是辅助和后台标记线程启动的绝对开始时间(ms)
	markStartTime int64

	// 需要启动专用标记线程的数量。它会在每个周期开始的时候计算而
	dedicatedMarkWorkersNeeded atomic.Int64

    // ..
}
```

```go
var work workType

type workType struct {
    // ..

	// 周期内已标记的字节大小，包含多种对象，
	// This includes bytes blackened in scanned objects, noscan objects
	// that go straight to black, and permagrey objects scanned by
	// markroot during the concurrent scan phase.
	// 由于标记过程中存在竞争，所以这个数字不是太准确，但也很接近
	bytesMarked uint64

	markrootNext uint32 // next markroot job
	markrootJobs uint32 // number of markroot jobs
    // ..
	

	// GC的工作模式
	mode gcMode


	// 在gc开始时，值为 memstats.heap_live
	initialHeapLive uint64

	// assistQueue is a queue of assists that are blocked because
	// there was neither enough credit to steal or enough work to
	// do.
	assistQueue struct {
		lock mutex
		q    gQueue
	}

	// gc 从 mark termination 转为 sweep时，要唤醒的goroutines
	sweepWaiters struct {
		lock mutex
		list gList
	}

	// 已完成的gc次数
	cycles atomic.Uint32

    // ...
}
```

## 启动 gc
```go
// go1.21.5/src/runtime/proc.go
func init() {
	go forcegchelper()
}

func forcegchelper() {
	forcegc.g = getg()
	lockInit(&forcegc.lock, lockRankForcegc)
	for {
		lock(&forcegc.lock)
		if forcegc.idle.Load() {
			throw("forcegc: phase error")
		}
        // ...
		// Time-triggered, fully concurrent.
		gcStart(gcTrigger{kind: gcTriggerTime, now: nanotime()})
	}
}
```


## GC 阶段


![](../../.asset/img/.gc_images/gc_cycle2.png)


```go
// /Users/python/go/go1.21.5/src/runtime/mgc.go
const (
	_GCoff             = iota // GC未运行; sweeping in background, write barrier disabled
	_GCmark                   // 标记中，启用写屏障 GC marking roots and workbufs: allocate black, write barrier ENABLED
	_GCmarktermination        // 标记终止，启用写屏障 GC mark termination: allocate black, P's help GC, write barrier ENABLED
)

var gcphase uint32
```

整个GC共四个阶段，每次开始时从上到下执行。第一步是清理上次未清理完的span，而不是直接标记阶段
```go
func GC() {
	// We consider a cycle to be: sweep termination, mark, mark
	// termination, and sweep. 

	// 从work全局变量里读取当前已GC的次数
	n := work.cycles.Load()
	// 确保第n次GC "清扫终止"、"标记阶段" 和 ”标记终止“,如果已标记完成会立即返回
	gcWaitOnMark(n)

	// We're now in sweep N or later. Trigger GC cycle N+1, which
	// will first finish sweep N if necessary and then enter sweep
	// termination N+1.
	gcStart(gcTrigger{kind: gcTriggerCycle, n: n + 1})

	// Wait for mark termination N+1 to complete.
	gcWaitOnMark(n + 1)

	// 第四阶段：清扫
	for work.cycles.Load() == n+1 && sweepone() != ^uintptr(0) {
		sweep.nbgsweep++
		Gosched()
	}

	// Callers may assume that the heap profile reflects the
	// just-completed cycle when this returns (historically this
	// happened because this was a STW GC), but right now the
	// profile still reflects mark termination N, not N+1.
	//
	// As soon as all of the sweep frees from cycle N+1 are done,
	// we can go ahead and publish the heap profile.
	//
	// First, wait for sweeping to finish. (We know there are no
	// more spans on the sweep queue, but we may be concurrently
	// sweeping spans, so we have to wait.)
	for work.cycles.Load() == n+1 && !isSweepDone() {
		Gosched()
	}

	// Now we're really done with sweeping, so we can publish the
	// stable heap profile. Only do this if we haven't already hit
	// another mark termination.
	mp := acquirem()
	cycle := work.cycles.Load()
	if cycle == n+1 || (gcphase == _GCmark && cycle == n+2) {
		mProf_PostSweep()
	}
	releasem(mp)
}
```

1. sweep termination（清理终止）

    - 会触发 STW ，所有的 P（处理器） 都会进入 safe-point（安全点）；

2. the mark phase（并发标记阶段）

    - **setGCPhase(_GCmark)**: 将 _GCoff GC 状态 改成 _GCmark，开启 Write Barrier （写入屏障）、mutator assists（协助线程），将根对象入队；
    - 恢复程序执行，mark workers（标记进程）和 mutator assists（协助线程）会开始并发标记内存中的对象。对于任何指针写入和新的指针值，都会被写屏障覆盖，而所有新创建的对象都会被直接标记成黑色；
    - GC 执行根节点的标记，这包括扫描所有的栈、全局对象以及不在堆中的运行时数据结构。扫描goroutine 栈绘导致 goroutine 停止，并对栈上找到的所有指针加置灰，然后继续执行 goroutine。
    - GC 在遍历灰色对象队列的时候，会将灰色对象变成黑色，并将该对象指向的对象置灰；
    - GC 会使用分布式终止算法（distributed termination algorithm）来检测何时不再有根标记作业或灰色对象，如果没有了 GC 会转为mark termination（标记终止）；

3. mark termination（标记终止）

    - **stopTheWorldWithSema(stwGCMarkTerm)** ，然后将 GC 阶段转为 _GCmarktermination,关闭 GC 工作线程以及 mutator assists（协助线程）；
    - 执行清理，如 flush mcache；

4. the sweep phase（清扫阶段）

    - 将 GC 状态转变至 _GCoff，初始化清理状态并关闭 Write Barrier（写入屏障）；
    - 恢复程序执行，从此开始新创建的对象都是白色的；
    - 后台并发清理所有的内存管理单元



```go
func gcStart(trigger gcTrigger) {
	// 判断当前G是否可抢占, 不可抢占时不触发GC
	mp := acquirem()
	if gp := getg(); gp == mp.g0 || mp.locks > 1 || mp.preemptoff != "" {
		releasem(mp)
		return
	}
	releasem(mp)
	mp = nil

	// 并发清除上一轮未清除的span【GC第一阶段：GC 清理终止】
	for trigger.test() && sweepone() != ^uintptr(0) {
		sweep.nbgsweep++
	}

	// 加锁，重新检查是否满足触发GC的条件，不满足的话，解锁并直接返回
	semacquire(&work.startSema)
	if !trigger.test() {
		semrelease(&work.startSema)
		return
	}

	// In gcstoptheworld debug mode, upgrade the mode accordingly.
	// We do this after re-checking the transition condition so
	// that multiple goroutines that detect the heap trigger don't
	// start multiple STW GCs.
	mode := gcBackgroundMode
	if debug.gcstoptheworld == 1 {
		mode = gcForceMode
	} else if debug.gcstoptheworld == 2 {
		mode = gcForceBlockMode
	}

	// 加锁，两个全局锁
	semacquire(&gcsema)
	semacquire(&worldsema)

	// For stats, check if this GC was forced by the user.
	// Update it under gcsema to avoid gctrace getting wrong values.
	work.userForced = trigger.kind == gcTriggerCycle

	if traceEnabled() {
		traceGCStart()
	}

	// Check that all Ps have finished deferred mcache flushes.
	for _, p := range allp {
		if fg := p.mcache.flushGen.Load(); fg != mheap_.sweepgen {
			println("runtime: p", p.id, "flushGen", fg, "!= sweepgen", mheap_.sweepgen)
			throw("p mcache not flushed")
		}
	}
    //启用后台 标记工作 线程
	gcBgMarkStartWorkers()

	// 在系统栈上运行gcResetMarkState()函数，实现在标记（concurrent或STW）之前重置全局状态，并重置所有Gp的栈扫描状态。
	systemstack(gcResetMarkState)

	work.stwprocs, work.maxprocs = gomaxprocs, gomaxprocs
	if work.stwprocs > ncpu {
		// This is used to compute CPU time of the STW phases,
		// so it can't be more than ncpu, even if GOMAXPROCS is.
		work.stwprocs = ncpu
	}
	work.heap0 = gcController.heapLive.Load()
	work.pauseNS = 0 // stw 的累计时间
	work.mode = mode
    // GC开始时间
	now := nanotime()
	work.tSweepTerm = now
	work.pauseStart = now
	systemstack(func() { stopTheWorldWithSema(stwGCSweepTerm) }) //停止所有运行中的G, 并禁止它们运行
    // !!!!!!!!!!!!!!!!
    //  从现在开始，世界已停止(STW)...
    // !!!!!!!!!!!!!!!!
	// Finish sweep before we start concurrent scan.
	systemstack(func() {
		finishsweep_m()
	})

	// clearpools before we start the GC. If we wait they memory will not be
	// reclaimed until the next GC cycle.
	clearpools()

	work.cycles.Add(1)

	// 本轮gc控制器状态初始化
	gcController.startCycle(now, int(gomaxprocs), trigger)

	// Notify the CPU limiter that assists may begin.
	gcCPULimiter.startGCTransition(true, now)

	// In STW mode, disable scheduling of user Gs. This may also
	// disable scheduling of this goroutine, so it may block as
	// soon as we start the world again.
	if mode != gcBackgroundMode {
		schedEnableUser(false)
	}

	// Enter concurrent mark phase and enable
	// write barriers.
	//
	// Because the world is stopped, all Ps will
	// observe that write barriers are enabled by
	// the time we start the world and begin
	// scanning.
	//
	// Write barriers must be enabled before assists are
	// enabled because they must be enabled before
	// any non-leaf heap objects are marked. Since
	// allocations are blocked until assists can
	// happen, we want enable assists as early as
	// possible.
	setGCPhase(_GCmark)

	//做标记准备工作
	gcBgMarkPrepare() // Must happen before assist enable.
	gcMarkRootPrepare() //计算扫描根对象的任务数量

	// 标记所有活跃小对象
	gcMarkTinyAllocs()

	// At this point all Ps have enabled the write
	// barrier, thus maintaining the no white to
	// black invariant. Enable mutator assists to
	// put back-pressure on fast allocating
	// mutators.
	atomic.Store(&gcBlackenEnabled, 1)

	// In STW mode, we could block the instant systemstack
	// returns, so make sure we're not preemptible.
	mp = acquirem()

	// Concurrent mark.
	systemstack(func() {
		//  重新启动世界，并更新相关信息
		now = startTheWorldWithSema()
		work.pauseNS += now - work.pauseStart
		work.tMark = now
		memstats.gcPauseDist.record(now - work.pauseStart)

		sweepTermCpu := int64(work.stwprocs) * (work.tMark - work.tSweepTerm)
		work.cpuStats.gcPauseTime += sweepTermCpu
		work.cpuStats.gcTotalTime += sweepTermCpu

		// Release the CPU limiter.
		gcCPULimiter.finishGCTransition(now)
	})
    // !!!!!!!!!!!!!!!
    // 世界已重新启动...
    // !!!!!!!!!!!!!!!

	// Release the world sema before Gosched() in STW mode
	// because we will need to reacquire it later but before
	// this goroutine becomes runnable again, and we could
	// self-deadlock otherwise.
	semrelease(&worldsema)
	releasem(mp)

	// Make sure we block instead of returning to user code
	// in STW mode.
	if mode != gcBackgroundMode {
		Gosched()
	}

	semrelease(&work.startSema)
}
```



## 问题 
### 1. 生产速度（申请内存） 大于消费速度（GC）？

解决：在标记开始的时候，收集器会默认抢占 25% 的 CPU 性能，剩下的75%会分配给程序执行。
但是一旦收集器认为来不及进行标记任务了，就会改变这个 25% 的性能分配。这个时候收集器会抢占程序额外的 CPU，这部分被抢占 goroutine 有个名字叫 Mark Assist。
而且因为抢占 CPU 的目的主要是 GC 来不及标记新增的内存，那么抢占正在分配内存的 goroutine 效果会更加好，所以分配内存速度越快的 goroutine 就会被抢占越多的资源。





## 观察GC方式

1.
```go
 go build main.go
 GODEBUG=gctrace=1 ./main
```
GODEBUG 变量可以控制运行时内的调试变量，参数以逗号分隔，格式为：name=val。本文着重点在 GC 的观察上，主要涉及 gctrace 参数，
我们通过设置 gctrace=1 后就可以使得垃圾收集器向标准错误流发出 GC 运行信息

2.
```shell
go tool trace trace.out
```

3. debug.ReadGCStats

4. runtime.ReadMemStats

```shell
gc 18 @17.141s 0%: 0.21+4.8+0.007 ms clock, 1.7+0/0.45/4.8+0.062 ms cpu, 1->1->1 MB, 4 MB goal, 8 P (forced)
gc 19 @18.151s 0%: 0.063+1.1+0.003 ms clock, 0.51+0/0.12/1.1+0.030 ms cpu, 1->1->1 MB, 4 MB goal, 8 P (forced)
gc 20 @19.157s 0%: 0.12+3.5+0.008 ms clock, 0.98+0/3.6/0.094+0.064 ms cpu, 1->1->1 MB, 4 MB goal, 8 P (forced)
```

涉及术语
- mark：标记阶段。
- markTermination：标记结束阶段。
- mutator assist：辅助 GC，是指在 GC 过程中 mutator 线程会并发运行，而 mutator assist 机制会协助 GC 做一部分的工作。
- heap_live：在 Go 的内存管理中，span 是内存页的基本单元，每页大小为 8kb，同时 Go 会根据对象的大小不同而分配不同页数的 span，而 heap_live 就代表着所有 span 的总大小。
- dedicated / fractional / idle：在标记阶段会分为三种不同的 mark worker 模式，分别是 dedicated、fractional 和 idle，它们代表着不同的专注程度，
  其中 dedicated 模式最专注，是完整的 GC 回收行为，fractional 只会干部分的 GC 行为，idle 最轻松。
  这里你只需要了解它是不同专注程度的 mark worker 就好了，详细介绍我们可以等后续的文章。

```shell
gc # @#s #%: #+#+# ms clock, #+#/#/#+# ms cpu, #->#-># MB, # MB goal, # P
```

含义
- 1 gc#：GC 执行次数的编号，每次叠加。
- 2 @#s：自程序启动后到当前的具体秒数。
- 3 #%：自程序启动以来在GC中花费的时间百分比。
- 4 #+...+#：GC 的标记工作共使用的 CPU 时间占总 CPU 时间的百分比。
- 5 #->#-># MB：分别表示 GC 启动时, GC 结束时, GC 活动时的堆大小.
- 6 #MB goal：下一次触发 GC 的内存占用阈值。
- 7 #P：当前使用的处理器 P 的数量

案例代码
```go

const capacity = 50000

var d interface{}

func main() {
	//value值方式
	//d = Value()

	// value指针方式
	d = pointer()
	for i := 0; i < 20; i++ {
		runtime.GC()
		time.Sleep(time.Second)
	}

}

func Value() interface{} {
	m := make(map[int]int, capacity)
	for i := 0; i < capacity; i++ {
		m[i] = i
	}
	return m
}

func pointer() interface{} {
	m := make(map[int]*int, capacity)
	for i := 0; i < capacity; i++ {
		v := i
		m[i] = &v
	}
	return m
}
```
执行结果分析
```shell
gc 7 @0.140s 1%: 0.031+2.0+0.042 ms clock, 0.12+0.43/1.8/0.049+0.17 ms cpu, 4->4->1 MB, 5 MB goal, 4 P
```


- gc 7：第 7 次 GC。
- @0.140s：当前是程序启动后的 0.140s。
- 1%：程序启动后到现在共花费 1% 的时间在 GC 上。
- 0.031+2.0+0.042 ms clock：
  0.031：表示单个 P 在 mark 阶段的 STW 时间。
  2.0：表示所有 P 的 mark concurrent（并发标记）所使用的时间。
  0.042：表示单个 P 的 markTermination 阶段的 STW 时间。

- 0.12+0.43/1.8/0.049+0.17 ms cpu：
  0.12：表示整个进程在 mark 阶段 STW 停顿的时间。
  0.43/1.8/0.049：0.43 表示 mutator assist 占用的时间，1.8 表示 dedicated + fractional 占用的时间，0.049 表示 idle 占用的时间。
  0.17ms：0.17 表示整个进程在 markTermination 阶段 STW 时间。

- 4->4->1 MB：
  4：表示开始 mark 阶段前的 heap_live 大小。
  4：表示开始 markTermination 阶段前的 heap_live 大小。
  1：表示被标记对象的大小。
- 5 MB goal：表示下一次触发 GC 回收的阈值是 5 MB。
- 4 P：本次 GC 一共涉及多少个 P。

- wall clock 是指开始执行到完成所经历的实际时间，包括其他程序和本程序所消耗的时间；cpu time 是指特定程序使用 CPU 的时间；他们存在以下关系：

    - wall clock < cpu time: 充分利用多核

    - wall clock ≈ cpu time: 未并行执行

    - wall clock > cpu time: 多核优势不明显

![](.next_gc_images/backend_mark.png)
- DedicatedMode代表处理器专门负责标记对象，不会被调度器抢占；
- FractionalMode代表协助后台标记，其在整个标记阶段只会花费一定部分时间执行，
- IdleMode 为当处理器没有查找到可以执行的 协程时，执行垃圾收集的标记任务直到被抢占


## 参考

- [Golang 1.16.2 GC源码分析](https://blog.haohtml.com/archives/26358/)