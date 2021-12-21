package main

import "fmt"

// 闭包基本使用
func adder() func(int) int {
	sum := 0
	innerFunc := func(x int) int {
		sum += x
		return sum
	}
	return innerFunc

}

func main() {
	pos, neg := adder(), adder()
	for i := 0; i < 10; i++ {
		fmt.Println(pos(i), neg(-2*i))
	}

}

/*	结果
0 0
1 -2
3 -6
6 -12
10 -20
15 -30
21 -42
28 -56
36 -72
45 -90

*/

/*
解析
	当用不同的参数调用adder函数得到（pos(i)，neg(i)）函数时，得到的结果是隔离的，
	也就是说每次调用adder返回的函数都将生成并保存一个新的局部变量sum。其实这里adder函数返回的就是闭包
*/
