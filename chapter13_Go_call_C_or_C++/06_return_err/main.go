package main

/*
#include <errno.h>
static void seterrno(int v){
	errno = v;
}
 */
import "C"
import "fmt"

//任意 C 函数 (即使是 void 函数) 可能会被在多赋值场景下被调用，以获取返回值 (如果有的话)，和 C errno 值 (作为 error 值)。
//如果 C 函数返回类型是 void，那么相应的值可以用 _ 代替。


func main()  {
	_,err := C.seterrno(9527)
	fmt.Println(err)
	
}