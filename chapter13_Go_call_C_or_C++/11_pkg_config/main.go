package main

//#cgo pkg-config: libnumber
//#include <stdlib.h>
//#include <number.h>
import "C"
import (
	"fmt"
)

func main() {
	fmt.Println(C.number_add_mod(10, 6, 12))
}
