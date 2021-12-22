package main

import (
	"fmt"
	_ "net/http/pprof" //引入init就注册了路由
	"os"
	"runtime/pprof"
	"time"
)

// 一段有问题的代码
func do() {
	var c chan int
	for {
		select {
		case v := <-c:
			fmt.Printf("我是有问题的那一行，因为收不到值：%v", v)
		default:
		}
	}
}

func main() {
	// 1. net/http.pprof在线使用
	// http://localhost:6060/debug/pprof/ 查看信息
	// 执行一段有问题的代码
	//for i := 0; i < 4; i++ {
	//	go do()
	//}
	//http.ListenAndServe("0.0.0.0:6060", nil)

	// 2 runtime/pprof 使用  适用于应用程序
	// go pprof 工具链配合 Graphviz 图形化工具可以将 runtime.pprof 包生成的数据转换为 PDF 格式，以图片的方式展示程序的性能分析结果
	// go tool pprof -pdf cpu.pprof 打印出pdf格式的文件
	file, err := os.Create("chapter02_goroutine/02_runtime/05pprof/cpu.pprof")
	if err != nil {
		fmt.Printf("创建采集文件失败, err:%v\n", err)
		return
	}
	// 进行cpu数据的获取
	_ = pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

	// 执行一段有问题的代码
	for i := 0; i < 4; i++ {
		go do()
	}
	// runtime.pprof 包在运行时对程序进行每秒 100 次的采样，最少采样 1 秒
	time.Sleep(1 * time.Second)
}

/*
展示参数：
	类型	描述

	cmdline	显示程序启动命令参数及其参数
	profile	cpu占用情况的采样信息，点击会下载文件
	trace	程序运行跟踪信息

源码：
	profiles.m = map[string]*Profile{
		"goroutine":    goroutineProfile,  //显示当前所有协程的堆栈信息
		"threadcreate": threadcreateProfile, // 系统线程创建情况的采样信息
		"heap":         heapProfile,  // 堆上的内存分配情况的采样信息
		"allocs":       allocsProfile,  //内存分配情况的采样信息
		"block":        blockProfile,  //阻塞操作情况的采样信息
		"mutex":        mutexProfile,  // 锁竞争情况的采样信息
	}
*/
