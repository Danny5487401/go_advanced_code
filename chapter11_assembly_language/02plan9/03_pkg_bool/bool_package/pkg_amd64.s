#include "textflag.h"
GLOBL ·BoolValue(SB),$1   // 未初始化

GLOBL ·TrueValue(SB),RODATA,$1   // var trueValue = true
DATA ·TrueValue(SB)/1,$1  // 非 0 均为 true

GLOBL ·FalseValue(SB),RODATA,$1  // var falseValue = true
DATA ·FalseValue(SB)/1,$0

