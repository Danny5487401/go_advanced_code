package main

import (
	"fmt"
	"unsafe"
)

func main() {
	// 基本api
	basicOperation()

}

func basicOperation() {
	// sizeof 占用的字节数 Byte
	fmt.Println(unsafe.Sizeof(true))       // 1
	fmt.Println(unsafe.Sizeof(int8(0)))    // 1
	fmt.Println(unsafe.Sizeof(int16(10)))  // 2
	fmt.Println(unsafe.Sizeof(int(10)))    // 8,cpu是64位的
	fmt.Println(unsafe.Sizeof(int32(190))) //4
	var DannyStr = "Danny"
	fmt.Println(unsafe.Sizeof(DannyStr)) //16
	var arrayInfo = [...]int{1, 3, 4}
	fmt.Println(unsafe.Sizeof(arrayInfo)) //24
	var sliceInfo = []int{1, 3, 4}
	fmt.Println(unsafe.Sizeof(sliceInfo)) //24

	// Offsetof
	user := User{Name: "Danny", Age: 23, gender: true}
	userNamePointer := unsafe.Pointer(&user)

	nNamePointer := (*string)(userNamePointer)
	*nNamePointer = "Joy"

	nAgePointer := (*uint32)(unsafe.Pointer(uintptr(userNamePointer) + unsafe.Offsetof(user.Age)))
	*nAgePointer = 25

	nGender := (*bool)(unsafe.Pointer(uintptr(userNamePointer) + unsafe.Offsetof(user.gender)))
	*nGender = false

	fmt.Printf("u.Name: %s, u.Age: %d,  u.Gender: %v\n", user.Name, user.Age, user.gender)

	// Alignof
	var b bool
	var i8 int8
	var i16 int16
	var i64 int64
	var f32 float32
	var s string
	var m map[string]string
	var p *int32

	fmt.Println(unsafe.Alignof(b))   //1
	fmt.Println(unsafe.Alignof(i8))  //1
	fmt.Println(unsafe.Alignof(i16)) //2
	fmt.Println(unsafe.Alignof(i64)) //8
	fmt.Println(unsafe.Alignof(f32)) //4
	fmt.Println(unsafe.Alignof(s))   //8
	fmt.Println(unsafe.Alignof(m))   //8
	fmt.Println(unsafe.Alignof(p))   //8

}

type User struct {
	Name   string
	Age    uint32
	gender bool // 男:true 女：false 就是举个例子别吐槽我这么用。。。。
}
