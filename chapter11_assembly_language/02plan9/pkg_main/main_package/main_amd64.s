#include "textflag.h"
// 下面是汇编中常见的几种标识符的使用方式（通常也适用于函数标识符）
//  GLOBL ·pkg_name1(SB),$1
//  GLOBL main·pkg_name2(SB),$1
//  GLOBL my/pkg·pkg_name(SB),$1
// Go汇编中可以定义仅当前文件可以访问的私有标识符（类似C语言中文件内static修饰的变量），以<>为后缀名
// GLOBL file_private<>(SB),$1

TEXT main·main(SB), $16-0
	MOVQ ·helloworld+0(SB), AX; MOVQ AX, 0(SP)
	MOVQ ·helloworld+8(SB), BX; MOVQ BX, 8(SP)
	CALL runtime·printstring(SB)
	CALL runtime·printnl(SB)
	RET


