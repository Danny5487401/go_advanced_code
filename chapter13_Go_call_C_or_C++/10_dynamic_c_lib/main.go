package main

//#cgo CFLAGS: -I./number
//#cgo LDFLAGS: -L./number -lnumber
//#include "number.h"
import "C"
import "fmt"

func main() {
	fmt.Println(C.number_add_mod(10, 6, 12))
}
