package main

import (
	"fmt"
	"unsafe"
)

/* Golang指针分为3种
1.  *类型:普通指针类型，用于传递对象地址，不能进行指针运算。
2.  unsafe.Pointer:通用指针类型，用于转换不同类型的指针，不能进行指针运算，不能读取内存存储的值（必须转换到某一类型的普通指针）。
3.  uintptr:用于指针运算，GC 不把 uintptr 当指针，uintptr 无法持有对象。uintptr 类型的目标会被回收。
	注意：uintptr是平台相关的，在32位系统下大小是4bytes，在64位系统下是8bytes。

unsafe.Pointer 是桥梁，可以让任意类型的指针实现相互转换，也可以将任意类型的指针转换为 uintptr 进行指针运算。
	unsafe.Pointer 可以让你的变量在不同的普通指针类型转来转去，也就是表示为任意可寻址的指针类型。
	而 uintptr 常用于与 unsafe.Pointer 打配合，用于做指针运算

1. unsafe.Pointer   通用指针

	（1）任何类型的指针都可以被转化为Pointer
	（2）Pointer可以被转化为任何类型的指针
	（3）uintptr可以被转化为Pointer
	（4）Pointer可以被转化为uintptr
	Note : 我们不可以直接通过*p来获取unsafe.Pointer指针指向的真实变量的值，因为我们并不知道变量的具体类型。
	和普通指针一样，unsafe.Pointer指针也是可以比较的，并且支持和nil常量比较判断是否为空指针


2. uintptr   整数类型
	定义: uintptr is an integer type that is large enough to hold the bit pattern of any 03PointerSetPrivateValue
	源码：type uintptr uintptr
*/

/* unsafe 包 两个类型，三个函数

type ArbitraryType int
type Pointer *ArbitraryType

func Sizeof(x ArbitraryType) uintptr
	unsafe.Sizeof接受任意类型的值(表达式)，返回其占用的字节数,这和c语言里面不同，
	Note:如果是slice，则不会返回这个slice在内存中的实际占用长度。
	c语言里面sizeof函数的参数是类型，而这里是一个表达式，比如一个变量

func Offsetof(x ArbitraryType) uintptr
	unsafe.Offsetof：返回结构体中元素所在内存的偏移量。

func Alignof(x ArbitraryType) uintptr
	Alignof返回变量对齐字节数量，Offsetof返回变量指定属性的偏移量，这个函数虽然接收的是任何类型的变量，
	但是有一个前提，就是变量要是一个struct类型，且还不能直接将这个struct类型的变量当作参数，只能将这个struct类型变量的属性当作参数


三个函数的参数均是ArbitraryType类型，就是接受任何类型的变量。
解释：
	ArbitraryType是int的一个别名，在Go中对ArbitraryType赋予特殊的意义。代表一个任意Go表达式类型。
	Pointer是int指针类型的一个别名，在Go中可以把Pointer类型，理解成任何指针的父类型。
*/

type Person struct {
	Age   int32  //年龄
	Name  string //名字
	Hobby string //爱好
}

func main() {
	var personA Person
	var a byte
	//Alignof 示例
	align := unsafe.Alignof(a)
	fmt.Println("align 变量对齐字节数量: ", align)
	//Sizeof 示例
	size := unsafe.Sizeof(a)
	fmt.Println("size变量在内存中占用的字节数: ", size)
	//Offsetof 示例
	offset := unsafe.Offsetof(personA.Name)
	fmt.Println("offset变量指定Name属性的偏移量: ", offset)

	//Pointer 示例

	personAddr := (uintptr)(unsafe.Pointer(&(personA)))
	fmt.Println("personAddr is: ", personAddr)

	nameAddr := (uintptr)(unsafe.Pointer(&(personA.Name)))
	fmt.Println("nameAddr is: ", nameAddr)

	personA.Age = 100
	personA.Hobby = "run"

	//指针操作 示例
	personAddr2 := nameAddr - offset
	fmt.Println("personAddr2 is: ", personAddr2)
	// 指向同一个区域
	personB := (*Person)(unsafe.Pointer(personAddr2))
	fmt.Println("personB.Age is :", personB.Age)

	//异常情况 示例
	//... 中间逻辑使personAaddr2指向不合法位置
	personB = (*Person)(unsafe.Pointer(uintptr(0)))
	fmt.Println("personB.Age is :", personB.Age)
}
