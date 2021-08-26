#include "textflag.h"

GLOBL ·Float32Value(SB),NOPTR,$4
DATA ·Float32Value+0(SB)/4,$1.5      // var float32Value = 1.5

GLOBL ·Float64Value(SB),NOPTR,$8
DATA ·Float64Value(SB)/8,$0x01020304 // bit 方式初始化

