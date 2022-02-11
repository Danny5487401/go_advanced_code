package main

//#cgo CFLAGS: -I./number
//#cgo LDFLAGS: -L${SRCDIR}/chapter13_Go_call_C_or_C++/10_dynamic_c_lib/number -llibnumber.so
//
//#include "number.h"
import "C"
import "fmt"

func main() {
	fmt.Println(C.number_add_mod(10, 6, 12))
}
