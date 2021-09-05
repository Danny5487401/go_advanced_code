package main

import (
	"go_advanced_code/chapter11_assembly_language/02plan9/14_closure/closure_package"
)

func main() {
	fnTwice := closure_package.NewTwiceFunClosure(1)
	println(fnTwice()) // 1*2 => 2
	println(fnTwice()) // 2*2 => 4
	println(fnTwice()) // 4*2 => 8
}
