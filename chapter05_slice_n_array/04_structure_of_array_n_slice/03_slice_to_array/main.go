package main

import "fmt"

func main() {
	//beforeGo_1_20()
	go_1_20()
}

// 1.20 之前方式
func beforeGo_1_20() {
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

// Note:转换后的数组长度要小于等于切片长度
func go_1_20() {
	var sl = []int{1, 2, 3, 4, 5, 6, 7}
	var arr = [7]int(sl)
	var parr = (*[7]int)(sl)
	fmt.Println(sl)  // [1 2 3 4 5 6 7]
	fmt.Println(arr) // [1 2 3 4 5 6 7]
	sl[0] = 11
	fmt.Println(arr)  // [1 2 3 4 5 6 7]
	fmt.Println(parr) // &[11 2 3 4 5 6 7]
}
