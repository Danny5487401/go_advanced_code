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
	printMemStats()

	// 添加数据
	initMap()
	runtime.GC()
	printMemStats()

	log.Println(len(intMap))
	// 删除数据
	for i := 0; i < cnt; i++ {
		delete(intMap, i)
	}
	log.Println(len(intMap))

	runtime.GC()
	printMemStats()

	// 释放map对象
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
源码注释
	type MemStats struct {
		// 一般统计
		Alloc      uint64 // 已申请且仍在使用的字节数
		TotalAlloc uint64 // 已申请的总字节数（已释放的部分也算在内）
		Sys        uint64 // 从系统中获取的字节数（下面XxxSys之和）
		Lookups    uint64 // 指针查找的次数
		Mallocs    uint64 // 申请内存的次数
		Frees      uint64 // 释放内存的次数
		// 主分配堆统计
		HeapAlloc    uint64 // 已申请且仍在使用的字节数
		HeapSys      uint64 // 从系统中获取的字节数
		HeapIdle     uint64 // 闲置span中的字节数
		HeapInuse    uint64 // 非闲置span中的字节数
		HeapReleased uint64 // 释放到系统的字节数
		HeapObjects  uint64 // 已分配对象的总个数
		// L低层次、大小固定的结构体分配器统计，Inuse为正在使用的字节数，Sys为从系统获取的字节数
		StackInuse  uint64 // 引导程序的堆栈
		StackSys    uint64
		MSpanInuse  uint64 // mspan结构体
		MSpanSys    uint64
		MCacheInuse uint64 // mcache结构体
		MCacheSys   uint64
		BuckHashSys uint64 // profile桶散列表
		GCSys       uint64 // GC元数据
		OtherSys    uint64 // 其他系统申请
		// 垃圾收集器统计
		NextGC       uint64 // 会在HeapAlloc字段到达该值（字节数）时运行下次GC
		LastGC       uint64 // 上次运行的绝对时间（纳秒）
		PauseTotalNs uint64
		PauseNs      [256]uint64 // 近期GC暂停时间的循环缓冲，最近一次在[(NumGC+255)%256]
		NumGC        uint32
		EnableGC     bool
		DebugGC      bool
		// 每次申请的字节数的统计，61是C代码中的尺寸分级数
		BySize [61]struct {
			Size    uint32
			Mallocs uint64
			Frees   uint64
		}
	}
*/
