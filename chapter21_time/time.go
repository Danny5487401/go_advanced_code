package main

import (
	"fmt"
	"time"
	_ "unsafe"
)

//大部分场景下，我们只需要时间戳
func main() {
	// 常规方式:涉及两次系统调用
	//runtime.walltime 和runtime.nanotime
	now := time.Now()
	fmt.Printf("纳秒是%v\n", now.UnixNano())

	// 更优雅的方式:
	sec, nsec := WallTime()
	fmt.Printf("秒是%v,纳秒是%v\n", sec, nsec)

}

//go:linkname WallTime runtime.walltime
func WallTime() (sec int64, nsec int32)

// Note: 系统调用在 go 里面相对来讲是比较重的。runtime会切换到g0栈中去执行这部分代码，time.Now方法在go<=1.16中有两次连续的系统调用。
//不过，go 官方团队的 lan 大佬已经发现并提交优化**pr**。优化后，这两次系统调将会合并在一起，减少一次g0栈的切换。
