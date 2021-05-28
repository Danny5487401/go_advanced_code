package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

/*
切片本身并不是动态数组或者数组指针。
它内部实现的数据结构通过指针引用底层数组，设定相关属性将数据读写操作限定在指定的区域内。切片本身是一个只读对象，其工作机制类似数组指针的一种封装。
切片（slice）是对数组一个连续片段的引用，所以切片是一个引用类型（因此更类似于 C++ 中的 Vector 类型，或者 Python 中的 list 类型）

结构
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
切片的结构体由3部分构成，Pointer 是指向一个数组的指针，len 代表当前切片的长度，cap 是当前切片的容量。cap 总是大于等于 len 的
 */

/*
1. 如果想从 slice 中得到一块内存地址:
s := make([]byte, 200)
ptr := unsafe.Pointer(&s[0])
 */

// 构造slice
/*Note
 1. Go 的内存地址中构造一个 slice:
var ptr unsafe.Pointer
var s1 = struct {
    addr uintptr
    len int
    cap int
}{uintptr(ptr), length, length}
s := *(*[]byte)(unsafe.Pointer(&s1))
 */

/*
2. Go 的反射中就存在一个与之对应的数据结构 SliceHeader,用它来构造一个 slice
var o []byte
sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&o))
sliceHeader.Cap = length
sliceHeader.Len = length
sliceHeader.Data = uintptr(ptr)
 */
func main()  {
	// 构造slice
	//方式一
	var ptr1 unsafe.Pointer
	var s1 = struct {
		addr uintptr
		len int
		cap int
	}{uintptr(ptr1), 3, 3}
	fmt.Printf("结构是%T,数值是%v\n",s1,s1)  // 结构是struct { addr uintptr; len int; cap int },数值是{0 3 3}[]
	s := *(*[]byte)(unsafe.Pointer(&s1))
	fmt.Printf("结构是%T,数值是%v\n",s,s) // 结构是[]uint8,数值是[]


	// 方式二
	var ptr2 unsafe.Pointer
	var o []byte
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&o))
	sliceHeader.Cap = 2
	sliceHeader.Len = 2
	sliceHeader.Data = uintptr(ptr2)
	fmt.Printf("结构是%T,数值是%v\n",sliceHeader,sliceHeader) // 结构是*reflect.SliceHeader,数值是&{0 2 2}

}
// 创建切片有两种形式，1. make 创建切片，空切片 2. 字面量也可以创建切片。
/*


func makeslice(et *_type, len, cap int) slice {
	// 根据切片的数据类型，获取切片的最大容量
	maxElements := maxSliceCap(et.size)
	// 比较切片的长度，长度值域应该在[0,maxElements]之间
	if len < 0 || uintptr(len) > maxElements {
		panic(errorString("makeslice: len out of range"))
	}
	// 比较切片的容量，容量值域应该在[len,maxElements]之间
	if cap < len || uintptr(cap) > maxElements {
		panic(errorString("makeslice: cap out of range"))
	}
	// 根据切片的容量申请内存
	p := mallocgc(et.size*uintptr(cap), et, true)
	// 返回申请好内存的切片的首地址
	return slice{p, len, cap}
}

func makeslice64(et *_type, len64, cap64 int64) slice {
	len := int(len64)
	if int64(len) != len64 {
		panic(errorString("makeslice: len out of range"))
	}

	cap := int(cap64)
	if int64(cap) != cap64 {
		panic(errorString("makeslice: cap out of range"))
	}

	return makeslice(et, len, cap)
}
*/


