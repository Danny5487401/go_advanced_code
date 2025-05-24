package main

import (
	"log"
	"runtime"
)

/*
空间收缩

	map 不会收缩 “不再使用” 的空间。就算把所有键值删除，它依然保留内存空间以待后用

runtime.ReadMemStats(&m)

	直接通过运行时的内存相关的 API 进行监控
*/
var lastTotalFreed uint64
var intMap map[int]int
var cnt = 8192

func main() {
	// 1. 初始化状态
	printMemStats()

	// 2. 添加数据
	initMap()
	runtime.GC()
	printMemStats()

	log.Println(len(intMap))
	// 3. 删除数据
	for i := 0; i < cnt; i++ {
		delete(intMap, i)
	}
	log.Println("元素数量", len(intMap))
	runtime.GC()
	printMemStats()

	// 4. 释放map对象
	intMap = nil
	runtime.GC()
	printMemStats()
}

func initMap() {
	intMap = make(map[int]int, cnt)

	for i := 0; i < cnt; i++ {
		intMap[i] = i
	}
}

// 打印内存状态
func printMemStats() {
	var m runtime.MemStats
	// 查看内存申请和分配统计信息
	runtime.ReadMemStats(&m)
	log.Printf("Alloc = %v TotalAlloc = %v  Just Freed = %v Sys = %v NumGC = %v\n",
		m.Alloc/1024, m.TotalAlloc/1024, ((m.TotalAlloc-m.Alloc)-lastTotalFreed)/1024, m.Sys/1024, m.NumGC)

	lastTotalFreed = m.TotalAlloc - m.Alloc
}

/*
术语解释
	Alloc：当前堆上对象占用的内存大小。
	TotalAlloc：堆上总共分配出的内存大小。
	Sys：程序从操作系统总共申请的内存大小。
	NumGC：垃圾回收运行的次数。
结论：
	Alloc代表了map占用的内存大小，这个结果表明，执行完delete后，map占用的内存并没有变小，Alloc依然是387，代表map的key和value占用的空间仍在map里.
	执行完map设置为nil，Alloc变为74，与刚创建的map大小基本是约等于。
提示：如长期使用 map 对象（比如用作 cache 容器），偶尔换成 “新的” 或许会更好。还有，int key 要比 string key 更快。

*/
