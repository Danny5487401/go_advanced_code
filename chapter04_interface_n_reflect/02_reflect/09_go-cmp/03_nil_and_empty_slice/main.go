package main

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func main() {
	var s1 []int
	var s2 = make([]int, 0)

	var m1 map[int]int
	var m2 = make(map[int]int)

	fmt.Println("s1 equals s2?", cmp.Equal(s1, s2))
	fmt.Println("m1 equals m2?", cmp.Equal(m1, m2))

	fmt.Println("s1 equals s2 with option?", cmp.Equal(s1, s2, cmpopts.EquateEmpty()))
	fmt.Println("m1 equals m2 with option?", cmp.Equal(m1, m2, cmpopts.EquateEmpty()))

	// 如果，我们想要实现无序切片的比较（即只要两个切片包含相同的值就认为它们相等），可以使用cmpopts.SortedSlice选项先对切片进行排序，然后再进行比较

	s3 := []int{1, 2, 3, 4}
	s4 := []int{4, 3, 2, 1}
	fmt.Println("s3 equals s4?", cmp.Equal(s3, s4))
	fmt.Println("s3 equals s4 with option?", cmp.Equal(s3, s4, cmpopts.SortSlices(func(i, j int) bool { return i < j })))

	m3 := map[int]int{1: 10, 2: 20, 3: 30}
	m4 := map[int]int{1: 10, 2: 20, 3: 30}
	fmt.Println("m3 equals m4?", cmp.Equal(m1, m2))
	fmt.Println("m3 equals m4 with option?", cmp.Equal(m3, m4, cmpopts.SortMaps(func(i, j int) bool { return i < j })))

}
