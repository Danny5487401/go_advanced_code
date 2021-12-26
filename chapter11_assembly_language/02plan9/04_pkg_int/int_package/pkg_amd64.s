#include "textflag.h"

GLOBL ·Id(SB),NOPTR,$8
DATA ·Id+0(SB)/1,$0x37
DATA ·Id+1(SB)/1,$0x25
DATA ·Id+2(SB)/1,$0x00
DATA ·Id+3(SB)/1,$0x00
DATA ·Id+4(SB)/1,$0x00
DATA ·Id+5(SB)/1,$0x00
DATA ·Id+6(SB)/1,$0x00
DATA ·Id+7(SB)/1,$0x00

GLOBL ·Int32Value(SB),NOPTR,$4
DATA ·Int32Value+0(SB)/1,$0x01  // 第0字节
DATA ·Int32Value+1(SB)/1,$0x02  // 第1字节
DATA ·Int32Value+2(SB)/2,$0x03  // 第3-4字节

GLOBL ·Uint32Value(SB),NOPTR,$4
DATA ·Uint32Value(SB)/4,$0x01020304 // 第1-4字节

// 最后一行的空行是必须的，否则可能报 unexpected EOF
