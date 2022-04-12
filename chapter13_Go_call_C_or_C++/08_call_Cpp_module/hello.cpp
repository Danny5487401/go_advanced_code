#include <iostream>

// 不过为了保证C++语言实现的SayHello函数满足C语言头文件hello.h定义的函数规范，我们需要通过extern "C"语句指示该函数的链接符号遵循C语言的规则。
extern "C" {
    #include "hello.h"
}


void SayHello(const char* s) {
    std::cout << s;
}