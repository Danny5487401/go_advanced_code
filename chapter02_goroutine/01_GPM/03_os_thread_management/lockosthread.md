<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [线程管理](#%E7%BA%BF%E7%A8%8B%E7%AE%A1%E7%90%86)
  - [背景](#%E8%83%8C%E6%99%AF)
  - [runtime.LockOSThread](#runtimelockosthread)
  - [runtime.UnlockOSThread()](#runtimeunlockosthread)
  - [应用-->cni](#%E5%BA%94%E7%94%A8--cni)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 线程管理

Go 语言的运行时会通过调度器改变线程的所有权，它也提供了 runtime.LockOSThread 和 runtime.UnlockOSThread 让我们有能力绑定 Goroutine 和线程完成一些比较特殊的操作。Goroutine 应该在调用操作系统服务或者依赖线程状态的非 Go 语言库时调用 runtime.LockOSThread 函数11，例如：C 语言图形库等

## 背景

golang的scheduler可以理解为公平协作调度和抢占的综合体，他不支持优先级调度。
当你开了几十万个goroutine，并且大多数协程已经在runq等待调度了, 那么如果你有一个重要的周期性的协程需要优先执行该怎么办？


可以借助runtime.LockOSThread()方法来绑定线程，绑定线程M后的好处在于，他可以由system kernel内核来调度，因为他本质是线程了


runtime.LockOSThread会锁定当前协程只跑在一个系统线程上，这个线程里也只跑该协程。


## runtime.LockOSThread 

```go
// go1.21.5/src/runtime/proc.go

// 调用 LockOSThread 将 绑定 当前 goroutine 到当前 操作系统 线程，此 goroutine 将始终在此线程执行，其它 goroutine 则无法在此线程中得到执行，
// 直到当前调用线程执行了 UnlockOSThread 为止（也就是说 LockOSThread 可以指定一个goroutine 独占 一个系统线程）；
// 如果调用者goroutine 在未解锁线程（未调用 UnlockOSThread）之前直接退出，则当前线程将直接被终止（也就是说线程被直接销毁）。
//
// 所有 init函数 都运行在启动线程。如果在一个 init函数 中调用了 LockOSThread 则导致 main 函数被执行在当前线程
//
// goroutine应该在调用依赖于每个线程状态的 OS服务 或 非Go库函数 之前调用 LockOSThread。
//
//go:nosplit
func LockOSThread() {
	if atomic.Load(&newmHandoff.haveTemplateThread) == 0 && GOOS != "plan9" {
		// If we need to start a new thread from the locked
		// thread, we need the template thread. Start it now
		// while we're in a known-good state.
		startTemplateThread()
	}
	gp := getg()
	gp.m.lockedExt++
	if gp.m.lockedExt == 0 {
		gp.m.lockedExt--
		panic("LockOSThread nesting overflow")
	}
	dolockOSThread()
}


//go:nosplit
func dolockOSThread() {
	if GOARCH == "wasm" {
		return // no threads on wasm yet
	}
	gp := getg()
	gp.m.lockedg.set(gp)
	gp.lockedm.set(gp.m)
}

```

runtime.dolockOSThread 会分别设置线程的 lockedg 字段和 Goroutine 的 lockedm 字段，这两行代码会绑定线程和 Goroutine。


## runtime.UnlockOSThread()

```go
//go:nosplit
func UnlockOSThread() {
	gp := getg()
	if gp.m.lockedExt == 0 {
		return
	}
	gp.m.lockedExt--
	dounlockOSThread()
}

//go:nosplit
func dounlockOSThread() {
	if GOARCH == "wasm" {
		return // no threads on wasm yet
	}
	gp := getg()
	if gp.m.lockedInt != 0 || gp.m.lockedExt != 0 {
		return
	}
	gp.m.lockedg = 0
	gp.lockedm = 0
}
```


## 应用-->cni

初始化调度
```go
// https://github.com/AliyunContainerService/terway/blob/c742f76b042a96949aadc8bd4610a2fb5aa0fead/plugin/terway/cni.go
func init() {
	runtime.LockOSThread()
}

func main() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.PluginSupports("0.3.0", "0.3.1", "0.4.0", "1.0.0"), bv.BuildString("terway"))
}
```

背景:在 Linux 中，不同的操作系统线程可能会设置不同的网络命名空间，而 Go 语言的协程会基于操作系统线程的负载以及其他信息动态地在不同的操作系统线程之间切换，这样可能会导致 Go 协程在意想不到的情况下切换到不同的网络命名空间中。
比较稳妥的做法是，利用 Go 语言提供的 runtime.LockOSThread() 函数保证特定的 Go 协程绑定到当前的操作系统线程中。

接口
```go
type NetNS interface {
	// Executes the passed closure in this object's network namespace,
	// attempting to restore the original namespace before returning.
	// However, since each OS thread can have a different network namespace,
	// and Go's thread scheduling is highly variable, callers cannot
	// guarantee any specific namespace is set unless operations that
	// require that namespace are wrapped with Do().  Also, no code called
	// from Do() should call runtime.UnlockOSThread(), or the risk
	// of executing code in an incorrect namespace will be greater.  See
	// https://github.com/golang/go/wiki/LockOSThread for further details.
	Do(toRun func(NetNS) error) error

	// Sets the current network namespace to this object's network namespace.
	// Note that since Go's thread scheduling is highly variable, callers
	// cannot guarantee the requested namespace will be the current namespace
	// after this function is called; to ensure this wrap operations that
	// require the namespace with Do() instead.
	Set() error

	// Returns the filesystem path representing this object's network namespace
	Path() string

	// Returns a file descriptor representing this object's network namespace
	Fd() uintptr

	// Cleans up this instance of the network namespace; if this instance
	// is the last user the namespace will be destroyed
	Close() error
}
```
具体实现
```go
// github.com/containernetworking/plugins/pkg/ns/ns_linux.go
func (ns *netNS) Do(toRun func(NetNS) error) error {
    // ..

	containedCall := func(hostNS NetNS) error {
		threadNS, err := GetCurrentNS()
		if err != nil {
			return fmt.Errorf("failed to open current netns: %v", err)
		}
		defer threadNS.Close()

		// switch to target namespace
		if err = ns.Set(); err != nil {
			return fmt.Errorf("error switching to ns %v: %v", ns.file.Name(), err)
		}
		defer func() {
			err := threadNS.Set() // switch back
			if err == nil {
				// Unlock the current thread only when we successfully switched back
				// to the original namespace; otherwise leave the thread locked which
				// will force the runtime to scrap the current thread, that is maybe
				// not as optimal but at least always safe to do.
				runtime.UnlockOSThread()
			}
		}()

		return toRun(hostNS)
	}

	// save a handle to current network namespace
	hostNS, err := GetCurrentNS()
	if err != nil {
		return fmt.Errorf("Failed to open current namespace: %v", err)
	}
	defer hostNS.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	// Start the callback in a new green thread so that if we later fail
	// to switch the namespace back to the original one, we can safely
	// leave the thread locked to die without a risk of the current thread
	// left lingering with incorrect namespace.
	var innerError error
	go func() {
		defer wg.Done()
		runtime.LockOSThread()
		innerError = containedCall(hostNS)
	}()
	wg.Wait()

	return innerError
}

```

调用runtime.LockOSThread()函数，锁定了执行闭包的当前线程。



```go
func setupVeth(netns ns.NetNS, br *netlink.Bridge, ifName string, ipAddress string) error {
  hostIface := &current.Interface{}                           ①

  err := netns.Do(func(hostNS ns.NetNS) error {               ②
    hostVeth, containerVeth, err := 
      ip.SetupVeth(ifName, 1500, hostNS)                      ③
    if err != nil {
      return err
    }
    hostIface.Name = hostVeth.Name                            ④

    return nil
  })
  if err != nil {
    return err
  }

  hostVeth, err := netlink.LinkByName(hostIface.Name)         ⑤
  if err != nil {
    return err
  }

  if err := netlink.LinkSetMaster(hostVeth, br); err != nil { ⑥
    return err
  }

  return nil
}
```

我们使用②时, Do函数在执行闭包之前，会先保存当前线程的network namespace，也就是宿主机对应的namespace，并作为参数传入闭包。然后在闭包执行结束时，把当前network namespace恢复成宿主机对应的namespace

## 参考




