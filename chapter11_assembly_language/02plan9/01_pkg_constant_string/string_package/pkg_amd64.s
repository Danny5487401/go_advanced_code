#include "textflag.h"

// 在Go汇编语言中，内存是通过SB伪寄存器定位。SB是Static base pointer的缩写，意为静态内存的开始地址.
// 我们可以将SB想象为一个和内容容量有相同大小的字节数组，所有的静态全局符号通常可以通过SB加一个偏移量定位，而我们定义的符号其实就是相对于SB内存开始地址偏移量。

// 对于SB伪寄存器，全局变量和全局函数的符号并没有任何区别。
// 其中$·NameData(SB)也是以$美元符号为前缀，因此也可以将它看作是一个常量，它对应的是NameData包变量的地址。
//  在汇编指令中，我们也可以通过LEA指令来获取NameData变量的地址
//   其实Go汇编语言中定义的数据并没有所谓的类型，每个符号只不过是对应一块内存而已，因此NameData符号也是没有类型的。
//   当Go语言的垃圾回收器在扫描到NameData变量的时候，无法知晓该变量内部是否包含指针，通过给NameData变量增加一个NOPTR标志，表示其中不会包含指针数据可以修复该错误：

// GLOBL ·Name(SB),$4
// 其中符号·Name以中点开头表示是当前包的变量，最终符号名为被展开为path/to/pkg.Name
GLOBL ·NameData(SB),NOPTR,$8

// DATA symbol+offset(SB)/width, value
DATA  ·NameData(SB)/8,$"gopher"

GLOBL ·Name(SB),NOPTR,$16
DATA  ·Name+0(SB)/8,$·NameData(SB)
DATA  ·Name+8(SB)/8,$6

