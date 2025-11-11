package main

import (
	"fmt"
	"strconv"
)

/*
【扩展阅读 a的典故】这是C语言遗留下的典故。C语言中没有string类型而是用字符数组(array)表示字符串，所以Itoa对很多C系的程序员很好理解。
*/
func stringIntTransfer() {
	// 1. Atoi():将字符串类型的整数转换为int类型
	s1 := "100"
	i1, err := strconv.Atoi(s1)
	if err != nil {
		fmt.Println("can't convert to int")
	} else {
		fmt.Printf("type:%T value:%#v\n", i1, i1) //type:int value:100
	}

	// 2. Itoa()函数用于将int类型数据转换为对应的字符串表示
	i2 := 200
	s2 := strconv.Itoa(i2)
	fmt.Printf("type:%T value:%#v\n", s2, s2) //type:string value:"200"
}

func parseFunction() {
	// Parse类函数用于转换字符串为给定类型的值：ParseBool()、ParseFloat()、ParseInt()、ParseUint()。
	b, _ := strconv.ParseBool("true")
	f, _ := strconv.ParseFloat("3.1415", 64)
	i, _ := strconv.ParseInt("-2", 10, 64)
	u, _ := strconv.ParseUint("2", 10, 64)
	fmt.Println(b, f, i, u)
}

func formatFunction() {
	// Format系列函数实现了将给定类型数据格式化为string类型数据的功能。

	s1 := strconv.FormatBool(true)

	// bitSize表示f的来源类型（32：float32、64：float64），会据此进行舍入。
	//fmt表示格式：’f’（-ddd.dddd）、’b’（-ddddp±ddd，指数为二进制）、’e’（-d.dddde±dd，十进制指数）、’E’（-d.ddddE±dd，十进制指数）、’g’（指数很大时用’e’格式，否则’f’格式）、’G’（指数很大时用’E’格式，否则’f’格式）。
	//prec控制精度（排除指数部分）：对’f’、’e’、’E’，它表示小数点后的数字个数；对’g’、’G’，它控制总的数字个数。如果prec 为-1，则代表使用最少数量的、但又必需的数字来表示f。
	s2 := strconv.FormatFloat(3.1415, 'E', -1, 64)
	s3 := strconv.FormatInt(-2, 16)
	s4 := strconv.FormatUint(2, 16)
	fmt.Println(s1, s2, s3, s4)
}

func main() {
	stringIntTransfer()

	parseFunction()

	formatFunction()
}
