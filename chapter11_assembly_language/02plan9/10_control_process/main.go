package main

import (
	"go_advanced_code/chapter11_assembly_language/02plan9/10_control_process/control_package"
)

func main() {
	//1。If
	println(control_package.If(0, 1, 2))
	// 2。汇编循环
	println(control_package.Sum([]int64{1, 2, 3, 4, 5}))
	println(control_package.LoopAdd(10, 0, 2))
}
