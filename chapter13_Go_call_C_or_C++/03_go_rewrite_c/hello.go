package main

/*
#include "./hello.h"
*/
import "C"
import "fmt"

//Go 语言实现 C 模块
//export SayHello
func SayHello(s *C.char)  {
	fmt.Println(C.GoString(s))
}
