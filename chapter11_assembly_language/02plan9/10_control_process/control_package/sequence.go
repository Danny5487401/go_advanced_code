package control_package

/*
顺序执行
	func main() {
		var a = 10
		println(a)

		var b = (a+a)*a
		println(b)
	}
-->汇编思维
	func main() {
		var a, b int

		a = 10
		runtime.printint(a)
		runtime.printnl()

		b = a
		b += b
		b *= a
		runtime.printint(b)
		runtime.printnl()
	}
// 顺序执行
// 一：未优化前的代码
//1.计算函数栈帧的大小。因为函数内有a、b两个int类型变量，同时调用的runtime·printint函数参数是一个int类型并且没有返回值，因此main函数的栈帧是3个int类型组成的24个字节的栈内存空间
	TEXT ·main(SB), $24-0

		// 2.开始处先将变量初始化为0值，其中a-8*2(SP)对应a变量、a-8*1(SP)对应b变量（因为a变量先定义，因此a变量的地址更小)
		MOVQ $0, a-8*2(SP) // a = 0
		MOVQ $0, b-8*1(SP) // b = 0

		// 将新的值写入a对应内存
		MOVQ $10, AX       // AX = 10
		MOVQ AX, a-8*2(SP) // a = AX

		// 3.为了输出a变量，需要将AX寄存器的值放到0(SP)位置
		// 以a为参数调用函数
		MOVQ AX, 0(SP)
		CALL runtime·printint(SB)
		CALL runtime·printnl(SB)

		// 4.在调用函数返回之后，全部的寄存器将被视为可能被调用的函数修改，因此我们需要从a、b对应的内存中重新恢复寄存器AX和BX
		// 函数调用后, AX/BX 寄存器可能被污染, 需要重新加载
		MOVQ a-8*2(SP), AX // AX = a
		MOVQ b-8*1(SP), BX // BX = b

		// 计算b值, 并写入内存
		MOVQ AX, BX        // BX = AX  // b = a
		ADDQ BX, BX        // BX += BX // b += a
		// 没有使用MULQ指令的原因是MULQ指令默认使用AX保存结果
		IMULQ AX, BX       // BX *= AX // b *= a
		MOVQ BX, b-8*1(SP) // b = BX

		// 以b为参数调用函数
		MOVQ BX, 0(SP)
		CALL runtime·printint(SB)
		CALL runtime·printnl(SB)

		RET


*/
