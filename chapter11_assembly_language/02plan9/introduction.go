package _2plan9
/*
Golang 的汇编是基于 Plan9 汇编的

1. 通用寄存器
下面是通用通用寄存器的名字在 IA64 和 plan9 中的对应关系：

	IA64	RAX	RBX	RCX	RDX	RDI	RSI	RBP	RSP	R8	R9	R10	R11	R12	R13	R14	RIP
	Plan9	AX	BX	CX	DX	DI	SI	BP	SP	R8	R9	R10	R11	R12	R13	R14	PC

Plan9 汇编的操作数方向和 Intel 汇编相反的，与 AT&T 类似。

2. 伪寄存器：
SB-> Static base pointer: global symbols.是一个虚拟寄存器，保存了静态基地址(static-base) 指针，即我们程序地址空间的开始地址；
NOSPLIT：向编译器表明不应该插入 stack-split 的用来检查栈需要扩张的前导指令；
FP->Frame pointer: arguments and locals. 使用形如 symbol+offset(FP) 的方式，引用函数的输入参数；
SP->Stack pointer: top of stack. plan9 的这个 SP 寄存器指向当前栈帧的局部变量的开始位置，使用形如 symbol+offset(SP) 的方式，
	引用函数的局部变量，注意：这个寄存器与上文的寄存器是不一样的，这里是伪寄存器，而我们展示出来的都是硬件寄存器.
	所以区分 SP 到底是指硬件 SP 还是指虚拟寄存器，需要以特定的格式来区分。eg：symbol+offset(SP) 则表示伪寄存器 SP。
	eg：offset(SP) 则表示硬件 SP
PC-> Program counter: jumps and branches.

Note:    virtual_mem_distribution 虚拟内存分布图




1. 静态数据区：存放的是全局变量与常量。这些变量的地址编译的时候就确定了（这也是使用虚拟地址的好处，如果是物理地址，这些地址编译的时候是不可能确定的）。
	Data 与 BSS 都属于这一部分。这部分只有程序中止（kill 掉、crasg 掉等）才会被销毁。
	a. BSS段->BSS segment:通常是指用来存放程序中未初始化的全局变量的一块内存区域。BSS是英文BlockStarted by Symbol的简称。
		BSS段属于静态内存分配。

	b. 数据段-> DATA segment通常是指用来存放程序中已初始化的全局变量的一块内存区域。数据段属于静态内存分配。

2. 代码区Text ->code segment/text segment：存放的就是我们编译后的机器码，一般来说这个区域只能是只读。

3. 栈区->stack：主要是 Golang 里边的函数、方法以及其本地变量存储的地方。这部分伴随函数、方法开始执行而分配，运行完后就被释放，
	特别注意这里的释放并不会清空内存。还有一个点需要记住栈一般是从高地址向低地址方向分配，
	换句话说：高地址属于栈底，低地址属于栈顶，它分配方向与堆是相反的。

4. 堆区->heap：像 C/C++ 语言，堆完全是程序员自己控制的。但是 Golang 里边由于有 GC 机制，我们写代码的时候并不需要关心内存是在栈还是堆上分配。
	Golang 会自己判断如果变量的生命周期在函数退出后还不能销毁或者栈上资源不够分配等等情况，就会被放到堆上。堆的性能会比栈要差一些。



逃逸分析：
如果变量被分配到栈上，会伴随函数调用结束自动回收，并且分配效率很高；其次分配到堆上，则需要 GC 进行标记回收。所谓逃逸就是指变量从栈上逃到了堆上。
go 也提供了更方便的命令来进行逃逸分析：go build -gcflags="-m"，如果真的是做逃逸分析，建议使用该命令，别折腾用汇编
 */


/* 使用
1. 变量声明
使用 DATA 结合 GLOBL 来定义一个变量。
	DATA	symbol+offset(SB)/width, value
使用 GLOBL 指令将变量声明为 global，额外接收两个参数，一个是 flag，另一个是变量的总大小。
	GLOBL divtab(SB), RODATA, $8

GLOBL 必须跟在 DATA 指令之后，下面是一个定义了多个 readonly 的全局变量的完整例子：
	DATA pi+0(SB)/8, $3.1415926
	GLOBL pi(SB), RODATA, $8


2. 函数声明
// 该声明一般写在任意一个 .go 文件中，例如：add.go
func add(a, b int) int

// 函数实现
// 该实现一般写在与声明同名的 _{Arch}.s 文件中，例如：add_amd64.s
TEXT pkgname·add(SB), NOSPLIT, $0-16
    MOVQ a+0(FP), AX
    MOVQ a+8(FP), BX
    ADDQ AX, BX
    MOVQ BX, ret+16(FP)
    RET

pkgname 包名可以不写，一般都是不写的，可以参考 go 的源码， 另外 add 前的 · 不是 .

                           参数及返回值大小
                                 |
 TEXT pkgname·add(SB),NOSPLIT,$0-16
         |     |               |
        包名  函数名    栈帧大小(局部变量+可能需要的额外调用函数的参数空间的总大小，
                               但不包括调用其它函数时的 ret address 的大小)


以上使用的 RODATA，NOSPLIT flag，还有其他的值，可以参考：https://golang.org/doc/asm#directives，
 */

/* 分析
很多时候我们无法确定一块代码是如何执行的，需要通过生成汇编、反汇编来研究
	// 编译
	go build -gcflags="-S"
	go tool compile -S hello.go
	go tool compile -N -S hello.go // 禁止优化
	// 反编译
	go tool objdump <binary>


 */
