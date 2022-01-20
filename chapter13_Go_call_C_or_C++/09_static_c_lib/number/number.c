#include "number.h"

int number_add_mod(int a, int b, int mod) {
    return (a+b)%mod;
}


// gcc -c -o number.o number.c
// ar rcs libnumber.a number.o

// ar命令可以用来创建、修改库，也可以从库中提出单个模块。
// 参数r：在库中插入模块（替换）。当插入的模块名已经在库中存在，则替换同名的模块。
// 参数c：创建一个库。不管库是否存在，都将创建。
// 参数s：创建目标文件索引，这在创建较大的库时能加快时间。