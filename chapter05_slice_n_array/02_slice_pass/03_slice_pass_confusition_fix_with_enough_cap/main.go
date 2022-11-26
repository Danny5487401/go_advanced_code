package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {

	s := make([]int, 3, 13)
	s[0] = 10
	s[1] = 20
	s[2] = 30

	fmt.Println("--------main start--------------")
	fmt.Println(s)
	fmt.Println(unsafe.Pointer(&s))
	sptr := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	fmt.Println(unsafe.Pointer(sptr.Data))

	changeSlice(s)

	fmt.Println("--------main after --------------")
	fmt.Println(s)



}

func changeSlice(s []int) {
	fmt.Println("--------changeSlice--------------")
	fmt.Println(s)
	fmt.Println(unsafe.Pointer(&s))
	sptr := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	fmt.Println(unsafe.Pointer(sptr.Data))

	for i := 0; i < 10; i++ {
		s = append(s, i)
	}

	fmt.Println(s)
	fmt.Println(unsafe.Pointer(&s))
	sptr1 := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	fmt.Println(unsafe.Pointer(sptr1.Data))

}

// 运行结果
