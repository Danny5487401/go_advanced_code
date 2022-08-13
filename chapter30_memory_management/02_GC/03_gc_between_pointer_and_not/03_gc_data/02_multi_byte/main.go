package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type S struct { // 起始地址
	a  uint8     // 0
	b  uintptr   // 8
	p1 *uint8    // 16
	c  [3]uint64 // 24
	d  uint32    // 48
	p2 *uint64   // 56
	p3 *uint8    // 64
	e  uint32    // 72
	p4 *uint64   // 80
}

func foo() *S {
	t := new(S)
	return t
}

func main() {
	t := foo()
	println("字节数目", unsafe.Sizeof(*t)) // 88
	typ := reflect.TypeOf(t)
	rtyp, ok := typ.Elem().(*reflect.Rtype)

	if !ok {
		println("error")
		return
	}
	fmt.Printf("runtime 类型信息：%#v\n", *rtyp)

	r := (*rtype)(unsafe.Pointer(rtyp))
	fmt.Printf("runtime 类型信息：%#v\n", *r)
	fmt.Printf("gc需要的byte数目: %d \n", *(r.gcdata))
	gcdata1 := (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(r.gcdata)) + 1))
	fmt.Printf("%d\n", *gcdata1)
}

type tflag uint8
type nameOff int32 // offset to a name
type typeOff int32 // offset to an *rtype

// gcdata也是非导出字段并且是一个指针，我们要想对其解引用，我们这里又在本地定义了一个本地rtype类型，用于输出gcdata指向的内存的值。
type rtype struct {
	size       uintptr
	ptrdata    uintptr // number of bytes in the type that can contain pointers
	hash       uint32  // hash of type; avoids computation in hash tables
	tflag      tflag   // extra type information flags
	align      uint8   // alignment of variable with this type
	fieldAlign uint8   // alignment of struct field with this type
	kind       uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal     func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata    *byte   // garbage collection data
	str       nameOff // string form
	ptrToThis typeOff // type for pointer to this type, may be zero
}
