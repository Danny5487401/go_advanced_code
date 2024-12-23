package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// 指针类型转换
/*
	[]byte和string其实内部的存储结构都是一样的，但 Go 语言的类型系统禁止他俩互换。如果借助unsafe.Pointer，我们就可以实现在零拷贝的情况下，
	将[]byte数组直接转换成string类型,实现字符串和 bytes 切片之间的转换，要求是 zero-copy

底层数据结构
type StringHeader struct{
	Data uintptr
	Len int
}

type SliceHeader struct{
	Data uintptr
	len int
	cap int
}

*/

func main() {
	/*
			func bytesToString(bs []byte) string {
			    return string(bs)
			}
		底层是调用了runtime.slicebytetostring()函数
	*/

	// 转换方式一:unsafe写法
	bytes := []byte{104, 101, 108, 108, 111}
	g := unsafe.Pointer(&bytes) //强制转换成unsafe.Pointer，编译器不会报错
	str := *(*string)(g)        //然后强制转换成string类型的指针，再将这个指针的值当做string类型取出来
	fmt.Println(str)            //输出 "hello"

	// 转换方式二: 反射方式
	fmt.Println("字符串转bytes:", String2bytes("Danny"))
	fmt.Println("bytes转字符串:", Bytes2string([]byte{104, 101, 108, 108, 111}))

	// 转换方式三: 1.20
	fmt.Println("v1_20 字符串转bytes:", String2ByteSlice("Danny"))
	fmt.Println("v1_20 bytes转字符串: ", ByteSlice2String([]byte{104, 101, 108, 108, 111}))
}

func String2bytes(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2string(b []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}
	return *(*string)(unsafe.Pointer(&sh))

}

func String2ByteSlice(str string) []byte {
	if str == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func ByteSlice2String(bs []byte) string {
	if len(bs) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}
