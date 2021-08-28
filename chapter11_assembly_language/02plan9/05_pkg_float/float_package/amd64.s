#include "textflag.h"


// Go汇编语言通常无法区分变量是否是浮点数类型，与之相关的浮点数机器指令会将变量当作浮点数处理。Go语言的浮点数遵循IEEE754标准，有float32单精度浮点数和float64双精度浮点数之分。
// IEEE754标准中，最高位1bit为符号位，然后是指数位（指数为采用移码格式表示），然后是有效数部分（其中小数点左边的一个bit位被省略）。
GLOBL ·Float32Value(SB),NOPTR,$4
DATA ·Float32Value+0(SB)/4,$1.5      // var float32Value = 1.5

GLOBL ·Float64Value(SB),NOPTR,$8
DATA ·Float64Value(SB)/8,$0x01020304 // bit 方式初始化

