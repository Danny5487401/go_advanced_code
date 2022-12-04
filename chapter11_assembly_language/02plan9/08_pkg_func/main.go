package main

import (
	"fmt"
	"github.com/Danny5487401/go_advanced_code/chapter11_assembly_language/02plan9/08_pkg_func/func_package"
)

func main() {
	// 在 .s 文件中是可以直接使用 .go 中定义的全局变量
	fmt.Println("返回全局变量", func_package.Get())

	// 交换函数
	var a, b int = func_package.Swap(1, 2)
	fmt.Println("Swap 简单返回值", a, b)
	fmt.Println("复杂返回值", func_package.Foo(true, 2))
}
