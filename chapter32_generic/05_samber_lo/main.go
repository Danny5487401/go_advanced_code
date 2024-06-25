package main

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

func main() {

	// 1 Filter
	fmt.Println(lo.Filter[string]([]string{"hello", "good bye", "world", "fuck", "fuck who"}, func(s string, _ int) bool {
		return !strings.Contains(s, "fuck")
	}))

	// 2 GroupBy
	fmt.Println(lo.GroupBy[int, int]([]int{1, 2, 3, 4, 5, 6, 7, 8}, func(i int) int {
		return i % 3
	}))

	// 3 Map的辅助函数
	// 创建一个包含map的所有key的切片
	fmt.Println(lo.Keys[string, int](map[string]int{"foo": 1, "bar": 2}))
	// lo.Values[string, int](map[string]int{"foo": 1, "bar": 2})
	fmt.Println(lo.Values[string, int](map[string]int{"foo": 1, "bar": 2}))

	// 返回由给定谓词函数过滤的相同类型的map。
	fmt.Println(lo.PickBy(map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(key string, value int) bool {
		return value%2 == 1
	}))

	// 将map转换为键值对数组
	fmt.Println(lo.Entries(map[string]int{"foo": 1, "bar": 2}))
	// 将键值对数组转换为map
	fmt.Println(lo.FromEntries([]lo.Entry[string, int]{
		{
			Key:   "foo",
			Value: 1,
		},
		{
			Key:   "bar",
			Value: 2,
		},
	}))

	// 将多个map从左到右合并
	fmt.Println(lo.Assign[string, int](
		map[string]int{"a": 1, "b": 2},
		map[string]int{"b": 3, "c": 4},
	))

	// 将map转换为切片
	fmt.Println(lo.MapToSlice(map[int]int64{1: 4, 2: 5, 3: 6}, func(k int, v int64) string {
		return fmt.Sprintf("%d_%d", k, v)
	}))

}
