#include "textflag.h"

// func output(int) (int, int, int)

// 栈结构
//  	------
//  	ret2 (8 bytes)
//  	------
//  	ret1 (8 bytes)
//  	------
//  	ret0 (8 bytes)
//  	------
//  	arg0 (8 bytes)
//  	------ FP
//  	ret addr (8 bytes)
//  	------
//  	caller BP (8 bytes)
//  	------ pseudo SP
//  	frame content (8 bytes)   内容为int类型
//  	------ hardware SP


TEXT ·Output(SB), $8-48
    MOVQ 24(SP), DX // 不带 symbol，这里的 SP 是硬件寄存器 SP
    MOVQ DX, ret3+24(FP) // 第三个返回值
    MOVQ perhapsArg1+16(SP), BX // 当前函数栈大小 > 0，所以 FP 在 SP 的上方 16 字节处
    MOVQ BX, ret2+16(FP) // 第二个返回值
    MOVQ arg1+0(FP), AX
    MOVQ AX, ret1+8(FP)  // 第一个返回值
    RET

