package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {

	s := []int{10, 20, 30}
	sliceInfo1 := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	fmt.Println("函数调用前指向的数据地址", sliceInfo1.Data)

	changeSlice(s)

	sliceInfo4 := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	fmt.Println("函数调用前指向的数据地址", sliceInfo4.Data)
}

func changeSlice(s []int) {
	sliceInfo2 := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	fmt.Println("扩容前指向的数据地址", sliceInfo2.Data)

	for i := 0; i < 10; i++ {
		s = append(s, i)
	}
	sliceInfo3 := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	fmt.Println("扩容后指向的数据地址", sliceInfo3.Data)

}

// 运行结果
//函数调用前指向的数据地址 1374389948152
//扩容前指向的数据地址 1374389948152
//扩容后指向的数据地址 1374390779904
//函数调用前指向的数据地址 1374389948152
