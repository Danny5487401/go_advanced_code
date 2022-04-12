package main

//#cgo CFLAGS: -I./number
//#cgo LDFLAGS: -L"./number" -lnumber
//#include "number.h"
import "C"
import "fmt"

func main() {
	fmt.Println(C.number_add_mod(10, 6, 12))
}

// 需要注意的是，在运行时需要将动态库放到系统能够找到的位置。
//1. 对于windows来说，可以将动态库和可执行程序放到同一个目录，或者将动态库所在的目录绝对路径添加到PATH环境变量中，动态链接库后缀为 .dll
//2. 对于macOS来说，需要设置DYLD_LIBRARY_PATH环境变量，动态链接库.dylib
//3. 而对于Linux系统来说，需要设置LD_LIBRARY_PATH环境变量，动态链接库后缀.so
