package main

import (
	"log"
	"runtime"
	"time"
)

func main() {

	// Go1.23 之前错误生成大量timer：造成内存泄漏
	var queue = make(chan string)
	// 错误写法
	//useWrongTimeAfter(queue)
	// 修改写法
	useRightNewTimer(queue)
	time.Sleep(time.Second * 3)
	close(queue)
	traceMemStats()
	time.Sleep(time.Second * 2)
}

// 在3分钟内容易重复创建对象，底层并没有删除对象，造成内存泄漏
func useWrongTimeAfter(queue <-chan string) {
	Running := true
	for Running {
		select {
		case _, ok := <-queue:
			Running = false
			if !ok {
				return
			}

		case <-time.After(1 * time.Second):
			// 超时退出
			return
		}
	}
}

// 正确的方式:重复利用对象
func useRightNewTimer(in <-chan string) {
	idleDuration := 1 * time.Second
	idleDelay := time.NewTimer(idleDuration)
	defer idleDelay.Stop()
	Running := true
	for Running {
		idleDelay.Reset(idleDuration)
		select {
		case _, ok := <-in:
			Running = false
			if !ok {
				return
			}

		case <-idleDelay.C:
			// handle `s`
			return
		}
	}
}

func traceMemStats() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	log.Printf("Alloc:%d(bytes) HeapIdle:%d(bytes) HeapReleased:%d(bytes)", ms.Alloc, ms.HeapIdle, ms.HeapReleased)
	// Go1.23 之前
	// 2024/08/15 14:41:07 Alloc:102368(bytes) HeapIdle:3604480(bytes) HeapReleased:3571712(bytes) // 错误写法
	// 2024/08/15 14:41:35 Alloc:102384(bytes) HeapIdle:3514368(bytes) HeapReleased:3481600(bytes) // 修改写法

	// Go1.23
	// 2024/08/15 14:39:45 Alloc:127312(bytes) HeapIdle:3375104(bytes) HeapReleased:3342336(bytes) // 错误写法
	// 2024/08/15 14:39:26 Alloc:128208(bytes) HeapIdle:3473408(bytes) HeapReleased:3440640(bytes) // 修改写法
}
