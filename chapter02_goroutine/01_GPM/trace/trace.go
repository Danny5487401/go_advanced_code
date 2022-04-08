package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"sync"
	"time"
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
	urls := []string{"0.0.0.0:5000", "0.0.0.0:5001", "0.0.0.0:5002"}
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mockSendToServer(url)
		}()

	}
	time.Sleep(time.Microsecond)
	wg.Wait()
}

func mockSendToServer(url string) {
	fmt.Printf("url是%v\n", url)
}

/*
方式一：
1. 运行 生成文件
	go run trace.go
2. 文件进行分析
	go tool trace trace.out
会打印
2021/12/10 17:55:21 Parsing trace...
2021/12/10 17:55:21 Splitting trace...
2021/12/10 17:55:21 Opening browser. Trace viewer is listening on http://127.0.0.1:57569


方式二：
1. go build .
2. GODEBUG=schedtrace=1000 ./trace
// 再看个复杂版本的，加上scheddetail=1可以打印更详细的trace信息 GODEBUG=schedtrace=1000,scheddetail=1 ./trace
*/
