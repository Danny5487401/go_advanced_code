#ifndef __HELLO_H__
#define __HELLO_H__

// 在Go1.10中CGO新增加了一个_GoString_预定义的C语言类型，用来表示Go语言字符串
// 下面代码== extern void SayHello(_GoString_ s);

extern void SayHello( char* s);

#endif