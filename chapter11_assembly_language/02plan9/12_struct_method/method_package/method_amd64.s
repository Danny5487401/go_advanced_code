#include "textflag.h"
// func (v MyInt) Twice() int
TEXT ·MyInt·Twice(SB), NOSPLIT, $0-16
	MOVQ a+0(FP), AX   // v
	ADDQ AX, AX        // AX *= 2
	MOVQ AX, ret+8(FP) // return v
	RET


// 增加一个接收参数是指针类型的Ptr方,汇编中的·(*MyInt)·Ptr。不过在Go汇编语言中，星号和小括弧都无法用作函数名字，也就是无法用汇编直接实现接收参数是指针类型的方法