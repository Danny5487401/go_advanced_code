package main

import (
	"fmt"
	"sync"
)



func main() {

	m := &sync.Map{}

	// 添加元素
	m.Store(1, "one")
	m.Store(2, "two")

	// 再次添加相同元素
	m.Store(1, "onePlus")

	// 获取元素1
	value, contains := m.Load(1)
	if contains {
		fmt.Printf("获取的结果%s\n", value.(string))
	}

	// 返回已存value，否则把指定的键值存储到map中
	value, loaded := m.LoadOrStore(3, "three")
	if !loaded {
		fmt.Printf("第一次写入%s\n", value.(string))
	}
	// 再次读取或写入 "3"，因为这个 key 已经存在，因此写入不成功，并且读出原值。
	value, loaded = m.LoadOrStore(3, "threePlus")
	if loaded {
		//已经加载成功
		fmt.Printf("已存在:%s\n", value.(string))
	}
	// 获取元素3
	value, contains = m.Load(3)
	if contains {
		fmt.Printf("最终存的数值%s\n", value.(string))
	}

	// 删除元素
	m.Delete(3)


	// 迭代所有元素
	fmt.Println("-----开始遍历------")
	m.Range(func(key, value interface{}) bool {
		fmt.Printf("%d: %s\n", key.(int), value.(string))
		return true
	})
}

/*
注意：
	sync.Map 没有提供获取 map 数量的方法，替代方法是在获取 sync.Map 时遍历自行计算数量，sync.Map 为了保证并发安全有一些性能损失，
	因此在非并发情况下，使用 map 相比使用 sync.Map 会有更好的性能
*/
