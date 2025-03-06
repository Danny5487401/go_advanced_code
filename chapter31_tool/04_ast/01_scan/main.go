package main

import (
	"fmt"
	"go/scanner"
	"go/token"
)

func main() {
	// src is the input that we want to tokenize.
	src := []byte("cos(x) + 1i*sin(x) // Danny")

	// Initialize the scanner.
	var s scanner.Scanner
	fset := token.NewFileSet()                      // positions are relative to fset
	file := fset.AddFile("", fset.Base(), len(src)) // register input "file", 要求 file 的大小等于src
	s.Init(file, src, nil /* no error handler */, scanner.ScanComments)

	// Repeated calls to Scan yield the token sequence found in the input.
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
	}

}

/*
1:1     IDENT   "cos"
1:4     (       ""
1:5     IDENT   "x"
1:6     )       ""
1:8     +       ""
1:10    IMAG    "1i"
1:12    *       ""
1:13    IDENT   "sin"
1:16    (       ""
1:17    IDENT   "x"
1:18    )       ""
1:20    COMMENT "// Danny"
1:28    ;       "\n"

*/
