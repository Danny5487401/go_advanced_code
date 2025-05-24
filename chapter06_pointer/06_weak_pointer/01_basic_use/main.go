package main

import (
	"fmt"
	"runtime"
	"weak"
)

type MyStruct struct {
	Data string
}

func main() {
	obj := &MyStruct{Data: "example"}
	wp := weak.Make(obj) // 创建弱指针
	val := wp.Value()    // 调用 wp.Value() 时，如果对象仍存活则返回强引用，否则返回 nil。
	if val != nil {
		fmt.Println(val.Data)
	} else {
		fmt.Println("对象已被垃圾回收")
	}
	obj = nil // 移除强引用
	runtime.GC()
	if wp.Value() == nil {
		fmt.Println("对象已被垃圾回收")
	} else {
		fmt.Println("对象仍然存活")
	}
}
