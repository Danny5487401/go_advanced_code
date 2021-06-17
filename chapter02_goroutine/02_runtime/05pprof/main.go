package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

/*
Go 中监控代码性能的有两个包：
	net/http/pprof
		使用场景：在线服务（一直运行着的程序）
	runtime/pprof
		使用场景：工具型应用（比如说定制化的分析小工具、集成到公司监控系统）
这两个包都是可以监控代码性能的， 只不过net/http/pprof是通过http端口方式暴露出来的，内部封装的仍然是runtime/pprof。

介绍：
	runtime/pprof中的程序来生成三种包含实时性数据的概要文件，分别是
	1. CPU概要文件
		在默认情况下，Go语言的运行时系统会以100 Hz的的频率对CPU使用情况进行取样。
	2. 内存概要文件
		内存概要文件用于保存在用户程序执行期间的内存使用情况。这里所说的内存使用情况，其实就是程序运行过程中堆内存的分配情况。
	3. 程序阻塞概要文件
		程序阻塞概要文件用于保存用户程序中的Goroutine阻塞事件的记录。
*/

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
	for i := 0; i < 4; i++ {
		go do()
	}
	http.ListenAndServe("0.0.0.0:6060", nil)

	// 2 runtime/pprof 使用
	//file, err := os.Create("chapter02_goroutine/02_runtime/05pprof/cpu.pprof")
	//if err != nil {
	//	fmt.Printf("创建采集文件失败, err:%v\n", err)
	//	return
	//}
	//// 进行cpu数据的获取
	//pprof.StartCPUProfile(file)
	//defer pprof.StopCPUProfile()
	//
	//// 执行一段有问题的代码
	//for i := 0; i < 4; i++ {
	//	go do()
	//}
	//time.Sleep(10 * time.Second)
}

/*
展示参数：

	类型	描述
	allocs	内存分配情况的采样信息
	blocks	阻塞操作情况的采样信息
	cmdline	显示程序启动命令参数及其参数
	goroutine	显示当前所有协程的堆栈信息
	heap	堆上的内存分配情况的采样信息
	mutex	锁竞争情况的采样信息
	profile	cpu占用情况的采样信息，点击会下载文件
	threadcreate	系统线程创建情况的采样信息
	trace	程序运行跟踪信息

*/
