package main

// 指针基本使用

import (
	"fmt"
	"reflect"
	"unsafe"
)

/*
	*x 称为 解引用 或者 间接引用.

	*x += *x 是通过借助 x 变量的地址, 来操作 x 对应的空间.

	不管是 x 还是 *p , 我们操作的都是同一个空间.
*/
func double(x *int) {
	*x += *x
	x = nil
}
func main() {
	// 通过指针修改值
	var a = 30
	double(&a)
	fmt.Println(a) //60

	p := &a
	double(p)
	fmt.Println(a, p == nil) //120 false

	// 指针类型转换
	fmt.Println("-----------")
	v1 := uint(12)
	v2 := int(13)

	fmt.Println("指针")
	fmt.Println(reflect.TypeOf(&v1)) //*uint
	fmt.Println(reflect.TypeOf(&v2)) //*int

	//使用unsafe.Pointer进行类型的转换
	q := (*uint)(unsafe.Pointer(&v2)) // &v2 *int-->*unit转换

	fmt.Println(reflect.TypeOf(q)) // *unit
	fmt.Println(*q)                //13值不变

}
