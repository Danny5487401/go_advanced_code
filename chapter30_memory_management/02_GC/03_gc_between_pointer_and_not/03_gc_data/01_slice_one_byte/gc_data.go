package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func bar() []*int {
	t := make([]*int, 8)
	return t
}

func main() {
	t := bar()
	v := reflect.TypeOf(t)

	rtyp, ok := v.(*reflect.Rtype)
	if !ok {
		println("error")
		return
	}

	r := (*rtype)(unsafe.Pointer(rtyp))
	fmt.Printf("%#v\n", *r)
	fmt.Printf("*gcdata = %d\n", *(r.gcdata))
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
