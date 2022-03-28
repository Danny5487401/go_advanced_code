package main

import (
	"fmt"
	"unsafe"
)

type Programmer struct {
	Name     string //名字
	Language string //爱好
}

type User struct {
	Name   string
	Age    uint32
	Gender bool // 男:true 女：false 就是举个例子别吐槽我这么用。。。。
}

func main() {
	// 基本操作
	BasicOperation()

	// 一。结构体操作
	StructOperation()

	// 二。切片操作
	SliceOperation()

	// 三。获取map的长度
	MapOperation()
}

func BasicOperation() {
	// sizeof
	fmt.Println(unsafe.Sizeof(true))           //1
	fmt.Println(unsafe.Sizeof(int8(0)))        //1
	fmt.Println(unsafe.Sizeof(int16(10)))      // 2
	fmt.Println(unsafe.Sizeof(int(10)))        // 8,cpu是64位的
	fmt.Println(unsafe.Sizeof(int32(190)))     //4
	fmt.Println(unsafe.Sizeof("Danny"))        //16
	fmt.Println(unsafe.Sizeof([]int{1, 3, 4})) //24

	// Offsetof
	user := User{Name: "Danny", Age: 23, Gender: true}
	userNamePointer := unsafe.Pointer(&user)

	nNamePointer := (*string)(userNamePointer)
	*nNamePointer = "Joy"

	nAgePointer := (*uint32)(unsafe.Pointer(uintptr(userNamePointer) + unsafe.Offsetof(user.Age)))
	*nAgePointer = 25

	nGender := (*bool)(unsafe.Pointer(uintptr(userNamePointer) + unsafe.Offsetof(user.Gender)))
	*nGender = false

	fmt.Printf("u.Name: %s, u.Age: %d,  u.Gender: %v\n", user.Name, user.Age, user.Gender)
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
	fmt.Println(unsafe.Alignof(i16)) // 2
	fmt.Println(unsafe.Alignof(i64)) //8
	fmt.Println(unsafe.Alignof(f32)) //4
	fmt.Println(unsafe.Alignof(s))   //8
	fmt.Println(unsafe.Alignof(m))   //8
	fmt.Println(unsafe.Alignof(p))   //8

}

func StructOperation() {
	p := Programmer{Name: "danny", Language: "Golang"}
	fmt.Println("修改前：", p)
	//获取 name的指针
	name := (*string)(unsafe.Pointer(&p))
	*name = "Joy"
	// offset使用获取language地址
	lang := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + unsafe.Offsetof(p.Language)))
	*lang = "Python"
	fmt.Println("修改后：", p)
	//异常情况 示例
	//... 中间逻辑使personAaddr2指向不合法位置
	//personB = (*Person)(unsafe.Pointer(uintptr(0)))
	//fmt.Println("personB.Age is :", personB.Age)
}

func SliceOperation() {
	// 有一个内存分配相关的事实：结构体会被分配一块连续的内存，结构体的地址也代表了第一个成员的地址
	/* runtime/slice.go
	type slice struct{
		array unsafe.Pointer
		len int
		cap int

	}
	func makeslice() slice  返回slice 结构体
	*/
	s := make([]int, 9, 20)
	var len1 = *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + unsafe.Sizeof(int(0))))
	fmt.Println("长度", len1, len(s))
	var cap1 = *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + uintptr(16)))
	fmt.Println("容量", cap1, cap(s))
	// 转换过程 Len: &s => pointer => uintptr => pointer => *int => int
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
