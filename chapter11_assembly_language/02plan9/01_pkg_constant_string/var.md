<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [变量](#%E5%8F%98%E9%87%8F)
  - [变量声明](#%E5%8F%98%E9%87%8F%E5%A3%B0%E6%98%8E)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 变量

在汇编里所谓的变量，一般是存储在 .rodata 或者 .data 段中的只读值。对应到应用层的话，就是已初始化过的全局的 const、var、static 变量/常量。

## 变量声明

Go汇编中使用 DATA 和 GLOBL 来定义一个变量(常量)。

DATA 用来指定对应内存中的值；
```cgo
DATA    symbol+offset(SB)/width, value
// 在Go汇编语言中，内存是通过SB伪寄存器定位。SB是Static base pointer的缩写，意为静态内存的开始地址。

// 其中当前包中Go语言定义的符号symbol，在汇编代码中对应·symbol，其中“·”中点符号为一个特殊的unicode符号
// 具体的含义是从symbol+offset偏移量开始，width宽度的内存，用value常量对应的值初始化。
```
其中symbol为变量在汇编语言中对应的标识符，offset是符号开始地址的偏移量，而不是相对于全局某个地址的偏移，width是要初始化内存的宽度大小，value是要初始化的值。

GLOBL 汇编指令用于定义名为 symbol 的全局变量，变量对应的内存宽度为 width，内存宽度部分必须用常量初始化。

```cgo
GLOBL ·symbol(SB), width
```
使用 GLOBL 指令将变量声明为 global，额外接收两个参数，一个是 flag，另一个是变量的总大小。

综合使用
```cgo
// 使用 DATA 结合 GLOBL 来定义一个变量
// GLOBL 必须跟在 DATA 指令之后:

// const age int32 = 8 // 8 为4个字节
DATA ·age+0(SB)/4, $8  
GLOBL ·age(SB), RODATA, $4

// const pi float64 = 3.1415926
DATA ·pi+0(SB)/8, $3.1415926 
GLOBL ·pi(SB), RODATA, $8

// var year int32 = 2020
DATA ·year+0(SB)/4, $2020 
GLOBL ·year(SB), RODATA, $4

// const hello string = "hello my world"
DATA ·hello+0(SB)/8, $"hello my" 
DATA ·hello+8(SB)/8, $"   world" 
GLOBL ·hello(SB), RODATA, $16 


// 引入了新的标记<>，这个跟在符号名之后，表示该全局变量只在当前文件中生效，类似于 C 语言中的 static。
DATA ·hello<>+0(SB)/8, $"hello my" 
DATA ·hello<>+8(SB)/8, $"   world" 
GLOBL ·hello<>(SB), RODATA, $16
```


DATA初始化内存时，width必须是1、2、4、8几个宽度之一，因为再大的内存无法一次性用一个uint64大小的值表示，
对于int32类型的count变量来说，我们既可以逐个字节初始化，也可以一次性初始化
```assembly

DATA ·count+0(SB)/1,$1
DATA ·count+1(SB)/1,$2
DATA ·count+2(SB)/1,$3
DATA ·count+3(SB)/1,$4

// or
DATA ·count+0(SB)/4,$0x04030201

```
正如之前所说，所有符号在声明时，其 offset 一般都是 0。有时也可能会想在全局变量中定义数组，或字符串，这时候就需要用上非 0 的 offset 了，例如:
```cgo
DATA bio<>+0(SB)/8, $"oh yes i"
DATA bio<>+8(SB)/8, $"am here "
GLOBL bio<>(SB), RODATA, $16
```
