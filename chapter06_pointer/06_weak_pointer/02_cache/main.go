package main

import (
	"fmt"
	"runtime"
	"sync"
	"weak"
)

type User struct {
	Name string
}

var cache sync.Map // map[int]weak.Pointer[*User]

func GetUser(id int) *User {
	// ① 先从缓存里取
	if wp, ok := cache.Load(id); ok {
		if u := wp.(weak.Pointer[User]).Value(); u != nil {
			fmt.Println("cache hit")
			return u
		}
	}

	// ② 真正加载（这里直接构造）
	u := &User{Name: fmt.Sprintf("user-%d", id)}
	cache.Store(id, weak.Make(u))
	fmt.Println("load from DB")
	return u
}

func main() {
	u := GetUser(1) // load from DB
	fmt.Println(u.Name)

	runtime.GC() // 即使立刻 GC，因 main 持有强引用，User 仍在

	u = nil      // 释放最后一个强引用
	runtime.GC() // 触发 GC，User 可能被回收

	_ = GetUser(1) // 如被回收，会再次 load from DB
}
