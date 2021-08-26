package main

import "go_advanced_code/chapter11_assembly_language/02plan9/pkg_cycle/cycle_package"

// 汇编循环
func main() {
	println(cycle_package.Sum([]int64{1, 2, 3, 4, 5}))
}
