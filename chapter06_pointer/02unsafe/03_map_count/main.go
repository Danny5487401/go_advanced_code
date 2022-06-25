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

	// 三. 获取map的长度
	MapOperation()
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
		func makemap(t *maptype, hint int64, h *hmap, bucket unsafe.Pointer) *hmap
		我们依然能通过 unsafe.Pointer 和 uintptr 进行转换，得到 hamp 字段的值，只不过，现在 count 变成二级指针
	*/
	mp := make(map[string]int)
	mp["danny"] = 1
	mp["Joy"] = 2
	count := **(**int)(unsafe.Pointer(&mp))
	// 转换过程&mp->pointer->**int->int
	fmt.Println("长度", count, len(mp))
}
