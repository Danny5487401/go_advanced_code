package main

// 普通的map，map保存到是int到int的映射，会执行delete删除map的每一项，执行垃圾回收，看内存是否被回收，map设置为nil，再看是否被回收

import (
	"log"
	"runtime"
)

var lastTotalFreed uint64
var intMap map[int]int
var cnt = 8192

func main() {
	printMemStats()

	initMap()
	runtime.GC()
	printMemStats()

	log.Println(len(intMap))
	for i := 0; i < cnt; i++ {
		delete(intMap, i)
	}
	log.Println(len(intMap))

	runtime.GC()
	printMemStats()

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

func printMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("Alloc = %v TotalAlloc = %v  Just Freed = %v Sys = %v NumGC = %v\n",
		m.Alloc/1024, m.TotalAlloc/1024, ((m.TotalAlloc-m.Alloc)-lastTotalFreed)/1024, m.Sys/1024, m.NumGC)

	lastTotalFreed = m.TotalAlloc - m.Alloc
}

/*
Alloc：当前堆上对象占用的内存大小。
TotalAlloc：堆上总共分配出的内存大小。
Sys：程序从操作系统总共申请的内存大小。
NumGC：垃圾回收运行的次数。
结论：
	Alloc代表了map占用的内存大小，这个结果表明，执行完delete后，map占用的内存并没有变小，Alloc依然是387，代表map的key和value占用的空间仍在map里.
	执行完map设置为nil，Alloc变为74，与刚创建的map大小基本是约等于。
*/
