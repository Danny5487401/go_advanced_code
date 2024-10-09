package main

import (
	"container/ring"
	"fmt"
)

func main() {
	// 创建一个环, 包含 3 个元素
	r := ring.New(3)
	fmt.Printf("ring: %+v\n", *r)

	// 初始化
	for i := 1; i <= 3; i++ {
		r.Value = i
		r = r.Next()
	}
	fmt.Printf("init ring: %+v\n", *r)

	// sum
	s := 0
	r.Do(func(i interface{}) {
		fmt.Println(i)
		s += i.(int)
	})
	fmt.Printf("sum ring: %d\n", s)
}
