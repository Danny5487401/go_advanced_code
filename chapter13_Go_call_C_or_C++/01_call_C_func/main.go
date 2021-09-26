package main

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

}
/*
分析：
	需要将Go的字符串传入C语言时，先通过C.CString将Go语言字符串对应的内存数据复制到新创建的C语言内存空间上。上面例子的处理思路虽然是安全的，
	但是效率极其低下（因为要多次分配内存并逐个复制元素），同时也极其繁琐
 */
