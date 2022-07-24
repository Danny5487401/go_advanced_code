package main

/*
#include "hello.h"
*/
import "C"

import "fmt"

//Go 语言实现 C 模块
// 通过CGO的//export SayHello指令将Go语言实现的函数SayHello导出为C语言函数

//export SayHello
func SayHello(s *C.char) {
	fmt.Println(C.GoString(s))
}
