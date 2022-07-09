package main

import (
	"fmt"
	"github.com/Danny5487401/go_advanced_code/chapter02_goroutine/02_runtime/05pprof/01_pprof/data"
	"log"
	"net/http"
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
	httpPprof()

	// 2 runtime/pprof 使用  适用于应用程序
	//runtimePprof()
}
func httpPprof() {
	go func() {
		for {
			time.Sleep(time.Second * 2)
			log.Println(data.Add("https://github.com/Danny"))
		}
	}()

	http.ListenAndServe("0.0.0.0:6060", nil)
}

func runtimePprof() {
	// go pprof 工具链配合 Graphviz 图形化工具可以将 runtime.pprof 包生成的数据转换为 PDF 格式，以图片的方式展示程序的性能分析结果
	// go tool pprof -pdf cpu.pprof 打印出pdf格式的文件

	file, err := os.Create("chapter02_goroutine/02_runtime/05pprof/01_pprof")
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
