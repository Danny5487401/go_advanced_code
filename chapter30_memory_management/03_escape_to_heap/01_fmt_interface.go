package main

import "fmt"

func main() {
	x := 42
	fmt.Println(x)
}

// go build -gcflags '-m -m'  01_fmt_interface.go

/*
./01_fmt_interface.go:7:13: inlining call to fmt.Println
./01_fmt_interface.go:7:13: x escapes to heap
./01_fmt_interface.go:7:13: []interface {}{...} does not escape
<autogenerated>:1: .this does not escape
*/
