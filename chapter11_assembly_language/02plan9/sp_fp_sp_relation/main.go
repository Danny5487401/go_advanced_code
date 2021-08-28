package main

import (
	"fmt"
	"go_advanced_code/chapter11_assembly_language/02plan9/sp_fp_sp_relation/sp_fp_package"
)

// 简单的代码证明伪 SP、伪 FP 和硬件 SP 的位置关系
func main() {
	a, b, c := sp_fp_package.Output(987654321)
	fmt.Println(a, b, c)
}
