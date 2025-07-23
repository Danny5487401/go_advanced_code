package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	// 由于是值复制的，run函数内部的wg和main函数内的wg不是同一个，会出现直接打印done的情况
	run(wg)
	wg.Wait()
	fmt.Println("done")
}

// 错误使用
func run(wg sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		go func(num int) {
			wg.Add(1)
			fmt.Println(num)
			wg.Done()
		}(i)
	}
}

/*
✗ go vet chapter09_design_pattern/06_Synchronization_mode/01_sync/11_nocpy/main.go
# command-line-arguments
# [command-line-arguments]
chapter09_design_pattern/06_Synchronization_mode/01_sync/11_nocpy/main.go:11:6: call of run copies lock value: sync.WaitGroup contains sync.noCopy
chapter09_design_pattern/06_Synchronization_mode/01_sync/11_nocpy/main.go:17:13: run passes lock by value: sync.WaitGroup contains sync.noCopy

*/
