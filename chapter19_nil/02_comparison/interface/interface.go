package main

import (
	"fmt"
	"unsafe"
)

//数据结构
// 空接口
//1.不含方法
type eface struct {
	_type *int           // int类型
	data  unsafe.Pointer //数据
}

/*
1.不含方法
type eface struct {
	_type *_type  // 类型
	data  unsafe.Pointer  //数据
}
2. 含方法
type iface struct {
	tab  *itab  // tab 是接口表指针，指向类型信息  --->动态类型
	data unsafe.Pointer // 数据指针，则指向具体的数据 --> 动态值
}
interface与nil进行比较比较的是结构体中指向类型的指针。当指向类型的指针为nil时，interface才为nil。
*/

func main() {
	// 不带方法， 用int类型作为例子
	var i interface{}
	fmt.Println(i == nil) // true

	(*eface)(unsafe.Pointer(&i)).data = unsafe.Pointer(uintptr(0x1))
	fmt.Println(i == nil) // true

	(*eface)(unsafe.Pointer(&i))._type = (*int)(unsafe.Pointer(uintptr(0x1)))
	fmt.Println(i == nil) // false

}
