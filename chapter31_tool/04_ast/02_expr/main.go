package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	expr := []byte(`a == 1 && b == 2`)
	fset := token.NewFileSet()
	exprAst, err := parser.ParseExprFrom(fset, "", expr, parser.PackageClauseOnly)
	if err != nil {
		fmt.Println(err)
		return
	}

	ast.Print(fset, exprAst)
}

/*
    0  *ast.BinaryExpr {
    1  .  X: *ast.BinaryExpr {
    2  .  .  X: *ast.Ident {
    3  .  .  .  NamePos: -
    4  .  .  .  Name: "a"
    5  .  .  }
    6  .  .  OpPos: -
    7  .  .  Op: ==
    8  .  .  Y: *ast.BasicLit {
    9  .  .  .  ValuePos: -
   10  .  .  .  Kind: INT
   11  .  .  .  Value: "1"
   12  .  .  }
   13  .  }
   14  .  OpPos: -
   15  .  Op: &&
   16  .  Y: *ast.BinaryExpr {
   17  .  .  X: *ast.Ident {
   18  .  .  .  NamePos: -
   19  .  .  .  Name: "b"
   20  .  .  }
   21  .  .  OpPos: -
   22  .  .  Op: ==
   23  .  .  Y: *ast.BasicLit {
   24  .  .  .  ValuePos: -
   25  .  .  .  Kind: INT
   26  .  .  .  Value: "2"
   27  .  .  }
   28  .  }
   29  }

*/
