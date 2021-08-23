package main

import (
	"fmt"
	"unsafe"
)

/*
一。Golang指针分为3种
	1.  *类型:普通指针类型，用于传递对象地址，不能进行指针运算。
	2.  unsafe.Pointer:通用指针类型，用于转换不同类型的指针，不能进行指针运算，不能读取内存存储的值（必须转换到某一类型的普通指针）。
	3.  uintptr:用于指针运算，GC 不把 uintptr 当指针，uintptr 无法持有对象。uintptr 类型的目标会被回收。
		注意：uintptr是平台相关的，在32位系统下大小是4bytes，在64位系统下是8bytes,所以不可移植.
		uintptr 并没有指针的语义，意思就是 uintptr 所指向的对象会被 gc 无情地回收

二。 为什么有 unsafe
	Go 语言类型系统是为了安全和效率设计的，有时，安全会导致效率低下。有了 unsafe 包，高阶的程序员就可以利用它绕过类型系统的低效。
	因此，它就有了存在的意义，阅读 Go 源码，会发现有大量使用 unsafe 包的例子。
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
源码分析
unsafe 包 两个类型，三个函数

	type ArbitraryType int
	type Pointer *ArbitraryType

	func Sizeof(x ArbitraryType) uintptr
		unsafe.Sizeof接受任意类型的值(表达式)，返回其占用的字节数,这和c语言里面不同，
		Note:如果是slice，则不会返回这个slice在内存中的实际占用长度，一个 slice 的大小则为 slice header 的大小
		c语言里面sizeof函数的参数是类型，而这里是一个表达式，比如一个变量。
		返回类型 x 所占据的字节数，但不包含 x 所指向的内容的大小。例如，对于一个指针，函数返回的大小为 8 字节（64位机上），


	func Offsetof(x ArbitraryType) uintptr
		unsafe.Offsetof： 返回结构体成员在内存中的位置离结构体起始处的字节数，所传参数必须是结构体的成员

	func Alignof(x ArbitraryType) uintptr
		Alignof返回变量对齐字节数量m，Offsetof返回变量指定属性的偏移量，它分配到的内存地址能整除 m.
		这个函数虽然接收的是任何类型的变量，但是有一个前提，就是变量要是一个struct类型，且还不能直接将这个struct类型的变量当作参数，
		只能将这个struct类型变量的属性当作参数


	三个函数的参数均是ArbitraryType类型，就是接受任何类型的变量,返回的结果都是 uintptr 类型，这和 unsafe.Pointer 可以相互转换。
	三个函数都是在编译期间执行，它们的结果可以直接赋给 const型变量。另外，因为三个函数执行的结果和操作系统、编译器相关，所以是不可移植的
解释：
	ArbitraryType是int的一个别名，在Go中对ArbitraryType赋予特殊的意义。代表一个任意Go表达式类型。实际上它类似于 C 语言里的 void*。
	Pointer是int指针类型的一个别名，在Go中可以把Pointer类型，理解成任何指针的父类型。

注意：
	高阶的 Gopher，怎么能不会使用 unsafe 包呢？它可以绕过 Go 语言的类型系统，直接操作内存。例如，一般我们不能操作一个结构体的未导出成员，
	但是通过 unsafe 包就能做到。unsafe 包让我可以直接读写内存，还管你什么导出还是未导
*/

type Programmer struct {
	Name     string //名字
	Language string //爱好
}

func main() {
	// 一。结构体操作
	p := Programmer{Name: "danny", Language: "Golang"}
	fmt.Println("修改前：", p)
	//获取 name的指针
	name := (*string)(unsafe.Pointer(&p))
	*name = "Joy"
	// offset使用获取language地址
	lang := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + unsafe.Offsetof(p.Name)))
	*lang = "Python"
	fmt.Println("修改后：", p)
	//异常情况 示例
	//... 中间逻辑使personAaddr2指向不合法位置
	//personB = (*Person)(unsafe.Pointer(uintptr(0)))
	//fmt.Println("personB.Age is :", personB.Age)

	// 二。切片操作
	/* runtime/slice.go
	type slice struct{
		array unsafe.Pointer
		len int
		cap int

	}
	func makeslice() slice  返回slice 结构体
	*/
	s := make([]int, 9, 20)
	var len1 = *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + uintptr(8)))
	fmt.Println("长度", len1, len(s))
	var cap1 = *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + uintptr(16)))
	fmt.Println("容量", cap1, cap(s))
	// 转换过程 Len: &s => pointer => uintptr => pointer => *int => int

	// 三。获取map的长度
	/*
		type hmap struct{
			count int
			flag uint8
			B	uint8
			....
		}
		和 slice 不同的是，makemap 函数返回的是 hmap 的指针
		func makemap()*map
		我们依然能通过 unsafe.Pointer 和 uintptr 进行转换，得到 hamp 字段的值，只不过，现在 count 变成二级指针
	*/
	mp := make(map[string]int)
	mp["danny"] = 1
	mp["Joy"] = 2
	count := **(**int)(unsafe.Pointer(&mp))
	// 转换过程&mp->pointer->**int->int
	fmt.Println("长度", count, len(mp))

}
