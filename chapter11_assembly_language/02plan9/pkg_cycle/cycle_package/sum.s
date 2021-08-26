#include "textflag.h"

// func sum(sl []int64) int64
TEXT ·Sum(SB), NOSPLIT, $0-32
    MOVQ $0, SI
    MOVQ sl+0(FP), BX // &sl[0], addr of the first elem
    MOVQ sl+8(FP), CX // len(sl)
    INCQ CX           // CX++, 因为要循环 len 次

start:
    DECQ CX       // CX--
    JZ   done
    ADDQ (BX), SI // SI += *BX
    ADDQ $8, BX   // 指针移动
    JMP  start

done:
    // 返回地址是 24 是怎么得来的呢？
    // 可以通过 go tool compile -S math.go 得知
    // 在调用 sum 函数时，会传入三个值，分别为:
    // slice 的首地址、slice 的 len， slice 的 cap
    // 不过我们这里的求和只需要 len，但 cap 依然会占用参数的空间
    // 就是 16(FP)
    MOVQ SI, ret+24(FP)
    RET

