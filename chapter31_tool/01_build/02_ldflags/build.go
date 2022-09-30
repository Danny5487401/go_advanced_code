package main

import "fmt"

var buildVer string

func main() {
	fmt.Println("link 传参数为", buildVer)

}

/*
'-s -w': 压缩编译后的体积
-s: 去掉符号表
-w: 去掉调试信息，不能gdb调试了
-X: 设置包中的变量值
*/

// Note: 不用前路径前缀 -ldflags "-X chapter31_tool.01_build.main.buildVer=1.0"
// go build -ldflags "-X main.buildVer=1.0" -o chapter31_tool/01_build/build chapter31_tool/01_build/build.go
