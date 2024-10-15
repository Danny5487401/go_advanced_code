<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [GOMAXPROCS](#gomaxprocs)
  - [确定 p 的数量过程](#%E7%A1%AE%E5%AE%9A-p-%E7%9A%84%E6%95%B0%E9%87%8F%E8%BF%87%E7%A8%8B)
  - [参考资料](#%E5%8F%82%E8%80%83%E8%B5%84%E6%96%99)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# GOMAXPROCS

GOMAXPROCS 是 Golang 提供的非常重要的一个环境变量设定。通过设定 GOMAXPROCS，用户可以调整 Runtime Scheduler 中 Processor（简称P）的数量。
由于每个系统线程，必须要绑定 P 才能真正地进行执行。所以 P 的数量会很大程度上影响 Golang Runtime 的并发表现。
GOMAXPROCS 在目前版本（1.12）的默认值是 CPU 核数，而以 Docker 为代表的容器虚拟化技术，会通过 cgroup 等技术对 CPU 资源进行隔离。
以 Kubernetes 为代表的基于容器虚拟化实现的资源管理系统，也支持这样的特性。



## 确定 p 的数量过程
```go
// /Users/python/go/go1.18/src/runtime/runtime2.go
var (
	allm       *m
	gomaxprocs int32
	ncpu       int32
	// ...
)
```

linux 获取 cpu 数量：系统调用 sched_getaffinity 获取亲和性，一个进程的CPU亲合力掩码决定了该进程将在哪个或哪几个CPU上运行.
```go
// go1.18/src/runtime/os_linux.go
func osinit() {
    ncpu = getproccount()
	// ... 
}


func getproccount() int32 {
	// This buffer is huge (8 kB) but we are on the system stack
	// and there should be plenty of space (64 kB).
	// Also this is a leaf, so we're not holding up the memory for long.
	// See golang.org/issue/11823.
	// The suggested behavior here is to keep trying with ever-larger
	// buffers, but we don't have a dynamic memory allocator at the
	// moment, so that's a bit tricky and seems like overkill.
	const maxCPUs = 64 * 1024
	var buf [maxCPUs / 8]byte
	r := sched_getaffinity(0, unsafe.Sizeof(buf), &buf[0])
	if r < 0 {
		return 1
	}
	n := int32(0)
	for _, v := range buf[:r] {
		for v != 0 {
			n += int32(v & 1)
			v >>= 1
		}
	}
	if n == 0 {
		n = 1
	}
	return n
}
```


调度器初始化
```go
// go1.21.5/src/runtime/proc.go
func schedinit() {
    // 前面是一些锁的初始化，可以忽略
    // ..

	// raceinit must be the first call to race detector.
	// In particular, it must be done before mallocinit below calls racemapshadow.
	gp := getg() // 从线程本地存储中获取当前正在运行的 g
	if raceenabled {
		gp.racectx, raceprocctx0 = raceinit()
	}
    // 设置最多启动 10000 个操作系统线程，即 10000 个 M
	sched.maxmcount = 10000

	// The world starts stopped.
	worldStopped()

	moduledataverify()
	stackinit()
	mallocinit()
	godebug := getGodebugEarly()
	initPageTrace(godebug) // must run after mallocinit but before anything allocates
	cpuinit(godebug)       // must run before alginit
	alginit()              // maps, hash, fastrand must not be used before this call
	fastrandinit()         // must run before mcommoninit
	mcommoninit(gp.m, -1)
	modulesinit()   // provides activeModules
	typelinksinit() // uses maps, activeModules
	itabsinit()     // uses activeModules
	stkobjinit()    // must run before GC starts

	sigsave(&gp.m.sigmask)
	initSigmask = gp.m.sigmask

	goargs()
	goenvs()
	secure()
	parsedebugvars()
	gcinit()

	// if disableMemoryProfiling is set, update MemProfileRate to 0 to turn off memprofile.
	// Note: parsedebugvars may update MemProfileRate, but when disableMemoryProfiling is
	// set to true by the linker, it means that nothing is consuming the profile, it is
	// safe to set MemProfileRate to 0.
	if disableMemoryProfiling {
		MemProfileRate = 0
	}

	lock(&sched.lock)
	sched.lastpoll.Store(nanotime())
	procs := ncpu
	if n, ok := atoi32(gogetenv("GOMAXPROCS")); ok && n > 0 {
		procs = n
	}
	//  procresize 函数具体创建 P
	if procresize(procs) != nil {
		throw("unknown runnable goroutine during bootstrap")
	}
	unlock(&sched.lock)

	// World is effectively started now, as P's can run.
	worldStarted()

    // ..
}
```

创建的 P 
```go
// 具体创建 P
func procresize(nprocs int32) *p {
	assertLockHeld(&sched.lock)
	assertWorldStopped()

	old := gomaxprocs
    // ...

	maskWords := (nprocs + 31) / 32

	// Grow allp if necessary.
	if nprocs > int32(len(allp)) {
		// Synchronize with retake, which could be running
		// concurrently since it doesn't run on a P.
		lock(&allpLock)
		if nprocs <= int32(cap(allp)) {
			// 如果需要的 P 个数小于 allp 的容量，缩小 allp 容量
			allp = allp[:nprocs]
		} else {
			// /否则，创建临时 nallp 数组，并 allp 全部拷贝至 nallp
			nallp := make([]*p, nprocs)
			// Copy everything up to allp's cap so we
			// never lose old allocated Ps.
			copy(nallp, allp[:cap(allp)])
			// 更新 allp
			allp = nallp
		}

		if maskWords <= int32(cap(idlepMask)) {
			idlepMask = idlepMask[:maskWords]
			timerpMask = timerpMask[:maskWords]
		} else {
			nidlepMask := make([]uint32, maskWords)
			// No need to copy beyond len, old Ps are irrelevant.
			copy(nidlepMask, idlepMask)
			idlepMask = nidlepMask

			ntimerpMask := make([]uint32, maskWords)
			copy(ntimerpMask, timerpMask)
			timerpMask = ntimerpMask
		}
		unlock(&allpLock)
	}

	// initialize new P's
	for i := old; i < nprocs; i++ {
		pp := allp[i]
		if pp == nil {
			// 申请新 P 对象
			pp = new(p)
		}
		// 初始化 P，每个 P 都有一个ID，该ID就是 P 在 allp 的索引
		pp.init(i)
		atomicstorep(unsafe.Pointer(&allp[i]), unsafe.Pointer(pp))
	}

	gp := getg()
	if gp.m.p != 0 && gp.m.p.ptr().id < nprocs {
		// continue to use the current P
		gp.m.p.ptr().status = _Prunning
		gp.m.p.ptr().mcache.prepareForSweep()
	} else {
		// release the current P and acquire allp[0].
		//
		// We must do this before destroying our current P
		// because p.destroy itself has write barriers, so we
		// need to do that from a valid P.
		if gp.m.p != 0 {
			//  继续使用当前的 P
			gp.m.p.ptr().m = 0
		}
		// 取第 0 号 p，清空m, 并设置p状态为 _Pidle
		gp.m.p = 0
		pp := allp[0]
		pp.m = 0
		pp.status = _Pidle
		acquirep(pp)
		if traceEnabled() {
			traceGoStart()
		}
	}

	// g.m.p is now set, so we no longer need mcache0 for bootstrapping.
	mcache0 = nil

	// 释放不在使用的 P 的资源，但不会收回 P，因为有可能处于系统调用的M还会关联 P
	for i := nprocs; i < old; i++ {
		pp := allp[i]
		pp.destroy()
		// can't free P itself because it can be referenced by an M in syscall
	}

	// Trim allp.
	if int32(len(allp)) != nprocs {
		lock(&allpLock)
		allp = allp[:nprocs]
		idlepMask = idlepMask[:maskWords]
		timerpMask = timerpMask[:maskWords]
		unlock(&allpLock)
	}

	// 下面这个 for 循环检查 allp，将 空闲的P 放入全局空闲P列表；
	var runnablePs *p
	for i := nprocs - 1; i >= 0; i-- {
		pp := allp[i]
		if gp.m.p.ptr() == pp {
			// allp[i] 已经跟 m0 关联了，所以要跳过
			continue
		}
		pp.status = _Pidle
		if runqempty(pp) {
			pidleput(pp, now)
		} else {
			pp.m.set(mget())
			pp.link.set(runnablePs)
			runnablePs = pp
		}
	}
	stealOrder.reset(uint32(nprocs))
	var int32p *int32 = &gomaxprocs // make compiler check that gomaxprocs is an int32
	atomic.Store((*uint32)(unsafe.Pointer(int32p)), uint32(nprocs))
	if old != nprocs {
		// Notify the limiter that the amount of procs has changed.
		gcCPULimiter.resetCapacity(now, nprocs)
	}
	// 返回所有可运行的 P 列表
	return runnablePs
}
```

procresize函数主要做了以下工作：

1. 根据传入的 nproc 参数，创建这么多 P，要么复用原来的P，要么新建
2. 初次启动时，该函数完成了 M0 与 P0 的关联
3. 检查 allp 列表，将所有空闲的 P 放入全局空闲 P 列表，返回所有可运行的 P 列表



## 参考资料
- [GOMAXPROCS 与容器的相处之道](https://gaocegege.com/Blog/maxprocs-cpu)