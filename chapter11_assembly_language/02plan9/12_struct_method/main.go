package main

import (
	"fmt"

	"github.com/Danny5487401/go_advanced_code/chapter11_assembly_language/02plan9/12_struct_method/method_package"
)

func main() {
	var a method_package.MyInt = 1
	fmt.Println(a.Twice())
}
