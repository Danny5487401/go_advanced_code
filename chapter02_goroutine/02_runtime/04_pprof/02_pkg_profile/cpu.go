package main

import (
	"time"

	"github.com/pkg/profile"
)

func joinSlice1() []string {

	const count = 100000

	var arr []string = make([]string, count)

	for i := 0; i < count; i++ {
		arr[i] = "arr"
	}

	return arr
}

func joinSlice2() []string {

	var arr []string

	for i := 0; i < 100000; i++ {
		// 故意造成多次的切片添加(append)操作, 由于每次操作可能会有内存重新分配和移动, 性能较低
		arr = append(arr, "arr")
	}

	return arr
}

func main() {
	// 开始性能分析, 返回一个停止接口
	stopper := profile.Start(profile.CPUProfile, profile.ProfilePath("."))

	// 在main()结束时停止性能分析
	defer stopper.Stop()

	// 分析的核心逻辑
	joinSlice2()

	// 让程序至少运行1秒
	time.Sleep(time.Second)
}

/*
# 第 1 行将 cpu.go 编译为可执行文件 cpu。
$ go build -o cpu cpu.go

# 第 2 行运行可执行文件，在当前目录输出 cpu.pprof 文件。
$ ./cpu

# 第 3 行，使用 go tool 工具链输入 cpu.pprof 和 cpu 可执行文件，生成 PDF 格式的输出文件，将输出文件重定向为 cpu.pdf 文件。
	这个过程中会调用 Graphviz 工具，Windows 下需将 Graphviz 的可执行目录添加到环境变量 PATH
$ go tool pprof --pdf cpu cpu.pprof > cpu.pdf



*/
