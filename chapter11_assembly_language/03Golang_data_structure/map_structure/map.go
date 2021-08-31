package main

import (
	"fmt"
	"sort"
)

/*
1. map哈希表源码解析

src/runtime/map.go

	type hmap struct {
		count     int    // 表示当前哈希表中元素的数量
		flags     uint8  // 表示哈希表的标记 1表示buckets正在被使用 2表示oldbuckets正在被使用 4表示哈希正在被写入 8表示哈希是等量扩容
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


	// A bucket for a Go map.桶结构体
	type bmap struct {
		// 每个元素hash值的高8位，如果tophash[0] < minTopHash，表示这个桶的搬迁状态
		tophash [bucketCnt]uint8
		// 接下来是8个key、8个value，但是我们不能直接看到；为了优化对齐，go采用了key放在一起，value放在一起的存储方式，
		// 再接下来是hash冲突发生时，下一个溢出桶的地址
	}
		tophash的存在是为了快速试错，毕竟只有8位，比较起来会快一点。
	桶的结构到底是怎样的？
		桶的结构体并不是上面提到的tophash [8]uint8,因为go是不支持泛型的，所以在编译过程中才会根据具体的类型确定,实际上桶的结构可以表示为
		type bmap struct {
			topbits  [8]uint8
			keys     [8]keytype
			values   [8]valuetype
			pad      uintptr
			overflow uintptr
		}

从定义可以看出，不同于STL中map以红黑树实现的方式，Golang采用了HashTable的实现，解决冲突采用的是链地址法。也就是说，使用数组+链表来实现map

2. 初始化 makemap
	func makemap(t *maptype, hint int, h *hmap) *hmap {
		//计算内存空间和判断是否内存溢出
		mem, overflow := math.MulUintptr(uintptr(hint), t.bucket.size)
		if overflow || mem > maxAlloc {
			hint = 0
		}

		// initialize Hmap
		if h == nil {
			h = new(hmap)
		}
		h.hash0 = fastrand()

		//计算出指数B,那么桶的数量表示2^B
		B := uint8(0)
		for overLoadFactor(hint, B) {
			B++
		}
		h.B = B

		if h.B != 0 {
			var nextOverflow *bmap
			//根据B去创建对应的桶和溢出桶
			h.buckets, nextOverflow = makeBucketArray(t, h.B, nil)
			if nextOverflow != nil {
				h.extra = new(mapextra)
				h.extra.nextOverflow = nextOverflow
			}
		}

		return h
	}

	a. 计算出需要的内存空间并且判断内存是否溢出
	b.hmap没有的情况进行初始化，并设置hash0表示hash因子
	c.计算出指数B,桶的数量表示为2^B,通过makeBucketArray去创建对应的桶和溢出桶


3. 赋值mapassign
	func mapassign(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer {

		.....
		//计算出hash值
		hash :=t.hasher(key,uintptr(h.hash0))

		//更新状态为正在写入
		h.flags ^= hashWriting

	again:
		//通过hash获取对应的桶
		bucket := hash & bucketMask(h.B)
		b :=(*bmap)(unsafe.Pointer(uintptr(h.buckets)+bucket*uintptr(t.bucketsize)))
		//计算出tophash
		top :=tophash(hash)

		var inserti *uint8//记录插入的tophash
		var insertk unsafe.Pointer//记录插入的key值地址
		var elem unsafe.Pointer//记录插入的value值地址

	bucketloop:
		for{
			for i :=uintptr(0);i < bucketCnt;i++{
				//判断tophash是否相等
				if b.tophash[i] != top {
					//如果tophash不相等并且等于空,则可以插入该位置
					if isEmpty(b.tophash[i]) && inserti == nil {
						inserti = &b.tophash[i]
						//获取对应插入key和value的指针地址
						insertk = add(unsafe.Pointer(b), dataOffset+i*uintptr(t.keysize))
						elem = add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.keysize)+i*uintptr(t.elemsize))
					}
					if b.tophash[i] == emptyRest {
						break bucketloop
					}
					continue
				}

				//走到这里,说明已经存在,获得指定的key和value在桶得位置地址
				k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.keysize))
				//如果是指针，则要转化为指针
				if t.indirectkey() {
					k = *((*unsafe.Pointer)(k))
				}
				//判断key值是否相等
				if !t.key.equal(key, k) {
					continue
				}
				// already have a mapping for key. Update it.
				//如果key值需要修改，那么修改key值
				if t.needkeyupdate() {
					typedmemmove(t.key, k, key)
				}
				//获取value元素地址
				elem = add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.keysize)+i*uintptr(t.elemsize))
				goto done

				//未找到可插入的位置,找一下有没溢出桶，如果有继续执行写入操作
				ovf := b.overflow(t)
				if ovf == nil{
					break
				}
				b = ovf
			}
		}

		if inserti == nil {
			//如果在正常桶和溢出桶中都未找到插入的位置，那么得到一个新的溢出桶执行插入
			newb := h.newoverflow(t, b)
			inserti = &newb.tophash[0]
			insertk = add(unsafe.Pointer(newb), dataOffset)
			elem = add(insertk, bucketCnt*uintptr(t.keysize))
		}

		.....
		//将key值信息插入桶中指定位置
		typedmemmove(t.key, insertk, key)
		*inserti = top//更新tophash值
		h.count++

	done:
		h.flags &^= hashWriting
		if t.indirectelem() {
			elem = *((*unsafe.Pointer)(elem))
		}
		return elem //返回value的指针地址
	}

	a. 计算key的hash值,通过hash的高八位和低B为分别确定tophash和桶的序号
		tophash是什么?
			tophash是用来快速定位key和value的位置的,在查找或删除过程如果高8位hash都不相等，那么就没必要再去比较key值是否相等了，效率相对会高一些。
		如何定位到哪个桶执行插入?
			例如哈希表对应2^4个桶,即B是4,某个key的hash二进制值是如下值，那么如图可知该key对应的tophash值为10001100,即140，
			桶的值为0111,即是桶的序号为7
			hash := 100011001101111001110010010110000001111010110000100101011010111
	b. 每个桶可以存储8个tophash、8个key、8个value,遍历桶中的tophash,如果tophash不相等且是空的,说明该位置可以插入，
		分别获取对应位置key和value的地址并更新tophash。

4. 扩容 hashGrow
	//判断是否扩容的条件
		哈希表不是正在扩容的状态
		元素的数量 > 2^B次方(桶的数量) * 6.5,6.5表示为装载因子,很容易理解装载因子最大为8(一个桶能装载的元素数量)
		溢出桶过多,当前已经使用的溢出桶数量 >=2^B次方(桶的数量) ,B最大为15
	if !h.growing() && (overLoadFactor(h.count+1, h.B) || tooManyOverflowBuckets(h.noverflow, h.B)) {
		hashGrow(t, h)
		goto again // Growing the table invalidates everything, so try again
	}
	扩容分为两种，一种是等量扩容和2倍扩容：

	a. 扩容时，会将原来的buckets搬到oldbuckets
5. 读取mapaccess
	go的哈希查找有两种方式,一种是不返回ok的对应的源码方法为runtime.mapaccess1,另外返回ok的函数对应源码方法为runtime.mapaccess2
	a. mapaccess1_fat 返回一个值和 mapaccess2_fat返回两个值
	b. hash分为高位和地位，先通过低位快速找到bucket，再通过高位进一步查找，对比具体的key
	c. 访问到oldbuckets的数据时，会迁移到buckets

6. 删除mapdelete
	a. 引入emptyOne和emptyRest，后者为了加速查找

*/

// 排序及传参
func main() {
	// 初始化
	var m = map[int]int{}
	m[43] = 1
	var n = map[string]int{}
	n["abc"] = 1
	fmt.Println(m, n) // map[43:1] map[abc:1]

	// 排序：无序变有序
	mapScore := map[string]int{
		"Tom":   78,
		"Mary":  60,
		"Kevin": 90,
		"Danny": 100,
	}
	len := len(mapScore)
	names := make([]string, 0, len)
	for key, _ := range mapScore {
		names = append(names, key)
	}
	// 按字母升序
	sort.Strings(names)
	for _, name := range names {
		fmt.Println("Key:", name, "Value:", mapScore[name])
	}

	// 传递
	fmt.Println("传递前Tom的值", mapScore["Tom"])
	modify(mapScore)
	fmt.Println("传递后Tom的值", mapScore["Tom"])
	// 运行结果是可以看出key为"Tom"的值被修改了,说明map是引用类型
}

//修改
func modify(mapScore map[string]int) {
	mapScore["Tom"] = 66
}
