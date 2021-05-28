package main

// 1. 数据结构及内存管理
//hashmap的定义位于 src/runtime/hashmap.go 中
/*
type hmap struct {
    count     int    // 元素的个数
    flags     uint8  // 状态标志
    B         uint8  // 可以最多容纳 6.5 * 2 ^ B 个元素，6.5为装载因子
    noverflow uint16 // 溢出的个数
    hash0     uint32 // 哈希种子

    buckets    unsafe.Pointer // 桶的地址
    oldbuckets unsafe.Pointer // 旧桶的地址，用于扩容
    nevacuate  uintptr        // 搬迁进度，小于nevacuate的已经搬迁
    overflow *[2]*[]*bmap
}

	其中，overflow是一个指针，指向一个元素个数为2的数组，数组的类型是一个指针，指向一个slice，slice的元素是桶(bmap)的地址，这些桶都是溢出桶；
		为什么有两个？因为Go map在hash冲突过多时，会发生扩容操作，为了不全量搬迁数据，使用了增量搬迁，[0]表示当前使用的溢出桶集合，
		[1]是在发生扩容时，保存了旧的溢出桶集合；overflow存在的意义在于防止溢出桶被gc。


// A bucket for a Go map.
type bmap struct {
    // 每个元素hash值的高8位，如果tophash[0] < minTopHash，表示这个桶的搬迁状态
    tophash [bucketCnt]uint8
    // 接下来是8个key、8个value，但是我们不能直接看到；为了优化对齐，go采用了key放在一起，value放在一起的存储方式，
    // 再接下来是hash冲突发生时，下一个溢出桶的地址
}
	tophash的存在是为了快速试错，毕竟只有8位，比较起来会快一点。

从定义可以看出，不同于STL中map以红黑树实现的方式，Golang采用了HashTable的实现，解决冲突采用的是链地址法。也就是说，使用数组+链表来实现map


 */

func main() {
	var m = map[int]int{}
	m[43] = 1
	var n = map[string]int{}
	n["abc"] = 1
	println(m, n)
}
