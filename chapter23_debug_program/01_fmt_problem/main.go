package main

import (
	"fmt"
)

func main() {

	generalVerb()

	intVerb()

	floatVerb()

	wideVerb()

	fmtProblem()
}

// 一般动词
func generalVerb() {
	p := Person{"Alice", 30}

	fmt.Printf("%%v 默认格式：%v\n", p)
	fmt.Printf("%%+v 包含字段名：%+v\n", p)
	fmt.Printf("%%#v Go 语法表示：%#v\n", p)
	fmt.Printf("%%T 类型表示：%T\n", p)
	fmt.Printf("%%%% 百分号：%%\n")

	// %v 和 %t 对布尔值的输出是相同的，但使用 %t 可以明确地表明你正在格式化布尔值。
	fmt.Printf("%%t 布尔值：%t\n", true)
}

// 整数类型
func intVerb() {
	num := 255
	fmt.Printf("%%b 二进制：%b\n", num)
	fmt.Printf("%%c 字符：%c\n", num)
	fmt.Printf("%%d 十进制：%d\n", num)
	fmt.Printf("%%o 八进制：%o\n", num)
	fmt.Printf("%%O 带前缀的八进制：%O\n", num)
	fmt.Printf("%%q 单引号字符：%q\n", num)
	fmt.Printf("%%x 小写十六进制：%x\n", num)
	fmt.Printf("%%X 大写十六进制：%X\n", num)
	fmt.Printf("%%U Unicode 格式：%U\n", num)
}

// 浮点数和复数
func floatVerb() {
	num := 12345.6789
	complexNum := 1.234 + 5.678i

	fmt.Printf("%%b 无小数点科学计数法：%b\n", num)
	fmt.Printf("%%e 科学计数法（小写 e）：%e\n", num)
	fmt.Printf("%%E 科学计数法（大写 E）：%E\n", num)
	fmt.Printf("%%f 十进制形式：%f\n", num)
	fmt.Printf("%%F 与 %%f 相同：%F\n", num)
	fmt.Printf("%%g 自动选择：%g\n", num)
	fmt.Printf("%%G 自动选择（大写）：%G\n", num)
	fmt.Printf("%%x 十六进制表示：%x\n", num)
	fmt.Printf("%%X 十六进制表示（大写）：%X\n", num)
	fmt.Println("----------")
	fmt.Printf("%%b 无小数点科学计数法（复数）：%b\n", complexNum)
	fmt.Printf("%%e 科学计数法（小写 e）（复数）：%e\n", complexNum)
	fmt.Printf("%%E 科学计数法（大写 E）（复数）：%E\n", complexNum)
	fmt.Printf("%%f 十进制形式（复数）：%f\n", complexNum)
	fmt.Printf("%%F 与 %%f 相同（复数）：%F\n", complexNum)
	fmt.Printf("%%g 自动选择（复数）：%g\n", complexNum)
	fmt.Printf("%%G 自动选择（大写）（复数）：%G\n", complexNum)
	fmt.Printf("%%x 十六进制表示（复数）：%x\n", complexNum)
	fmt.Printf("%%X 十六进制表示（大写）（复数）：%X\n", complexNum)
}

func stringVerb() {
	str := "hello"
	str1 := `"hello"`
	data := []byte{72, 101, 108, 108, 111}

	fmt.Printf("%%s 原始字节：%s\n", str)
	fmt.Printf("%%q 双引号字符串：%q\n", str)
	fmt.Printf("%%q 双引号字符串：%q\n", str1)
	fmt.Printf("%%x 小写十六进制：%x\n", data)
	fmt.Printf("%%X 大写十六进制：%X\n", data)
}

func sliceVerb() {
	slice := []int{10, 20, 30}

	fmt.Printf("%%p 切片首元素地址：%p\n", &slice[0])
	fmt.Printf("%%b 二进制表示（切片）：%b\n", slice)
	fmt.Printf("%%d 十进制表示（切片）：%d\n", slice)
	fmt.Printf("%%o 八进制表示（切片）：%o\n", slice)
	fmt.Printf("%%x 小写十六进制表示（切片）：%x\n", slice)
	fmt.Printf("%%X 大写十六进制表示（切片）：%X\n", slice)
}

func precisionVerb() {
	num := 123.456789
	str := "GoLang Programming"

	fmt.Printf("默认精度的浮点数：%f\n", num)    // 输出 123.456789 默认精度 6
	fmt.Printf("精度 2 的浮点数：%.2f\n", num) // 输出 123.46 精度为 2
	fmt.Printf("精度 5 的字符串：%.5s\n", str) // 输出 'GoLan' 精度为 5，截断字符串
}

// 问题
func fmtProblem() {
	// 1. 结构体中含有指针对象
	ins := Instance{
		A: "AAAA",
		B: 1000,
		C: &Inner{
			D: "DDDD",
			E: "EEEE",
		},
	}
	fmt.Println(ins) // {AAAA 1000 0x14000070000}
	// 由于 C 字段是指针，所以打印出来的是一个地址0xc000054020，而地址背后的数据却被隐藏了。显然，这对程序排查非常不友好

	// 2. 数组或者map中是指针对象时
	arr := [...]*Demo{{100, "Python"}, {200, "Golang"}}
	fmt.Printf("%v\n-----------------分割线-----------\n", arr) // [0x140000ac000 0x140000ac018]

	// 3. 循环结构
	c := &Circular{1, nil}
	c.next = &Circular{2, c}

	fmt.Printf("%+v\n----------------分割线-------------------\n", c) // &{a:1 next:0x140000a6030}
}

// 宽度打印
func wideVerb() {
	num := 123.45
	str := "GoLang"

	fmt.Printf("宽度 10 的浮点数：%10f\n", num) // 输出 123.450000 右对齐，宽度为 10
	fmt.Printf("宽度 10 的字符串：%10s\n", str) // 输出 '      GoLang' 右对齐，宽度为 10
}

type Circular struct {
	a    int
	next *Circular
}

type Instance struct {
	A string
	B int
	C *Inner
}

type Inner struct {
	D string
	E string
}

type Demo struct {
	a int
	b string
}

type Person struct {
	Name string
	Age  int
}
