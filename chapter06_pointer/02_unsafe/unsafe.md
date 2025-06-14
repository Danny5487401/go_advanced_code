<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Golang指针](#golang%E6%8C%87%E9%92%88)
  - [为什么有 unsafe](#%E4%B8%BA%E4%BB%80%E4%B9%88%E6%9C%89-unsafe)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [Go 1.17 之前](#go-117-%E4%B9%8B%E5%89%8D)
      - [1. func Sizeof(x ArbitraryType) uintptr](#1-func-sizeofx-arbitrarytype-uintptr)
      - [2. func Offsetof(x ArbitraryType) uintptr](#2-func-offsetofx-arbitrarytype-uintptr)
      - [3. func Alignof(x ArbitraryType) uintptr](#3-func-alignofx-arbitrarytype-uintptr)
      - [总结](#%E6%80%BB%E7%BB%93)
    - [1.17 新变化](#117-%E6%96%B0%E5%8F%98%E5%8C%96)
  - [应用](#%E5%BA%94%E7%94%A8)
    - [1. map 源码中的应用](#1-map-%E6%BA%90%E7%A0%81%E4%B8%AD%E7%9A%84%E5%BA%94%E7%94%A8)
      - [简单应用](#%E7%AE%80%E5%8D%95%E5%BA%94%E7%94%A8)
      - [复杂应用](#%E5%A4%8D%E6%9D%82%E5%BA%94%E7%94%A8)
    - [2. atomic.value中应用](#2-atomicvalue%E4%B8%AD%E5%BA%94%E7%94%A8)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Golang指针



## 为什么有 unsafe
Go 语言类型系统是为了安全和效率设计的，有时，安全会导致效率低下。有了 unsafe 包，高阶的程序员就可以利用它绕过类型系统的低效。
因此，它就有了存在的意义，阅读 Go 源码，会发现有大量使用 unsafe 包的例子。
unsafe.Pointer 是桥梁，可以让任意类型的指针实现相互转换，也可以将任意类型的指针转换为 uintptr 进行指针运算。
unsafe.Pointer 可以让你的变量在不同的普通指针类型转来转去，也就是表示为任意可寻址的指针类型。
而 uintptr 常用于与 unsafe.Pointer 打配合，用于做指针运算

1. unsafe.Pointer   通用指针

- （1）任何类型的指针都可以被转化为Pointer
- （2）Pointer可以被转化为任何类型的指针
- （3）uintptr可以被转化为Pointer
- （4）Pointer可以被转化为uintptr

Note : 我们不可以直接通过*p来获取unsafe.Pointer指针指向的真实变量的值，因为我们并不知道变量的具体类型。
   和普通指针一样，unsafe.Pointer指针也是可以比较的，并且支持和nil常量比较判断是否为空指针


2. uintptr  整数类型

定义: uintptr is an integer type that is large enough to hold the bit pattern of any 03PointerSetPrivateValue

源码：
```go
type uintptr uintptr
```

Note:uintptr 并没有指针的语义，意思就是 uintptr 所指向的对象会被 gc 无情地回收。
而 unsafe.Pointer 有指针语义，可以保护它所指向的对象在“有用”的时候不会被垃圾回收

## 源码分析

### Go 1.17 之前
unsafe包 两个类型，三个函数
```go
type ArbitraryType int
type Pointer *ArbitraryType
```

ArbitraryType是int的一个别名，在Go中对ArbitraryType赋予特殊的意义。代表一个任意Go表达式类型。实际上它类似于 C 语言里的 void*。
Pointer是int指针类型的一个别名，在Go中可以把Pointer类型，理解成任何指针的父类型。

在Go 1.17 之前，unsafe标准库包只提供了三个函数：
- func Alignof(variable ArbitraryType) uintptr。 此函数用来取得一个值在内存中的地址对齐保证（address alignment guarantee）
- func Offsetof(selector ArbitraryType) uintptr。 此函数用来取得一个结构体值的某个字段的地址相对于此结构体值的地址的偏移。
- func Sizeof(variable ArbitraryType) uintptr。 此函数用来取得一个值的尺寸（亦即此值的类型的尺寸）


#### 1. func Sizeof(x ArbitraryType) uintptr
```go
// Sizeof takes an expression x of any type and returns the size in bytes
// of a hypothetical variable v as if v was declared via var v = x.
// The size does not include any memory possibly referenced by x.
// For instance, if x is a slice, Sizeof returns the size of the slice
// descriptor, not the size of the memory referenced by the slice.
// The return value of Sizeof is a Go constant.
func Sizeof(x ArbitraryType) uintptr

```
unsafe.Sizeof接受任意类型的值(表达式)，返回其占用的字节数,这和c语言里面不同，

如果是slice，则不会返回这个slice在内存中的实际占用长度，一个 slice 的大小则为 slice header 的大小.
c语言里面sizeof函数的参数是类型，而这里是一个表达式，比如一个变量。
```C
int a=10;
int arr=[1,2,3];
char str[]="hello";
int len_a = sizeof(a);
int len_arr = sizeof(arr);
int len_str = sizeof(str)
 
printf("len_a=%d,len_arr=%d,len_str=%d\n",len_a,len_arr,len_str)

```
返回类型 x 所占据的字节数，但不包含 x 所指向的内容的大小。例如，对于一个指针，函数返回的大小为 8 字节（64位机上）

#### 2. func Offsetof(x ArbitraryType) uintptr
```go
func Offsetof(x ArbitraryType) uintptr
//unsafe.Offsetof： 返回结构体成员在内存中的位置离结构体起始处的字节数，所传参数必须是结构体的成员
```

#### 3. func Alignof(x ArbitraryType) uintptr
```go
func Alignof(x ArbitraryType) uintptr
//Alignof返回变量对齐字节数量m，Offsetof返回变量指定属性的偏移量，它分配到的内存地址能整除 m.
//这个函数虽然接收的是任何类型的变量，但是有一个前提，就是变量要是一个struct类型，且还不能直接将这个struct类型的变量当作参数，
//只能将这个struct类型变量的属性当作参数
```
规则：
- 对于任意类型的变量 x ，unsafe.Alignof(x) 至少为 1。
- 对于 struct 类型的变量 x，计算 x 每一个字段 f 的 unsafe.Alignof(x.f)，unsafe.Alignof(x) 等于其中的最大值。
```go
type Bar struct {
    x int32 // 4
    y *Foo  // 8
    z bool  // 1
}
// 结构体变量b1的对齐系数
fmt.Println(unsafe.Alignof(b1))   // 8
// b1每一个字段的对齐系数
fmt.Println(unsafe.Alignof(b1.x)) // 4：表示此字段须按4的倍数对齐
fmt.Println(unsafe.Alignof(b1.y)) // 8：表示此字段须按8的倍数对齐
fmt.Println(unsafe.Alignof(b1.z)) // 1：表示此字段须按1的倍数对
```  

- 对于 array 类型的变量 x，unsafe.Alignof(x) 等于构成数组的元素类型的对齐倍数。


#### 总结

三个函数的参数均是ArbitraryType类型，就是接受任何类型的变量,返回的结果都是 uintptr 类型，这和 unsafe.Pointer 可以相互转换。

三个函数都是在编译期间执行，它们的结果可以直接赋给 const型变量。另外，因为三个函数执行的结果和操作系统、编译器相关，所以是不可移植的

### 1.17 新变化 
Go 1.17引入了一个新类型和两个新函数。 此新类型为IntegerType。它的定义如下。 此类型不代表着一个具体类型，它只是表示任意整数类型

Go 1.17引入的两个函数为：
- func Add(ptr Pointer, len IntegerType) Pointer。 此函数在一个（非安全）指针表示的地址上添加一个偏移量，然后返回表示新地址的一个指针。 此函数以一种更正规的形式部分地覆盖了下面将要介绍的使用模式3中展示的合法用法。
- func Slice(ptr *ArbitraryType, len IntegerType) []ArbitraryType。 此函数用来从一个任意（安全）指针派生出一个指定长度的切片

Go 1.20进一步引入了三个函数：
- func String(ptr *byte, len IntegerType) string。 此函数用来从一个任意（安全）byte指针派生出一个指定长度的字符串。
- func StringData(str string) *byte。 此函数用来获取一个字符串底层字节序列中的第一个byte的指针。
- func SliceData(slice []ArbitraryType) *ArbitraryType。 此函数用来获取一个切片底层元素序列中的第一个元素的指针


## 应用
例如，一般我们不能操作一个结构体的未导出成员，但是通过 unsafe 包就能做到。
unsafe 包让我可以直接读写内存，还管你什么导出还是未导出

### 1. map 源码中的应用

#### 简单应用
```go
//mapaccess1、mapassign、mapdelete 函数中，需要定位 key 的位置，会先对 key 做哈希运算。
b := (*bmap)(unsafe.Pointer(uintptr(h.buckets) + (hash&m)*uintptr(t.bucketsize)))
```

h.buckets 是一个 unsafe.Pointer，将它转换成 uintptr，然后加上 (hash&m)*uintptr(t.bucketsize)，
二者相加的结果再次转换成 unsafe.Pointer，最后，转换成 bmap指针，得到 key 所落入的 bucket 位置

#### 复杂应用
```go
// store new key/value at insert position
if t.indirectkey {
   kmem := newobject(t.key)
   *(*unsafe.Pointer)(insertk) = kmem
   insertk = kmem
}
if t.indirectvalue {
   vmem := newobject(t.elem)
   *(*unsafe.Pointer)(val) = vmem
}
typedmemmove(t.key, insertk, key)
```

	这段代码是在找到了 key 要插入的位置后，进行“赋值”操作。insertk 和 val 分别表示 key 和 value 所要“放置”的地址。
	如果 t.indirectkey 为真，说明 bucket 中存储的是 key 的指针，因此需要将 insertk 看成 指针的指针，
	这样才能将 bucket 中的相应位置的值设置成指向真实 key 的地址值，也就是说 key 存放的是指针
注意
    uintptr 并没有指针的语义，意思就是 uintptr 所指向的对象会被 gc 无情地回收。而 unsafe.Pointer 有指针语义，可以保护它所指向的对象在“有用”的时候不会被垃圾回收

###  2. atomic.value中应用

atomic/value.go中定义了一个ifaceWords结构，其中typ和data字段类型就是unsafe.Poniter，
这里使用unsafe.Poniter类型的原因是传入的值就是interface{}类型，使用unsafe.Pointer强转成ifaceWords类型，这样可以把类型和值都保存了下来
```go
// ifaceWords is interface{} internal representation.
type ifaceWords struct {
    typ  unsafe.Pointer
    data unsafe.Pointer
}
// Load returns the value set by the most recent Store.
// It returns nil if there has been no call to Store for this Value.
func (v *Value) Load() (x interface{}) {
    vp := (*ifaceWords)(unsafe.Pointer(v))
    for {
        typ := LoadPointer(&vp.typ) // 读取已经存在值的类型
    /**
    ..... 中间省略
    **/
    // First store completed. Check type and overwrite data.
    if typ != xp.typ { //当前类型与要存入的类型做对比
       panic("sync/atomic: store of inconsistently typed value into Value")
    }
}
```

## 参考

- [非类型安全指针](https://gfw.go101.org/article/unsafe.html)