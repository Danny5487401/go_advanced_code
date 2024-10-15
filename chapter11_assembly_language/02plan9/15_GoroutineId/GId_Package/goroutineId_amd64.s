// 根据官方的Go汇编语言文档，每个运行的Goroutine结构的g指针保存在当前运行Goroutine的线程局部存储(Thread Local Storage,TLS)中
#include "textflag.h"

// func getg() unsafe.Pointer
TEXT ·getg(SB), NOSPLIT, $0-8
	MOVQ TLS, CX
  MOVQ 0(CX)(TLS*1), AX  // 其实TLS类似线程局部存储的地址，地址对应的内存里的数据才是g指针-->MOVQ (TLS), AX
	MOVQ AX, ret+0(FP)
	RET

