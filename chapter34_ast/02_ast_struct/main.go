package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()
	// 来源可以是 file 或则 src
	f, err := parser.ParseFile(fset, "chapter34_ast/02_ast_struct/hello/hello.go", nil, parser.Mode(0))
	if err != nil {
		fmt.Printf("err = %s", err)
	}
	ast.Print(fset, f)

}
