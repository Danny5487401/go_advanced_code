package main

import (
	"fmt"
	"syscall"
)

func main() {
	// 函数查看当前进程 PID 的例子
	pid, _, _ := syscall.Syscall(39, 0, 0, 0) // 用不到的就补上 0
	fmt.Println("Process id: ", pid)
}
