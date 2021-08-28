package main

var helloworld = "你好, 世界"

/*
TEXT ·main(SB), $16-0用于定义main函数，其中$16-0表示main函数的帧大小是16个字节（对应string头部结构体的大小，用于给runtime·printstring函数传递参数），
0表示main函数没有参数和返回值。main函数内部通过调用运行时内部的runtime·printstring(SB)函数来打印字符串。然后调用runtime·printnl打印换行符号
Go语言函数在函数调用时，完全通过栈传递调用参数和返回值。先通过MOVQ指令，将helloworld对应的字符串头部结构体的16个字节复制到栈指针SP对应的16字节的空间，然后通过CALL指令调用对应函数。最后使用RET指令表示当前函数返回
*/
