package main

import "fmt"

var buildVer string

func main() {
	fmt.Println("link 传参数为", buildVer)

}

// Note: 不用前路径前缀 -ldflags "-X chapter31_tool.01_build.main.buildVer=1.0"
// go build -ldflags "-X main.buildVer=1.0" -o chapter31_tool/01_build/build chapter31_tool/01_build/build.go
