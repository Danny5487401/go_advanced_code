# LockOSThread


大部分的用户代码并不需要线程级的操作。但某些情况下，当需要 使用 cgo 调用 C 端图形库（如 GLib）时，甚至需要将某个 Goroutine 用户态代码一直在主线程上执行。

runtime提供了一个LockOSThread的函数，该方法的作用是可以让当前协程绑定并独立一个线程 M。


我们知道golang的scheduler可以理解为公平协作调度和抢占的综合体，他不支持优先级调度。当你开了几十万个goroutine，并且大多数协程已经在runq等待调度了, 那么如果你有一个重要的周期性的协程需要优先执行该怎么办？


可以借助runtime.LockOSThread()方法来绑定线程，绑定线程M后的好处在于，他可以由system kernel内核来调度，因为他本质是线程了。



## 第三方应用--flannel 源码

```go
// https://github.com/flannel-io/flannel/blob/aba86a5775f9ce103a3fc5958da3e31dafd5cf50/pkg/ns/ns.go

//go:build !windows
// +build !windows


package ns

import (
	"runtime"
	"testing"

	"github.com/vishvananda/netns"
)

func SetUpNetlinkTest(t *testing.T) func() {
	// new temporary namespace so we don't pollute the host
	// lock thread since the namespace is thread local
	runtime.LockOSThread()
	var err error
	ns, err := netns.New()
	if err != nil {
		t.Fatalf("Failed to create newns: %v", err)
	}

	return func() {
		ns.Close()
		runtime.UnlockOSThread()
	}
}
```


## 参考链接
1. [官方lock child goroutines to same OS thread ](https://github.com/golang/go/issues/23758)
2. [Goroutines, threads, and thread IDs](https://dunglas.dev/2022/05/goroutines-threads-and-thread-ids/)