package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

/* unsafe 包 两个类型，三个函数

type ArbitraryType int
type Pointer *ArbitraryType

//  unsafe.Sizeof接受任意类型的值(表达式)，返回其占用的字节数,这和c语言里面不同，
	c语言里面sizeof函数的参数是类型，而这里是一个表达式，比如一个变量
func Sizeof(x ArbitraryType) uintptr

// 	unsafe.Offsetof：返回结构体中元素所在内存的偏移量。
func Offsetof(x ArbitraryType) uintptr

//  Alignof返回变量对齐字节数量Offsetof返回变量指定属性的偏移量，这个函数虽然接收的是任何类型的变量，
	但是有一个前提，就是变量要是一个struct类型，且还不能直接将这个struct类型的变量当作参数，只能将这个struct类型变量的属性当作参数
func Alignof(x ArbitraryType) uintptr

三个函数的参数均是ArbitraryType类型，就是接受任何类型的变量。
解释：
	ArbitraryType是int的一个别名，在Go中对ArbitraryType赋予特殊的意义。代表一个任意Go表达式类型。
	Pointer是int指针类型的一个别名，在Go中可以把Pointer类型，理解成任何指针的父类型。
 */


// 普通指针类型转换
func main()  {
	v1 := uint(12)
	v2 := int(13)

	fmt.Println(reflect.TypeOf(v1)) //uint
	fmt.Println(reflect.TypeOf(v2)) //int

	fmt.Println(reflect.TypeOf(&v1)) //*uint
	fmt.Println(reflect.TypeOf(&v2)) //*int


	p := (*uint)(unsafe.Pointer(&v2)) //使用unsafe.Pointer进行类型的转换

	fmt.Println(reflect.TypeOf(p)) // *unit
	fmt.Println(*p) //13
}
