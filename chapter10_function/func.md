<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [函数式编程概论](#%E5%87%BD%E6%95%B0%E5%BC%8F%E7%BC%96%E7%A8%8B%E6%A6%82%E8%AE%BA)
  - [背景](#%E8%83%8C%E6%99%AF)
  - [函数](#%E5%87%BD%E6%95%B0)
  - [函数使用](#%E5%87%BD%E6%95%B0%E4%BD%BF%E7%94%A8)
    - [1. 定义函数类型](#1-%E5%AE%9A%E4%B9%89%E5%87%BD%E6%95%B0%E7%B1%BB%E5%9E%8B)
    - [2. 声明函数类型的变量和为变量赋值](#2-%E5%A3%B0%E6%98%8E%E5%87%BD%E6%95%B0%E7%B1%BB%E5%9E%8B%E7%9A%84%E5%8F%98%E9%87%8F%E5%92%8C%E4%B8%BA%E5%8F%98%E9%87%8F%E8%B5%8B%E5%80%BC)
    - [3. 函数作为其他函数入参](#3-%E5%87%BD%E6%95%B0%E4%BD%9C%E4%B8%BA%E5%85%B6%E4%BB%96%E5%87%BD%E6%95%B0%E5%85%A5%E5%8F%82)
    - [4. 函数作为返回值+动态创建](#4-%E5%87%BD%E6%95%B0%E4%BD%9C%E4%B8%BA%E8%BF%94%E5%9B%9E%E5%80%BC%E5%8A%A8%E6%80%81%E5%88%9B%E5%BB%BA)
    - [5. 匿名函数](#5-%E5%8C%BF%E5%90%8D%E5%87%BD%E6%95%B0)
    - [6. closure 闭包](#6-closure-%E9%97%AD%E5%8C%85)
      - [定义](#%E5%AE%9A%E4%B9%89)
      - [解析](#%E8%A7%A3%E6%9E%90)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 函数式编程概论

## 背景
硬件性能的提升以及编译技术和虚拟机技术的改进，一些曾被性能问题所限制的动态语言开始受到关注，Python、Ruby 和 Lua 等语言都开始在应用中崭露头角。
伴随动态语言的流行，函数式编程也再次进入了我们的视野。
函数式编程是一种编程模型，他将计算机运算看做是数学中函数的计算，并且避免了状态以及变量的概念

## 函数
一等公民：
函数作为变量对待。也就说，函数与变量没有差别，它们是一样的，变量出现的地方都可以替换成函数，并且编译也是可以通过的，没有任何语法问题。

## 函数使用
1. 函数可以定义函数类型
2. 函数可以赋值给变量
3. 高阶函数---可以作为入参也可以作为返回值
4. 动态创建函数
5. 匿名函数
6. 闭包

### 1. 定义函数类型
```go
type Operation func(a,b int) int
// -----Operation :type name类型名称
//-----func(a,b int) int:signature函数签名
func Add func(a,b int) int{
    return a+b
}
//符合函数签名的函数
```

### 2. 声明函数类型的变量和为变量赋值
```go
var op Operation
op = Add
fmt.Println(op(1,2))
```
变量op是Operation类型的，可以把Add作为值赋值给变量op，执行op等价于执行Add。

### 3. 函数作为其他函数入参
```go
type Calculator struct {
    v int
}
func (c Calculator)Do(op Operation,a int){
    c.v = op(c.v,a)
}
func main(){
   var calc Calculator
   calc.Do(add,1)
}
```


### 4. 函数作为返回值+动态创建
```go
type Operation func(b int)int
func Add(b int) Operation{
   addB := func(a int)int{
      return a + b
    }
   return addB
}

type Calculator struct {
    v int
}
func (c Calculator)Do(op Operation){
    c.v = op(c.v)
}
func main(){
   var calc Calculator
   calc.Do(add(1)) //c.v = 1
}
```


### 5. 匿名函数
```go
func(a int)int{}
   func Add(b int) Operation{
   return func(a int)int{
       return a + b
   }
}
```


### 6. closure 闭包

#### 定义
闭包是由函数及其相关引用环境组合而成的实体(即：闭包=函数+引用环境 closure = anonymous function + closure conetxt。)   

#### 解析
闭包只是在形式和表现上像函数，但实际上不是函数。函数是一些可执行的代码，这些代码在函数被定义后就确定了，不会在执行时发生变化。
所以一个函数只有一个实例。闭包在运行时可以有多个实例，不同的引用环境和相同的函数组合可以产生不同的实例。
所谓引用环境是指在程序执行中的某个点所有处于活跃状态的约束所组成的集合。
不同的引用环境和相同的函数组合可以产生不同的实例。

```go
type Operation func(b int) int
func Add(b int) Operation{
   addB := func(a int) int {
       return a + b
   }
   return addB
}
```

比如匿名函数里直接使用了变量b，该匿名函数也是闭包函数。
Note:一个函数可以是匿名函数，但不是闭包函数，因为闭包有时是有副作用的。


闭包是函数和它所引用的环境。那么是不是可以表示为一个结构体呢
```go
type Closure struct {
    F func()() 
    i *int
}
```

查看汇编
```asm
        0x001c 00028 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:7)        MOVD    $type:int(SB), R0
        0x0024 00036 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:7)        PCDATA  $1, $0
        0x0024 00036 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:7)        CALL    runtime.newobject(SB)
        0x0028 00040 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:7)        MOVD    R0, main.&sum-8(SP)
        0x002c 00044 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:7)        MOVD    ZR, (R0)
        0x0030 00048 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        MOVD    $type:noalg.struct { F uintptr; X0 *int }(SB), R0
        0x0038 00056 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        PCDATA  $1, $1
        0x0038 00056 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        CALL    runtime.newobject(SB)
        0x003c 00060 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        MOVD    R0, main..autotmp_3-16(SP)
        0x0040 00064 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        MOVD    $main.adder.func1(SB), R1
        0x0048 00072 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        MOVD    R1, (R0)
        0x004c 00076 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        MOVD    main.&sum-8(SP), R1
        0x0050 00080 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        MOVD    main..autotmp_3-16(SP), R2
        0x0054 00084 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        PCDATA  $0, $-2
        0x0054 00084 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        MOVB    (R2), R27
        0x0058 00088 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        PCDATA  $0, $-1
        0x0058 00088 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        PCDATA  $0, $-3
        0x0058 00088 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        MOVWU   runtime.writeBarrier(SB), R3
        0x0060 00096 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        PCDATA  $0, $-1
        0x0060 00096 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        PCDATA  $0, $-2
        0x0060 00096 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        CBZW    R3, 104
        0x0064 00100 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        JMP     108
        0x0068 00104 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        JMP     128
        0x006c 00108 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        CALL    runtime.gcWriteBarrier2(SB)
        0x0070 00112 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        MOVD    R1, (R25)
        0x0074 00116 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        MOVD    8(R2), R3
        0x0078 00120 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        MOVD    R3, 8(R25)
        0x007c 00124 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        JMP     128
        0x0080 00128 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        MOVD    R1, 8(R2)
        0x0084 00132 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        PCDATA  $0, $-1
        0x0084 00132 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        MOVD    main..autotmp_3-16(SP), R1
        0x0088 00136 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:8)        MOVD    R1, main.innerFunc-24(SP)
        0x008c 00140 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:12)       MOVD    main.innerFunc-24(SP), R0
        0x0090 00144 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:12)       MOVD    R0, main.~r0-32(SP)
        0x0094 00148 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:12)       MOVD    -8(RSP), R29
        0x0098 00152 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:12)       MOVD.P  80(RSP), R30
        0x009c 00156 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:12)       RET     (R30)
        0x00a0 00160 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:12)       NOP
        0x00a0 00160 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:6)        PCDATA  $1, $-1
        0x00a0 00160 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:6)        PCDATA  $0, $-2
        0x00a0 00160 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:6)        MOVD    R30, R3
        0x00a4 00164 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:6)        CALL    runtime.morestack_noctxt(SB)
        0x00a8 00168 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:6)        PCDATA  $0, $-1
        0x00a8 00168 (/Users/python/Downloads/git_download/go_advanced_code/chapter10_function/01_func_application/01_closure/main.go:6)        JMP     0
        0x0000 90 0b 40 f9 ff 63 30 eb c9 04 00 54 fe 0f 1b f8  ..@..c0....T....
        0x0010 fd 83 1f f8 fd 23 00 d1 ff 17 00 f9 00 00 00 90  .....#..........
        0x0020 00 00 00 91 00 00 00 94 e0 23 00 f9 1f 00 00 f9  .........#......
        0x0030 00 00 00 90 00 00 00 91 00 00 00 94 e0 1f 00 f9  ................
        0x0040 01 00 00 90 21 00 00 91 01 00 00 f9 e1 23 40 f9  ....!........#@.
        0x0050 e2 1f 40 f9 5b 00 80 39 1b 00 00 90 63 03 40 b9  ..@.[..9....c.@.
        0x0060 43 00 00 34 02 00 00 14 06 00 00 14 00 00 00 94  C..4............
        0x0070 21 03 00 f9 43 04 40 f9 23 07 00 f9 01 00 00 14  !...C.@.#.......
        0x0080 41 04 00 f9 e1 1f 40 f9 e1 1b 00 f9 e0 1b 40 f9  A.....@.......@.
        0x0090 e0 17 00 f9 fd 83 5f f8 fe 07 45 f8 c0 03 5f d6  ......_...E..._.
        0x00a0 e3 03 1e aa 00 00 00 94 d6 ff ff 17 00 00 00 00  ................
        rel 28+8 t=R_ADDRARM64 type:int+0
        rel 36+4 t=R_CALLARM64 runtime.newobject+0
        rel 48+8 t=R_ADDRARM64 type:noalg.struct { F uintptr; X0 *int }+0
        rel 56+4 t=R_CALLARM64 runtime.newobject+0
        rel 64+8 t=R_ADDRARM64 main.adder.func1+0
        rel 88+8 t=R_ARM64_PCREL_LDST32 runtime.writeBarrier+0
        rel 108+4 t=R_CALLARM64 runtime.gcWriteBarrier2+0
        rel 164+4 t=R_CALLARM64 runtime.morestack_noctxt+0


```

type:noalg.struct { F uintptr; X0 *int }+0


## 参考

