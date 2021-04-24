package main

// Get 返回 Pool 中的任意一个对象。如果 Pool 为空，则调用 New 返回一个新创建的对象。

import (
	"fmt"
	"sync"
)

func main()  {
	// 建立对象
	var pipe = &sync.Pool{New:func()interface{}{return "Hello, danny"}}

	// 准备放入的字符串
	val := "Hello,World!"

	// 放入
	pipe.Put(val)

	// 取出
	first := pipe.Get().(string)
	fmt.Println(first)

	// 再取就没有了,会自动调用NEW
	second := pipe.Get().(string)
	fmt.Println(second)
}

