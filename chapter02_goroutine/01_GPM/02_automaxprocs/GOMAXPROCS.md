<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [GOMAXPROCS](#gomaxprocs)
  - [确定 p 的数量过程](#%E7%A1%AE%E5%AE%9A-p-%E7%9A%84%E6%95%B0%E9%87%8F%E8%BF%87%E7%A8%8B)
  - [参考资料](#%E5%8F%82%E8%80%83%E8%B5%84%E6%96%99)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# GOMAXPROCS

GOMAXPROCS 是 Golang 提供的非常重要的一个环境变量设定。通过设定 GOMAXPROCS，用户可以调整 Runtime Scheduler 中 Processor（简称P）的数量。
由于每个系统线程，必须要绑定 P 才能真正地进行执行。所以 P 的数量会很大程度上影响 Golang Runtime 的并发表现。GOMAXPROCS 在目前版本（1.12）的默认值是 CPU 核数，而以 Docker 为代表的容器虚拟化技术，会通过 cgroup 等技术对 CPU 资源进行隔离。以 Kubernetes 为代表的基于容器虚拟化实现的资源管理系统，也支持这样的特性。



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
/Users/python/go/go1.18/src/runtime/os_linux.go
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
// The bootstrap sequence is:
//
//	call osinit
//	call schedinit
//	make & queue new G
//	call runtime·mstart
//
// The new G calls runtime·main.
func schedinit() {
    procs := ncpu
    if n, ok := atoi32(gogetenv("GOMAXPROCS")); ok && n > 0 {
        procs = n
    }
}
```




## 参考资料
- [GOMAXPROCS 与容器的相处之道](https://gaocegege.com/Blog/maxprocs-cpu)