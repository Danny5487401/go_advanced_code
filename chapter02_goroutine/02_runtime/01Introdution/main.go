package main
// runtime的基本使用

/*
golang 的 runtime 核心功能包括以下内容:

1. 协程(goroutine)调度(并发调度模型)
2. 垃圾回收(GC)
3. 内存分配
4. 使得 golang 可以支持如 pprof、trace、race 的检测
5. 支持 golang 的内置类型 channel、map_test、slice、string等的实现

 */

import (
	"fmt"
	"runtime"
)

func main()  {
	//1. 获取goroot目录：
	fmt.Println("GOROOT-->",runtime.GOROOT()) //E:\go

	//2. 获取操作系统
	fmt.Println("os/platform-->",runtime.GOOS) // GOOS--> darwin，mac系统   windows

	//3.获取逻辑cpu的数量
	fmt.Println("逻辑CPU的核数：",runtime.NumCPU()) //4
	//4.设置go程序执行的最大的：[1,256]
	n := runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(n)
	// 5. Gosched()
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Printf("goroutine。。。%d\n",i)
		}

	}()

	for i := 0; i < 4; i++ {
		//让出时间片，先让别的协议执行，它执行完，再回来执行此协程
		runtime.Gosched()
		fmt.Printf("main。%d\n",i)
	}

}
