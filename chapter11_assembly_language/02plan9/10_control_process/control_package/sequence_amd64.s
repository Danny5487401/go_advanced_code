
// 顺序执行
// 优化后的代码
// 我们并不需要a、b两个临时变量分配两个内存空间，而且也不需要在每个寄存器变化之后都要写入内存


//1.栈帧大小从24字节减少到16字节。唯一需要保存的是a变量的值
TEXT ·main(SB), $16-0
	// var temp int

	// 将新的值写入a对应内存
	MOVQ $10, AX        // AX = 10
	MOVQ AX, temp-8(SP) // temp = AX

	// 以a为参数调用函数
	CALL runtime·printint(SB)
	CALL runtime·printnl(SB)

	// 函数调用后, AX 可能被污染, 需要重新加载
	MOVQ temp-8*1(SP), AX // AX = temp

	// 计算b值, 不需要写入内存
	MOVQ AX, BX        // BX = AX  // b = a
	ADDQ BX, BX        // BX += BX // b += a
	IMULQ AX, BX       // BX *= AX // b *= a

    MOVQ BX, 0(SP)
    CALL runtime·printint(SB)
    CALL runtime·printnl(SB)

    RET