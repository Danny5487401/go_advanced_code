package main

import (
	"fmt"
	"os"
	"runtime/trace"
)

// trace侧重于分析goroutine的调度
func main() {
	// 创建trace文件
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 启动trace goroutine
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	// main
	fmt.Println("Hello trace")
}

// 运行
//go run trace.go
// go tool trace trace.out
