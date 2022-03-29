package main

import (
	"go_advanced_code/chapter11_assembly_language/02plan9/08_pkg_func/func_package"
)

func main() {
	// 在 .s 文件中是可以直接使用 .go 中定义的全局变量
	println(func_package.Get())

	// 交换函数
	println(func_package.Swap(1, 2))
	println(func_package.Foo(true, 2))
}
