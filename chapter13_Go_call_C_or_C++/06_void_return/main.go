package main

//static void noreturn(){}
import "C"
import "fmt"

func main(){
	// 获取一个 void 类型的 C 函数的返回值
	x,_ := C.noreturn()
	fmt.Printf("%#v\n",x) // main._Ctype_void{}
}
