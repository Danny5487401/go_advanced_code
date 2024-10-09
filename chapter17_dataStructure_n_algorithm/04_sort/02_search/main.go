package main

import (
	"fmt"
	"sort"
)

func main() {
	// 使用 binary search

	// SearchInts 查找数据
	arr1 := []int{2, 1, 6, 5, 3}
	// 错误使用: 直接使用，应该先排序
	sort.Ints(arr1)
	fmt.Println(arr1)
	foundPosition := sort.SearchInts(arr1, 6)
	unfoundPosition := sort.SearchInts(arr1, 4)
	fmt.Println("找到元素索引位置是", foundPosition)
	fmt.Println("找不到元素索引位置是", unfoundPosition)

	// sort.Search 用于从一个已经排序的数组中找到某个值所对应的索引
	arrInts := []int{13, 35, 56, 79}
	target := 35

	// 等价 sort.SearchInts
	// findPos := sort.SearchInts(arrInts,35)
	findPos := sort.Search(len(arrInts), func(i int) bool {
		// 注意不是 arrInts[i] == 35
		return arrInts[i] >= target
	})

	if findPos < len(arrInts) && arrInts[findPos] == target {
		fmt.Printf("target is present at data[%v] \n", findPos)
	} else {
		fmt.Printf("target is not present at data[%v] \n", findPos)
		// but i is the index where it would be inserted.
	}
}
