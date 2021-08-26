package _1asm_basic_idea

/*

汇编语言（assembly language）:用于电子计算机、微处理器、微控制器或其他可编程器件的低级语言，亦称为符号语言
	1.助记符（Mnemonics）代替机器指令的操作码
	2.用地址符号（Symbol）或标号（Label）代替指令或操作数的地址

汇编器：
	1.windows: Microsoft 宏汇编器（称为 MASM）,TASM（Turbo 汇编器），NASM（Netwide 汇编器）和 MASM32（MASM 的一种变体）
	2.linux: GAS（GNU 汇编器）和 NASM,NASM 的语法与 MASM 的最相似

汇编器和链接器:
	汇编器（assembler）是一种工具程序，用于将汇编语言源程序转换为机器语言
	链接器（linker）把汇编器生成的单个文件组合为一个可执行程序。
	调试器（debugger），使程序员可以在程序运行时，单步执行程序并检查寄存器和内存状态。

汇编语言与机器语言有什么关系:
	机器语言（machine language）是一种数字语言， 专门设计成能被计算机处理器（CPU）理解。所有 x86 处理器都理解共同的机器语言。
	汇编语言（assembly language）包含用短助记符如 ADD、MOV、SUB 和 CALL 书写的语句。汇编语言与机器语言是一对一（one-to-one）的关系：
		每一条汇编语言指令对应一条机器语言指令。寄存器（register）是 CPU 中被命名的存储位置，用于保存操作的中间结果。


优点
	1、因为用汇编语言设计的程序最终被转换成机器指令，故能够保持机器语言的一致性，直接、简捷，并能像机器指令一样访问、控制计算机的各种硬件设备，
		如磁盘、存储器、CPU、I/O端口等。使用汇编语言，可以访问所有能够被访问的软、硬件资源。
	2、目标代码简短，占用内存少，执行速度快，是高效的程序设计语言，经常与高级语言配合使用，以改善程序的执行速度和效率，弥补高级语言在硬件控制方面的不足，应用十分广泛


缺点
	1、汇编语言是面向机器的，处于整个计算机语言层次结构的底层，故被视为一种低级语言，通常是为特定的计算机或系列计算机专门设计的。不同的处理器有不同的汇编语言语法和编译器，编译的程序无法在不同的处理器上执行，缺乏可移植性；
	2、难于从汇编语言代码上理解程序设计意图，可维护性差，即使是完成简单的工作也需要大量的汇编语言代码，很容易产生bug，难于调试；
	3、使用汇编语言必须对某种处理器非常了解，而且只能针对特定的体系结构和处理器进行优化，开发效率很低，周期长且单调。

语言组成
	1. 数据传送指令:通用数据传送指令MOV、条件传送指令CMOVcc、堆栈操作指令PUSH/PUSHA/PUSHAD/POP/POPA/POPAD、
		交换指令XCHG/XLAT/BSWAP、地址或段描述符选择子传送指令LEA/LDS/LES/LFS/LGS/LSS等
	2. 整数和逻辑运算指令:加法指令ADD/ADC、减法指令SUB/SBB、加一指令INC、减一指令DEC、比较操作指令CMP、
		乘法指令MUL/IMUL、除法指令DIV/IDIV、符号扩展指令CBW/CWDE/CDQE、十进制调整指令DAA/DAS/AAA/AAS、逻辑运算指令NOT/AND/OR/XOR/TEST等。
	3. 移位指令:将寄存器或内存操作数移动指定的次数。包括逻辑左移指令SHL、逻辑右移指令SHR、算术左移指令SAL、算术右移指令SAR、循环左移指令ROL、循环右移指令ROR等
	4. 位操作指令:位测试指令BT、位测试并置位指令BTS、位测试并复位指令BTR、位测试并取反指令BTC、位向前扫描指令BSF、位向后扫描指令BSR等
	5. 条件设置指令: 这不是一条具体的指令，而是一个指令簇，包括大约30条指令，用于根据EFLAGS寄存器的某些位状态来设置一个8位的寄存器或者内存操作数。
		比如SETE/SETNE/SETGE等等
	6. 控制转移指令:这部分包括无条件转移指令JMP、条件转移指令Jcc/JCXZ、循环指令LOOP/LOOPE/LOOPNE、过程调用指令CALL、子过程返回指令RET、
		中断指令INTn、INT3、INTO、IRET等。注意，Jcc是一个指令簇，包含了很多指令，用于根据EFLAGS寄存器的某些位状态来决定是否转移；
		INT n是软中断指令，n可以是0到255之间的数，用于指示中断向量号。
	7. 串操作指令:这部分指令用于对数据串进行操作，包括串传送指令MOVS、串比较指令CMPS、串扫描指令SCANS、串加载指令LODS、串保存指令STOS，
		这些指令可以有选择地使用REP/REPE/REPZ/REPNE和REPNZ的前缀以连续操作
	8. 输入输出指令: 这部分指令用于同外围设备交换数据，包括端口输入指令IN/INS、端口输出指令OUT/OUTS


寄存器。CPU 本身只负责运算，不负责储存数据。数据一般都储存在内存之中，CPU 要用的时候就去内存读写数据。但是，CPU 的运算速度远高于内存的读写速度，为了避免被拖慢，CPU 都自带一级缓存和二级缓存。基本上，CPU 缓存可以看作是读写速度较快的内存。

	但是，CPU 缓存还是不够快，另外数据在缓存里面的地址是不固定的，CPU 每次读写都要寻址也会拖慢速度。因此，除了缓存之外，CPU 还自带了寄存器（register），用来储存最常用的数据。也就是说，那些最频繁读写的数据（比如循环变量），都会放在寄存器里面，CPU 优先读写寄存器，再由寄存器跟内存交换数据

寄存器种类：
	EAX
	EBX
	ECX
	EDX

	EDI
	ESI

	EBP
	ESP  前面七个都是通用的。ESP 寄存器有特定用途，保存当前 Stack 的地址。32位 CPU、64位 CPU 这样的名称，其实指的就是寄存器的大小

内存模型：
	1.Heap： 查看heap.png
	用户主动请求而划分出来的内存区域，叫做 Heap（堆）。它由起始地址开始，从低位（地址）向高位（地址）增长。
		Heap 的一个重要特点就是不会自动消失，必须手动释放，或者由垃圾回收机制来回收

	2. Stack
	Stack 是由于函数运行而临时占用的内存区域.Stack 是由内存区域的结束地址开始，从高位（地址）向低位（地址）分配。

math.go
	package main
	import "fmt"

	func add(a, b int) int // 汇编函数声明

	func sub(a, b int) int // 汇编函数声明

	func mul(a, b int) int // 汇编函数声明

	func main() {
		fmt.Println(add(10, 11))
		fmt.Println(sub(99, 15))
		fmt.Println(mul(11, 12))
	}
math.s
	#include "textflag.h" // 因为我们声明函数用到了 NOSPLIT 这样的 flag，所以需要将 textflag.h 包含进来

	// func add(a, b int) int
	TEXT ·add(SB), NOSPLIT, $0-24
		MOVQ a+0(FP), AX // 参数 a
		MOVQ b+8(FP), BX // 参数 b
		ADDQ BX, AX    // AX += BX
		MOVQ AX, ret+16(FP) // 返回
		RET

	// func sub(a, b int) int
	TEXT ·sub(SB), NOSPLIT, $0-24
		MOVQ a+0(FP), AX
		MOVQ b+8(FP), BX
		SUBQ BX, AX    // AX -= BX
		MOVQ AX, ret+16(FP)
		RET

	// func mul(a, b int) int
	TEXT ·mul(SB), NOSPLIT, $0-24
		MOVQ  a+0(FP), AX
		MOVQ  b+8(FP), BX
		IMULQ BX, AX    // AX *= BX
		MOVQ  AX, ret+16(FP)
		RET
		// 最后一行的空行是必须的，否则可能报 unexpected EOF

*/
