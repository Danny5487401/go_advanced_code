package main

//#cgo CFLAGS: -I./number
//#cgo LDFLAGS: -L${SRCDIR}/number -lnumber
//
//#include "number.h"
import "C"
import "fmt"

func main() {
	fmt.Println(C.number_add_mod(10, 5, 12))
}

/*
两个#cgo命令，分别是编译和链接参数。

CFLAGS通过-I./number将number库对应头文件所在的目录加入头文件检索路径.
LDFLAGS通过-L${SRCDIR}/number将编译后number静态库所在目录加为链接库检索路径，-lnumber表示链接libnumber.a静态库。

*/
