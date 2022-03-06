package main

// #cgo pkg-config: libhello
// #include < stdlib.h >
// #include < hello_world.h >
import "C"
import (
	"fmt"
)

func main() {
	fmt.Println(C.number_add_mod(10, 6, 12))
}
