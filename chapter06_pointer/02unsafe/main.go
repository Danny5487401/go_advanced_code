package main

import (
	"fmt"
	"unsafe"
)

type Programmer struct {
	Name     string //名字
	Language string //爱好
}

func main() {
	// 一。结构体操作
	StructOperation()

	// 二。切片操作
	SliceOperation()

	// 三。获取map的长度
	MapOperation()
}

func StructOperation() {
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
}

func SliceOperation() {
	// 有一个内存分配相关的事实：结构体会被分配一块连续的内存，结构体的地址也代表了第一个成员的地址
	/* runtime/slice.go
	type slice struct{
		array unsafe.Pointer
		len int
		cap int

	}
	func makeslice() slice  返回slice 结构体
	*/
	s := make([]int, 9, 20)
	var len1 = *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + unsafe.Sizeof(int(0))))
	fmt.Println("长度", len1, len(s))
	var cap1 = *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + uintptr(16)))
	fmt.Println("容量", cap1, cap(s))
	// 转换过程 Len: &s => pointer => uintptr => pointer => *int => int
}

func MapOperation() {
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
