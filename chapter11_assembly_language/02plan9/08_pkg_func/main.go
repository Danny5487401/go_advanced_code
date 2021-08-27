package main

import (
	"go_advanced_code/chapter11_assembly_language/02plan9/08_pkg_func/func_package"
)

func main() {
	println(func_package.Get())
	// 交换函数
	println(func_package.Swap(1, 2))
}
