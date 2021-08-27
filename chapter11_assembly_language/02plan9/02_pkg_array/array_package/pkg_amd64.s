#include "textflag.h"
// 汇编中数组也是一种非常简单的类型。Go语言中数组是一种有着扁平内存结构的基础类型。因此[2]byte类型和[1]uint16类型有着相同的内存结构
// 汇编代码中并不需要NOPTR标志，因为Go编译器会从Go语言语句声明的[2]int类型中推导出该变量内部没有指针数据
GLOBL ·Num(SB),RODATA,$16
DATA ·Num+0(SB)/8,$0
DATA ·Num+8(SB)/8,$0

