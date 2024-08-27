<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [time.sleep](#timesleep)
  - [源码](#%E6%BA%90%E7%A0%81)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# time.sleep

## 源码
```go
// go1.21.5/src/runtime/runtime2.go

//go:linkname timeSleep time.Sleep
func timeSleep(ns int64) {
	if ns <= 0 { //判断入参是否正常
		return
	}

	gp := getg() //获取当前的goroutine
	t := gp.timer
	if t == nil { //如果不存在timer，new一个
		t = new(timer)
		gp.timer = t
	}
	t.f = goroutineReady //后面唤醒时候会用到，修改goroutine状态为goready
	t.arg = gp
	t.nextwhen = nanotime() + ns //记录唤醒时间
	if t.nextwhen < 0 { // check for overflow.
		t.nextwhen = maxWhen
	}
    //调用gopark挂起goroutine
	gopark(resetForSleep, unsafe.Pointer(t), waitReasonSleep, traceBlockSleep, 1)
}
```

当触发完gopark方法，会调用releasem(mp)方法释放当前goroutine与m的连接后，该goroutine脱离当前的m挂起，进入gwaiting状态，不在任何运行队列上。



## 参考

[图解 Go语言 time.Sleep 的实现原理](https://mp.weixin.qq.com/s/02w-k5YgYxMC_gxbRdNpJQ)