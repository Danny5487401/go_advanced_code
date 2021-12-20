package main

// runtime的基本使用

/*
一。 runtime 核心功能包括以下内容:
	1. 协程(goroutine)调度(并发调度模型)：linux的调度为CPU找到可运行的线程. 而Go的调度是为M(线程)找到P(内存, 执行票据)和可运行的G.
	2. 垃圾回收(GC)
	3. 内存分配
	4. 使得 golang 可以支持如 pprof、trace、race 的检测
	5. map, channel, string等内置类型及反射的实现.
	6. 操作系统及CPU相关的操作的封装(信号处理, 系统调用, 寄存器操作, 原子操作等), CGO;
二. 特点：
	1.与Java, Python不同, Go并没有虚拟机的概念, Runtime也直接被编译 成native code.
	2.go对系统调用的指令进行了封装, 可不依赖于glibc
	3. 用户代码与Runtime代码在执行的时候并没有明显的界限, 都是函数调用
	4. 一些go的关键字被编译器编译成runtime包下的函数.
		go-->newproc
		new->newobject
		make->makeslice,makechan,makemap,makemap_small
		<- 代表chansend1
		-> 代表chanrecv1

*/

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	//1. 获取GOROOT环境变量：
	fmt.Println("GOROOT-->", runtime.GOROOT()) //E:\go

	//2. 获取操作系统
	fmt.Println("os/platform-->", runtime.GOOS) // GOOS--> darwin，mac系统   windows

	//3.获取逻辑cpu的数量
	fmt.Println("逻辑CPU的核数：", runtime.NumCPU()) //4
	//4.设置最大可同时执行的最大CPU数：[1,256]
	n := runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("最大CPU数:", n)
	// 5. Gosched()
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Printf("goroutine。。。%d\n", i)
		}

	}()

	for i := 0; i < 4; i++ {
		//让出时间片，先让别的协议执行，它执行完，再回来执行此协程
		runtime.Gosched()
		fmt.Printf("main。%d\n", i)
	}
	// 6. 获取版本号
	fmt.Println(runtime.Version())
	// 7. 变量绑定方法,当垃圾回收的时候进行监听
	var i *Student = new(Student)
	runtime.SetFinalizer(i, func(i *Student) {
		println("垃圾回收了")
	})
	// 立即执行一次垃圾回收
	runtime.GC()
	time.Sleep(time.Second)
	// 8. 获取程序调用go协程的栈踪迹历史 func Stack(buf []byte, all bool) int
	//Stack将调用其的go程的调用栈踪迹格式化后写入到buf中并返回写入的字节数。
	//若all为true，函数会在写入当前go程的踪迹信息后，将其它所有go程的调用栈踪迹都格式化写入到buf中
	go showRecord()
	time.Sleep(time.Second)
	buf := make([]byte, 10000)
	runtime.Stack(buf, true)
	fmt.Println(string(buf))

	// 9.runtime.Caller函数可以获取当前函数的调用者列表
	// 获取当前函数或者上层函数的标识号、文件名、调用方法在当前文件中的行号
	show(2)
}

type Student struct {
	name string
}

func showRecord() {
	ticker := time.Tick(time.Second)
	for t := range ticker {
		fmt.Println(t)
	}
}

func show(depth int) {
	for skip := 0; skip < depth; skip++ {
		// runtime.Caller 返回值为 程序计数器，文件位置，行号，是否能恢复信息
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		// 根据PC寄存器表示的指令位置，通过runtime.FuncForPC函数获取函数的基本信息
		p := runtime.FuncForPC(pc)
		fnfile, fnline := p.FileLine(0)

		fmt.Printf("skip = %d, pc = 0x%08X\n", skip, pc)
		fmt.Printf("  func: file = %s, line = L%03d, name = %s, entry = 0x%08X\n", fnfile, fnline, p.Name(), p.Entry())
		fmt.Printf("  call: file = %s, line = L%03d\n", file, line)
	}
}
