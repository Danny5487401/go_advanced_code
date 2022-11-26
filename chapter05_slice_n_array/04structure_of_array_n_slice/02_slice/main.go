package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	// slice占用的内存大小为24byte
	var s = []int{1, 2, 3}
	fmt.Println(unsafe.Sizeof(s))

	sliceInfo := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	fmt.Println(sliceInfo.Data, sliceInfo.Len, sliceInfo.Cap)
	arr := (*[3]int)(unsafe.Pointer(sliceInfo.Data))
	for i := 0; i < len(s); i++ {
		fmt.Printf("%d ", arr[i])
	}
	fmt.Println()

}
