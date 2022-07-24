package main

//static const char* cs = "hello";
import "C"
import "github.com/Danny5487401/go_advanced_code/chapter13_Go_call_C_or_C++/04_import_other_pkg/cgo_helper"

func main() {
	cgo_helper.PrintCString(C.cs)
}

/*
这段代码是不能正常工作的，因为当前main包引入的C.cs变量的类型是当前main包的cgo构造的虚拟的C包下的*char类型（具体点是*C.char，更具体点是*main.C.char），
它和cgo_helper包引入的*C.char类型（具体点是*cgo_helper.C.char）是不同的。

在Go语言中方法是依附于类型存在的，不同Go包中引入的虚拟的C包的类型却是不同的（main.C不等cgo_helper.C），
这导致从它们延伸出来的Go类型也是不同的类型（*main.C.char不等*cgo_helper.C.char），这最终导致了前面代码不能正常工作
*/
