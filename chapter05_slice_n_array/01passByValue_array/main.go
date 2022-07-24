package main

import "fmt"

// 在 Go 中，与 C 数组变量隐式作为指针使用不同，Go 数组是值类型，赋值和函数传参操作都会复制整个数组数据

func main() {
	arrayA := [2]int{100, 200}
	var arrayB [2]int
	// 地址变化
	arrayB = arrayA
	fmt.Printf("arrayA : %p , %v\n", &arrayA, arrayA) // arrayA : 0xc00000a0a0 , [100 200]
	fmt.Printf("arrayB : %p , %v\n", &arrayB, arrayB) // arrayB : 0xc00000a0b0 , [100 200]

	testArray(arrayA)
}

func testArray(x [2]int) {
	// 值传递，变地址
	fmt.Printf("func Array : %p , %v\n", &x, x) // func Array : 0xc00000a100 , [100 200]
}

// 结论 ：三个内存地址都不同，这也就验证了 Go 中数组赋值和函数传参都是值复制的
// 缺点： 假想每次传参都用数组，那么每次数组都要被复制一遍。
//		如果数组大小有 100万，在64位机器上就需要花费大约 800W 字节，即 8MB 内存。
//		这样会消耗掉大量的内存。于是乎有人想到，函数传参用数组的指针
