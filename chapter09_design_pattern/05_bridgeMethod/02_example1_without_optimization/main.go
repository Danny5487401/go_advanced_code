package main

import (
	"fmt"
)

type Set struct {
	impl map[string]bool
}

func NewSet() *Set {
	return &Set{make(map[string]bool)}
}

// 添加的功能
func (s *Set) Add(key string) {
	s.impl[key] = true
}

func (s *Set) Iter(f func(key string)) {
	for key := range s.impl {
		f(key)
	}
}

func main() {
	s := NewSet()
	s.Add("hello")
	s.Add("world")
	s.Iter(func(key string) {
		fmt.Printf("key: %s\n", key)
	})
}
// Set这个结构体里，封装了其他底层实现，Set规定了它所提供的方法，但是底层具体的实现，确是由 map_test[string_test]bool 真正 提供的，
//	但是Set的使用者并不知道这个事实，因此对调用者而言，Set实现提供了功能，但是没有暴露底层实现。
