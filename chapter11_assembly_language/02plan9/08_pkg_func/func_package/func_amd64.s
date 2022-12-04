#include "textflag.h"



// ·a(SB)，表示该符号需要链接器来帮我们进行重定向(relocation)，如果找不到该符号，会输出 relocation target not found 的错误。
TEXT ·Get(SB), NOSPLIT, $0-8
    MOVQ ·a(SB), AX
    MOVQ AX, ret+0(FP)
    RET

// 交换值：使用FP寄存器
TEXT ·Swap(SB), $0
	MOVQ a+0(FP), AX     // AX = a
	MOVQ b+8(FP), BX     // BX = b
	MOVQ BX, ret0+16(FP) // ret0 = BX
	MOVQ AX, ret1+24(FP) // ret1 = AX
	RET

// 生成[]byte返回：使用FP寄存器
TEXT ·Foo(SB), $0
	MOVQ a+0(FP),       AX // a
	MOVQ b+2(FP),       BX // b
	MOVQ c_dat+8*1(FP), CX // c.Data
	MOVQ c_len+8*2(FP), DX // c.Len
	MOVQ c_cap+8*3(FP), DI // c.Cap
	RET
// 其中a和b参数之间出现了一个字节的空洞，b和c之间出现了4个字节的空洞。出现空洞的原因是要保证每个参数变量地址都要对齐到相应的倍数





