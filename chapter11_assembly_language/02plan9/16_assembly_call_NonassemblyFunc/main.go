package main

import (
	"fmt"

	"github.com/Danny5487401/go_advanced_code/chapter11_assembly_language/02plan9/16_assembly_call_NonassemblyFunc/assembly_package"
)

// 汇编调用非汇编函数
func main() {
	s := assembly_package.Output(10, 13)
	fmt.Println(s)
}
