package main

import (
	"fmt"
	"unsafe"
)

type eface struct {
	itab, data uintptr
}

func main() {
	var a interface{} = nil

	var b interface{} = (*int)(nil)

	x := 5
	var c interface{} = (*int)(&x)

	ia := *(*eface)(unsafe.Pointer(&a))
	ib := *(*eface)(unsafe.Pointer(&b))
	ic := *(*eface)(unsafe.Pointer(&c))

	fmt.Println(ia, ib, ic)

	fmt.Println(*(*int)(unsafe.Pointer(ic.data)))
}
