package main

import (
	"go_advanced_code/chapter11_assembly_language/02plan9/pkg_int/int_package"
)

func main() {
	println(int_package.Id)
}

// 最后一行的空行是必须的，否则可能报 unexpected EOF
