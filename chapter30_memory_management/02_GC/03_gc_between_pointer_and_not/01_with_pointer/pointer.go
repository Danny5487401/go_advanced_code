package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	a := make([]*int, 1e9)

	for i := 0; i < 10; i++ {
		start := time.Now()
		runtime.GC()
		fmt.Printf("GC took %s\n", time.Since(start))
	}

	// 程序中调用runtime.KeepAlive函数用于保证在该函数调用点之前切片a不会被GC释放掉。
	runtime.KeepAlive(a)
}
