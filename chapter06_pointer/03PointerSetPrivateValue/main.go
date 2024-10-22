package main

import (
	"fmt"
	"runtime"
	"strconv"
	"unsafe"

	"github.com/Danny5487401/go_advanced_code/chapter06_pointer/03PointerSetPrivateValue/p"
)

// unsafe.pointer用于访问操作结构体的私有变量
// 我是想通过unsafe包来实现对Version的成员i和j赋值，然后通过PutI()和PutJ()来打印观察输出结果。
func main() {

	// 1.分配一段内存(会按类型的零值来清零)，并返回一个指针
	var v *p.Version = new(p.Version)

	// 2. 将指针v转成通用指针，再转成int32指针。这里就看到了unsafe.Pointer的作用了，您不能直接将v转成int32类型的指针，那样将会panic
	var i *int32 = (*int32)(unsafe.Pointer(v))

	// 3. 赋值得解引用
	*i = int32(98)

	// 4. i是int32类型，也就是说i占4个字节。所以j是相对于v偏移了4个字节。Note:我的机器64位，占8个字节
	// 如果是32位CPU就是4个字节，如果是64位就是8个字节，由CPU的位数决定，然后按照公式1字节 = 8位计算
	fmt.Println("cpu架构:", runtime.GOARCH)  // amd64
	fmt.Println("int位数:", strconv.IntSize) //int位数 64
	//		您可以用uintptr(4)或uintptr(unsafe.Sizeof(int32(0)))来做这个事
	//var j *int64 = (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + uintptr(unsafe.Sizeof(int32(5)))))
	var j *int64 = (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + uintptr(8)))

	*j = int64(99)
	v.PrintJ() // i=98
	v.PrintJ() // 如果j=0 不是99，注意字节数
}
