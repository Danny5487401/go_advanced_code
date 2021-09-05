
#include "textflag.h"

TEXT ·ptrToFunc(SB), NOSPLIT, $0-16
	MOVQ ptr+0(FP), AX // AX = ptr
	MOVQ AX, ret+8(FP) // return AX
	RET

TEXT ·asmFunTwiceClosureAddr(SB), NOSPLIT, $0-8
	LEAQ ·asmFunTwiceClosureBody(SB), AX // AX = ·asmFunTwiceClosureBody(SB)
	MOVQ AX, ret+0(FP)                   // return AX
	RET


// 采用NEEDCTXT标志定义的汇编函数表示需要一个上下文环境
//在AMD64环境下是通过DX寄存器来传递这个上下文环境指针，也就是对应FunTwiceClosure结构体的指针
TEXT ·asmFunTwiceClosureBody(SB), NOSPLIT|NEEDCTXT, $0-8
	MOVQ 8(DX), AX
	ADDQ AX   , AX        // AX *= 2
	MOVQ AX   , 8(DX)     // ctx.X = AX
	MOVQ AX   , ret+0(FP) // return AX
	RET

