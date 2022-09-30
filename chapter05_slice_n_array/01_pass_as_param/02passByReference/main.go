package main

import "fmt"

func main() {
	arrayA := [2]int{100, 200}
	fmt.Printf("起始arrayA : %p , %v\n", &arrayA, arrayA) // 起始arrayA : 0xc00000a0a0 , [100 200]
	// 1.传数组指针
	testArrayPoint1(&arrayA)

	arrayB := arrayA[:]
	fmt.Printf("arrayB : %p , %v\n", &arrayB, arrayB) // arrayB : 0xc0000044c0 , [100 300]]

	// 2.传切片
	testArrayPoint2(&arrayB)

	fmt.Printf("最终arrayA : %p , %v\n", &arrayA, arrayA) // 最终arrayA : 0xc00000a0a0 , [100 400]
}

func testArrayPoint1(x *[2]int) {
	// 1.传数组指针，地址不变
	fmt.Printf("func Array1 : %p , %v\n", x, *x) // func Array1 : 0xc00000a0a0 , [100 200]

	// 增加100
	(*x)[1] += 100
}

func testArrayPoint2(x *[]int) {
	// 2.传指针切片
	fmt.Printf("func Array2 : %p , %v\n", x, *x) // func Array2 : 0xc0000044c0 , [100 300]

	// 增加100
	(*x)[1] += 100
}

/*
数组指针
	优点：就算是传入10亿的数组，也只需要再栈上分配一个8个字节的内存给指针就可以了

	缺点： 第一行和第三行指针地址都是同一个，万一原数组的指针指向更改了，那么函数里面的指针指向都会跟着更改

切片指针
	优点：用切片传数组参数，既可以达到节约内存的目的，也可以达到合理处理好共享内存的问题。
	打印结果第二行就是切片，切片的指针和原来数组的指针是不同的。
*/
