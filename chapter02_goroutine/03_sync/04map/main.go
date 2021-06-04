package main

import (
	"fmt"
	"sync"
)

/*
背景：
	Go语言中的 map 在并发情况下，只读是线程安全的，同时读写是线程不安全的。
做法：
	Go语言在 1.9 版本中提供了一种效率较高的并发安全的 sync.Map，sync.Map 和 map 不同，不是以语言原生形态提供，而是在 sync 包下的特殊结构。
结构：
	type Map struct {
    mu Mutex

  	// 后面是readOnly结构体，依靠map实现，仅仅只用来读
    read atomic.Value // readOnly

    // 这个map主要用来写的，部分时候也承担读的能力
    dirty map[interface{}]*entry

    // 记录自从上次更新了read之后，从read读取key失败的次数
    misses int
}
使用：
	1. 使用Store(interface {}，interface {})添加元素。
	2. 使用Load(interface {}) interface {}检索元素。
	3. 使用Delete(interface {})删除元素。
	4. 使用LoadOrStore(interface {}，interface {}) (interface {}，bool)检索或添加之前不存在的元素。
		如果键之前在map中存在，则返回的布尔值为true。
	5. 使用Range遍历元素。
sync.Map特性：

	1.无须初始化，直接声明即可。
	2.sync.Map 不能使用 map 的方式进行取值和设置等操作，而是使用 sync.Map 的方法进行调用，Store 表示存储，Load 表示获取，Delete 表示删除。
	3.使用 Range 配合一个回调函数进行遍历操作，通过回调函数返回内部遍历出来的值，Range 参数中回调函数的返回值在需要继续迭代遍历时，返回 true，终止迭代遍历时，返回 false。

*/

func main() {

	m := &sync.Map{}

	// 添加元素
	m.Store(1, "one")
	m.Store(2, "two")

	// 获取元素1
	value, contains := m.Load(1)
	if contains {
		fmt.Printf("%s\n", value.(string))
	}

	// 返回已存value，否则把指定的键值存储到map中
	value, loaded := m.LoadOrStore(3, "three")
	if !loaded {
		fmt.Printf("%s\n", value.(string))
	}

	m.Delete(3)

	// 迭代所有元素
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
