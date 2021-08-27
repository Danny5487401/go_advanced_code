#include "textflag.h"

TEXT ·Foo(SB), $32-0
	MOVQ a-32(SP),      AX // a
	MOVQ b-30(SP),      BX // b
	MOVQ c_data-24(SP), CX // c.Data
	MOVQ c_len-16(SP),  DX // c.Len
	MOVQ c_cap-8(SP),   DI // c.Cap
	RET


// 局部变量中先定义的变量c离伪SP寄存器对应的地址最近，最后定义的变量a离伪SP寄存器最远。有两个因素导致出现这种逆序的结果：
//1.一个从Go语言函数角度理解，先定义的c变量地址要比后定义的变量的地址更大；
//2.另一个是伪SP寄存器对应栈帧的底部，而X86中栈是从高向低生长的，所以最先定义有着更大地址的c变量离栈的底部伪SP更近。