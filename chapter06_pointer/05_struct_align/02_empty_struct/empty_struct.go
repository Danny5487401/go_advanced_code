package main

import (
	"fmt"
	"unsafe"
)

type emptyStructAsFirst struct {
	a struct{}
	b int32 // 4
}

func main() {
	fmt.Println("空结构体作为第一个元素", unsafe.Sizeof(emptyStructAsFirst{})) // 4
	fmt.Println("空结构体作为最后一个元素", unsafe.Sizeof(emptyStructAsLast{})) // 8
}

type emptyStructAsLast struct {
	a int32 // 4

	// 当 struct{} 作为结构体最后一个字段时，需要内存对齐。因为如果有指针指向该字段, 返回的地址将在结构体之外，如果此指针一直存活不释放对应的内存，就会有内存泄露的问题（该内存不因结构体释放而释放），所以当struct{}作为结构体成员中最后一个字段时，要填充额外的内存保证安全
	b struct{}
}
