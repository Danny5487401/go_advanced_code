#include "textflag.h"

// 在Go汇编语言中，内存是通过SB伪寄存器定位。SB是Static base pointer的缩写，意为静态内存的开始地址
// 我们可以将SB想象为一个和内容容量有相同大小的字节数组，所有的静态全局符号通常可以通过SB加一个偏移量定位，而我们定义的符号其实就是相对于SB内存开始地址偏移量。
//对于SB伪寄存器，全局变量和全局函数的符号并没有任何区别。

// GLOBL ·count(SB),$4
// 其中符号·count以中点开头表示是当前包的变量，最终符号名为被展开为path/to/pkg.count
GLOBL ·NameData(SB),NOPTR,$8
// DATA symbol+offset(SB)/width, value
DATA  ·NameData(SB)/8,$"gopher"

GLOBL ·Name(SB),NOPTR,$16
DATA  ·Name+0(SB)/8,$·NameData(SB)
DATA  ·Name+8(SB)/8,$6

