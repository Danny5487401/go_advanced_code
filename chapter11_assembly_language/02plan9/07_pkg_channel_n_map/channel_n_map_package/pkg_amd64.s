#include "textflag.h"

// map/channel等类型并没有公开的内部结构，它们只是一种未知类型的指针，无法直接初始化
GLOBL ·M(SB),RODATA,$8  // var m map[string]int
DATA  ·M+0(SB)/8,$0

GLOBL ·Ch(SB),RODATA,$8 // var ch chan int
DATA  ·Ch+0(SB)/8,$0

//其实在runtime包中为汇编提供了一些辅助函数。比如在汇编中可以通过runtime.makemap和runtime.makechan内部函数来创建map和chan变量。辅助函数的签名如
// func makemap(mapType *byte, hint int, mapbuf *any) (hmap map[any]any)
// func makechan(chanType *byte, size int) (hchan chan any)
