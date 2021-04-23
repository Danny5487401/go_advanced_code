package main

import (
	"fmt"
	"sync"
)

/*
1. 使用Store(interface {}，interface {})添加元素。
2. 使用Load(interface {}) interface {}检索元素。
3. 使用Delete(interface {})删除元素。
4. 使用LoadOrStore(interface {}，interface {}) (interface {}，bool)检索或添加之前不存在的元素。
	如果键之前在map中存在，则返回的布尔值为true。
5. 使用Range遍历元素。

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
