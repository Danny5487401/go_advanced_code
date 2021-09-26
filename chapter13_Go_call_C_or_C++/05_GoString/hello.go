package main

//extern void SayHello(_GoString_ s);
import "C"
import "fmt"


func main()  {
	C.SayHello("Hello World\n")

}

//export SayHello
func SayHello(s string)  {
	fmt.Print(s)
}