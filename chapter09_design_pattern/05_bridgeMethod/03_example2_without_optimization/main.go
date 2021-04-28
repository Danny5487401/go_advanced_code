// 高并发版本

package main

import (
	"fmt"
	"sync"
)

type Set struct {
	impl sync.Map
}

func NewSet() *Set {
	return &Set{sync.Map{}}
}

func (s *Set) Add(key string) {
	s.impl.Store(key, true)
}

func (s *Set) Iter(f func(key string)) {
	s.impl.Range(func(key, value interface{}) bool {
		f(key.(string))
		return true
	})
}

func main() {
	s := NewSet()
	s.Add("hello")
	s.Add("world")
	s.Iter(func(key string) {
		fmt.Printf("key: %s\n", key)
	})
}

// 特点：我们更改了底层实现，但是 main 函数作为调用者，却不需要更改任何一行代码。这就是桥接模式的威力！