package _2plan9

// go tool compile -S defer.go
func demo1() int {
	ret := -1
	defer func() {
		ret = 1
	}()
	return ret
}

/*
	// 编译的汇编 demo1 部分代码
"".demo1 STEXT size=158 args=0x8 locals=0x30
        0x0000 00000 (.\defer.go:5)   TEXT    "".demo1(SB), ABIInternal, $48-8

        // 栈的初始化操作，以及 GC相关的标记等等操作，有兴趣的可以自己研究以下。
        ...

        // 15(SP) 不知道为什么要操作这个，望大佬解释，本人猜测可能跟 deferreturn 有关。
        0x002c 00044 (.\defer.go:5)   MOVB    $0, ""..autotmp_3+15(SP)

        // 这里对 56(SP) 地址进行了赋值操作写了个 0，这个位置其实是返回值地址
        0x0031 00049 (.\defer.go:5)   MOVQ    $0, "".~r0+56(SP)

        // 16(SP) 临时变量 ret，将 -1 写入到了栈中。
			0x003a 00058 (.\defer.go:6)   MOVQ    $-1, "".ret+16(SP)

        // 猜测与 deferreturn 有关。
        0x0043 00067 (.\defer.go:7)   LEAQ    "".demo1.func1·f(SB), AX
        0x004a 00074 (.\defer.go:7)   MOVQ    AX, ""..autotmp_4+32(SP)

        // 将 16(SP) 的地址给了 AX 寄存器，这个地址里存的是 -1
        0x004f 00079 (.\defer.go:7)   LEAQ    "".ret+16(SP), AX

        // 将 AX 寄存器里的 16(SP) 的地址给了 24(SP)
        0x0054 00084 (.\defer.go:7)   MOVQ    AX, ""..autotmp_5+24(SP)
        0x0059 00089 (.\defer.go:7)   MOVB    $1, ""..autotmp_3+15(SP)

        // 将 16(SP) 的值给了 AX 寄存器，这个地址里存的是 -1
        0x005e 00094 (.\defer.go:10)  MOVQ    "".ret+16(SP), AX

        // 将 AX 的值给了 56(SP), 56(SP) 上面说过了是返回值地址， 所以当前的返回值是 -1
        // 这里也是最后一次操作 56(SP)，所以最终的返回值是 -1
        0x0063 00099 (.\defer.go:10)  MOVQ    AX, "".~r0+56(SP)
        0x0068 00104 (.\defer.go:10)  MOVB    $0, ""..autotmp_3+15(SP)

        // 24(SP) 的值给了 AX，24(SP) 存储的是 16(SP) 的地址， 也就是临时变量的地址
        0x006d 00109 (.\defer.go:10)  MOVQ    ""..autotmp_5+24(SP), AX

        // 将 AX 的值给了  0(SP)， 也就是将 16(SP) 的地址给了 0(SP)
        // 这里可以 0(SP) 为调用 demo1.func1 的入参
        0x0072 00114 (.\defer.go:10)  MOVQ    AX, (SP)
        0x0076 00118 (.\defer.go:10)  PCDATA  $1, $1

        // 调用 demo1.func1
        0x0076 00118 (.\defer.go:10)  CALL    "".demo1.func1(SB)
        0x007b 00123 (.\defer.go:10)  MOVQ    40(SP), BP
        0x0080 00128 (.\defer.go:10)  ADDQ    $48, SP
        0x0084 00132 (.\defer.go:10)  RET
        0x0085 00133 (.\defer.go:10)  CALL    runtime.deferreturn(SB)
        0x008a 00138 (.\defer.go:10)  MOVQ    40(SP), BP
        0x008f 00143 (.\defer.go:10)  ADDQ    $48, SP
        0x0093 00147 (.\defer.go:10)  RET
        0x0094 00148 (.\defer.go:10)  NOP
        0x0094 00148 (.\defer.go:5)   PCDATA  $1, $-1
        0x0094 00148 (.\defer.go:5)   PCDATA  $0, $-2
        0x0094 00148 (.\defer.go:5)   CALL    runtime.morestack_noctxt(SB)
        0x0099 00153 (.\defer.go:5)   PCDATA  $0, $-1
        0x0099 00153 (.\defer.go:5)   JMP     0

"".demo1.func1 STEXT nosplit size=13 args=0x8 locals=0x0
        // 这里的 $0-8 就是只有一个参数没有返回值， go 代码中 defer 后面的函数
        0x0000 00000 (.\defer.go:8)   TEXT    "".demo1.func1(SB), NOSPLIT|ABIInternal, $0-8
        0x0000 00000 (.\defer.go:8)   FUNCDATA        $0, gclocals·1a65e721a2ccc325b382662e7ffee780(SB)
        0x0000 00000 (.\defer.go:8)   FUNCDATA        $1, gclocals·69c1753bd5f81501d95132d08af04464(SB)

        // 将 8(SP) 的值给了 AX 寄存器，也就是将 16(SP) 的地址给了 AX
        0x0000 00000 (.\defer.go:9)   MOVQ    "".&ret+8(SP), AX

        // 将 1 给了 AX 寄存器保存的地址的位置上。这个操作像 *a = 1
        0x0005 00005 (.\defer.go:9)   MOVQ    $1, (AX)
        0x000c 00012 (.\defer.go:10)  RET
*/
