package main

import (
	"fmt"
	"go/parser"
	"go/token"
)

func main() {
	var src string = `
package main

import "fmt"
import https "net/http"
import . "go/token"
import _ "go/parser"

`

	fl := token.NewFileSet()
	pl, _ := parser.ParseFile(fl, "example.go", src, parser.ImportsOnly)
	for idx, spec := range pl.Imports {
		fmt.Printf("%d\t %s\t %#v\n", idx, spec.Name, spec.Path)
	}
}

/*
0        <nil>   &ast.BasicLit{ValuePos:22, Kind:9, Value:"\"fmt\""}
1        https   &ast.BasicLit{ValuePos:41, Kind:9, Value:"\"net/http\""}
2        .       &ast.BasicLit{ValuePos:61, Kind:9, Value:"\"go/token\""}
3        _       &ast.BasicLit{ValuePos:81, Kind:9, Value:"\"go/parser\""}
*/
