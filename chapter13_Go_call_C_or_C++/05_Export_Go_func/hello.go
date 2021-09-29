package main

// extern表明变量或者函数是定义在其他其他文件中的

//extern void SayHello(_GoString_ s);
import "C"
import "fmt"


func main()  {
	C.SayHello("Hello World\n")

}

// 导出go函数,main包导出的函数会在_cgo_export.h声明
// 导出 C函数的名字要和 Go 函数的名字保持一致，同时函数的参数和返回值类型要尽量采用 C 语言类型。

//export SayHello
func SayHello(s string)  {
	fmt.Print(s)
}

