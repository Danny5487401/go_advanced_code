#include "textflag.h"

// 中点 · 比较特殊，是一个 unicode 的中点，该点在 mac 下的输入方法是 option+shift+9


// 函数标识符通过TEXT汇编指令定义，表示该行开始的指令定义在TEXT内存段。
// TEXT语句后的指令一般对应函数的实现，但是对于TEXT指令本身来说并不关心后面是否有指令。
// 因此TEXT和LABEL定义的符号是类似的，区别只是LABEL是用于跳转标号，但是本质上他们都是通过标识符映射一个内存地址。
// 常见的NOSPLIT主要用于指示叶子函数不进行栈分裂，NOSPLIT对应Go语言中的//go:nosplit注释。


// ·a(SB)，表示该符号需要链接器来帮我们进行重定向(relocation)，如果找不到该符号，会输出 relocation target not found 的错误。
TEXT ·Get(SB), NOSPLIT, $0-8
    MOVQ ·a(SB), AX
    MOVQ AX, ret+0(FP)
    RET

TEXT ·Swap(SB), $0
	MOVQ a+0(FP), AX     // AX = a
	MOVQ b+8(FP), BX     // BX = b
	MOVQ BX, ret0+16(FP) // ret0 = BX
	MOVQ AX, ret1+24(FP) // ret1 = AX
	RET


TEXT ·Foo(SB), $0
	MOVQ a+0(FP),       AX // a
	MOVQ b+2(FP),       BX // b
	MOVQ c_dat+8*1(FP), CX // c.Data
	MOVQ c_len+8*2(FP), DX // c.Len
	MOVQ c_cap+8*3(FP), DI // c.Cap
	RET
// 其中a和b参数之间出现了一个字节的空洞，b和c之间出现了4个字节的空洞。出现空洞的原因是要保证每个参数变量地址都要对齐到相应的倍数



