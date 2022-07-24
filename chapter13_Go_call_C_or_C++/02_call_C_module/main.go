package main

//#cgo CFLAGS: -I./hello
//#cgo LDFLAGS: -L./ -lhello
//#include "hello.h"
import "C"

func main() {
	C.SayHello(C.CString("Hello World\n"))

}
