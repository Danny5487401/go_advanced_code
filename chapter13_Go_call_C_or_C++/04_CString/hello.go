package main

//extern void SayHello(char* s);
import "C"
import (
	"fmt"
)


func main()  {
	C.SayHello(C.CString("Hello World\n"))


}

//export SayHello
func SayHello(s *C.char)  {
	fmt.Println(C.GoString(s))
}