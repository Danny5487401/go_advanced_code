<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [runtime](#runtime)
  - [一. runtime 核心功能包括以下内容:](#%E4%B8%80-runtime-%E6%A0%B8%E5%BF%83%E5%8A%9F%E8%83%BD%E5%8C%85%E6%8B%AC%E4%BB%A5%E4%B8%8B%E5%86%85%E5%AE%B9)
  - [二. 特点：](#%E4%BA%8C-%E7%89%B9%E7%82%B9)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# runtime
## 一. runtime 核心功能包括以下内容:
1. 协程(goroutine)调度(并发调度模型)：linux的调度为CPU找到可运行的线程. 而Go的调度是为M(线程)找到P(内存, 执行票据)和可运行的G.
2. 垃圾回收(GC)
3. 内存分配
4. 使得 golang 可以支持如 pprof、trace、race 的检测
5. map, channel, string等内置类型及反射的实现.
6. 操作系统及CPU相关的操作的封装(信号处理, 系统调用, 寄存器操作, 原子操作等), CGO;

## 二. 特点：
1.与Java, Python不同, Go并没有虚拟机的概念, Runtime也直接被编译 成native code.
2.go对系统调用的指令进行了封装, 可不依赖于glibc
3. 用户代码与Runtime代码在执行的时候并没有明显的界限, 都是函数调用
4. 一些go的关键字被编译器编译成runtime包下的函数.
```css
go-->newproc
new->newobject
make->makeslice,makechan,makemap,makemap_small
<- 代表chansend1
-> 代表chanrecv1

```

