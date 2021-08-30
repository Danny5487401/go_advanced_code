package main

import (
	"go_advanced_code/chapter11_assembly_language/02plan9/10_control_process/control_package"
)

// 汇编循环
func main() {
	println(control_package.Sum([]int64{1, 2, 3, 4, 5}))
}
