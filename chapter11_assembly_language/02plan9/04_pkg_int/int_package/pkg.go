package int_package

var Id int

var int32Value int32

var uint32Value uint32

// 汇编定义变量时初始化数据并不区分整数是否有符号。只有在CPU指令处理该寄存器数据时，才会根据指令的类型来取分数据的类型或者是否带有符号位

/*
DATA symbol+offset(SB)/width, value
	其中symbol为变量在汇编语言中对应的标识符，offset是符号开始地址的偏移量，width是要初始化内存的宽度大小，value是要初始化的值。
	其中当前包中Go语言定义的符号symbol，在汇编代码中对应·symbol，其中“·”中点符号为一个特殊的unicode符号


	在Go汇编语言中，内存是通过SB伪寄存器定位。SB是Static base pointer的缩写，意为静态内存的开始地址。

	具体的含义是从symbol+offset偏移量开始，width宽度的内存，用value常量对应的值初始化。
	DATA初始化内存时，width必须是1、2、4、8几个宽度之一，因为再大的内存无法一次性用一个uint64大小的值表示
	对于int32类型的count变量来说，我们既可以逐个字节初始化，也可以一次性初始化
	DATA ·count+0(SB)/1,$1
	DATA ·count+1(SB)/1,$2
	DATA ·count+2(SB)/1,$3
	DATA ·count+3(SB)/1,$4

	// or

	DATA ·count+0(SB)/4,$0x04030201


pkg_amd64.s的后缀名表示AMD64环境下的汇编代码文件
*/
