package main

import "fmt"

func main() {
	// 数组转切片
	arrayToSlice()
	// 切片转数组, 需求来自这个 issue：https://github.com/golang/go/issues/395
	//sliceToArrayBeforeGo_1_20()
	sliceToArrayInGo_1_20()
}

// 1.20 之前方式
func sliceToArrayBeforeGo_1_20() {
	// 获取切片的底层数组地址
	var sl = []int{1, 2, 3, 4, 5, 6, 7}
	// parr就是指向切片sl底层数组的指针
	var parr = (*[7]int)(sl)
	// arr则是底层数组的一个副本
	var arr = *(*[7]int)(sl)
	fmt.Println(sl)  // [1 2 3 4 5 6 7]
	fmt.Println(arr) // [1 2 3 4 5 6 7]

	sl[0] = 11
	fmt.Println(sl)    // [11 2 3 4 5 6 7]
	fmt.Println(arr)   // [1 2 3 4 5 6 7]
	fmt.Println(*parr) // [11 2 3 4 5 6 7]
}

// 1.20 方式
func sliceToArrayInGo_1_20() {
	var sl = []int{1, 2, 3, 4, 5, 6, 7}
	var parr = (*[7]int)(sl) // Note:转换后的数组长度要小于等于切片长度
	var arr = [7]int(sl)     // Note:转换后的数组长度要小于等于切片长度
	fmt.Println(sl)          // [1 2 3 4 5 6 7]
	fmt.Println(arr)         // [1 2 3 4 5 6 7]

	sl[0] = 11
	fmt.Println(arr)  // [1 2 3 4 5 6 7]
	fmt.Println(parr) // &[11 2 3 4 5 6 7]
}

func arrayToSlice() {
	// slice 表达式（slice expressions）可以从一个数组得到一个切片
	// a[low : high : max] max 可以省略
	a := [5]int{1, 2, 3, 4, 5}
	s := a[1:4]
	fmt.Println(s) // [2 3 4]
}
