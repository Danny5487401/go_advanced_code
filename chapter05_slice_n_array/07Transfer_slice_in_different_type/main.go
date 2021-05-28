package main
// 不同类型的切片间互转

import (
	"fmt"
	"reflect"
	"unsafe"
)


func sliceConvert(origSlice interface{},newSliceType reflect.Type)interface{}  {
	sv := reflect.ValueOf(origSlice)
	if sv.Kind() != reflect.Slice{
		panic(fmt.Sprintf("Invalid origSlice(Non-slice value of type %T)", origSlice))
	}
	if newSliceType.Kind() != reflect.Slice {
		panic(fmt.Sprintf("Invalid newSliceType(non-slice type of type %T)", newSliceType))
	}

	//生成新类型的切片
	newSlice := reflect.New(newSliceType)

	//hdr指向到 新生成切片 的SliceHeader
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(newSlice.Pointer()))

	var newElemSize = int(sv.Type().Elem().Size())/int(newSliceType.Elem().Size())
	//设置SliceHeader的Cap，Len，以及数组的ptr
	hdr.Cap = sv.Cap() * newElemSize
	hdr.Len = sv.Len() * newElemSize
	hdr.Data = uintptr(sv.Pointer())

	return newSlice.Elem().Interface()
}

func main()  {
	var int32Slice = []int32{
		1,2,3,4,5,6,7,8,
	}
	var byteSlice []uint8
	byteSlice = sliceConvert(int32Slice, reflect.TypeOf(byteSlice)).([]uint8)
	fmt.Println(byteSlice)
}