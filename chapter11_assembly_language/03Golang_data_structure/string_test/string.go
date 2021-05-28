package main

import (
	"fmt"
	"unicode/utf8"
	"unsafe"
)

//  字符串底层结构定义在源码runtime包下的 string.go 文件中

// src/runtime/string.go
// 由于runtime.stringStruct结构是非导出的，我们不能直接使用。
// 所以我在代码中手动定义了一个stringStruct结构体，字段与runtime.stringStruct完全相同
type stringStruct struct {
	str unsafe.Pointer
	len int
}
// str：一个指针，指向存储实际字符串的内存地址。
// len：字符串的长度。与切片类似，在代码中我们可以使用len()函数获取这个值。注意，len存储实际的字节数，而非字符数。
//	所以对于非单字节编码的字符，结果可能让人疑惑。后面会详细介绍多字节字符



//go:noinline
func stringParam(s string) {}


func main() {
	// 制造string
	s := "Hello World!"
	fmt.Println(*(*stringStruct)(unsafe.Pointer(&s))) // {0x102ad4f91 12}

	for _, b := range s {
		fmt.Printf("%v\t",b)  // 72      101     108     108     111     32      87      111     114     108     100     33
	}
	fmt.Println("")
	// 对于使用非 ASCII 字符的字符串，我们可以使用标准库的 unicode/utf8 包中的RuneCountInString()方法获取实际字符数
	s1 := "Hello World!"
	s2 := "你好，中国"

	fmt.Println(utf8.RuneCountInString(s1)) // 12
	fmt.Println(utf8.RuneCountInString(s2)) // 5
}
