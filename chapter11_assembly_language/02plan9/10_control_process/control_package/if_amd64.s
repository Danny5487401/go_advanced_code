
TEXT ·If(SB), NOSPLIT, $0-32
	MOVQ ok+8*0(FP), CX // ok
	MOVQ a+8*1(FP), AX  // a
	MOVQ b+8*2(FP), BX  // b

	CMPQ CX, $0         // 比较
	JZ   L              // if ok == 0, goto L
	MOVQ AX, ret+24(FP) // return a
	RET

// 在跳转指令中，跳转的目标一般是通过一个标号表示。不过在有些通过宏实现的函数中，更希望通过相对位置跳转，这时候可以通过PC寄存器的偏移量来计算临近跳转的位置
L:
	MOVQ BX, ret+24(FP) // return b
	RET