package main

/*
背景：部门产品业务功能采用Golang开发，但是有些功能是用c写的，比如说net-snmp，bfd协议等等，
	像这些如果使用GO语言重编的话，既有实现的复杂度也需要相当长的时间，好在GO语言提供了CGO机制，使得能够在go代码中直接调用C的库函数，
	大大提高了效率，减少了重复开发工作,此外还支持在C语言中调用GO函数，这一点还是蛮强大的

*/

// Go语言调用C函数例子

// 注意： 引用的C头文件需要在注释中声明，紧接着注释需要有import "C"，且这一行和注释之间不能有空格

/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
void myprint(char* s) {
	printf("%s\n", s);
}
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

func main() {

	// func C.CString(string) *C.char  将go的字符串转换为C语言的char*类型。
	cs := C.CString("Hello World\n")

	C.myprint(cs)

	// C.CString 返回的空间由C语言的malloc分配，使用完毕后需要用free释放。
	//  C语言的free参数是void*类型，对应go语言的unsafe.Pointer。
	defer C.free(unsafe.Pointer(cs))

	fmt.Println("call C.sleep for 3s")
	C.sleep(3)
	// 当前进程执行的cgo调用次数
	println("当前进程调用c方法的次数:", runtime.NumCgoCall())
	return
}
