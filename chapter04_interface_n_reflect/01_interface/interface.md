# interface 源码分析

## 背景：
接口的定义可以说是一种规范，是一组方法的集合，通常在代码设计的层面，对多个组件有共性的方法进行抽象(共性可以分为横向和纵向)引入一层中间层，
解除上游与下游的耦合关系，让代码可读性更高并不用关心方法的具体实现，同时借助接口也可以实现多态。

共性可以分为横向和纵向的
- 纵向：
    例如动物这个对象可以向下细分为狗和猫，它们有共同的行为可以跑。
- 横向
    再或者数据库的连接可以抽象为接口，可以支持mysql、oracle等

## 源码分类

interface的定义在1.15.3源码包runtime中,interface的定义分为两种，
- 不带方法的runtime.eface
- 带方法的runtime.iface

### 1. runtime.eface表示不含方法的interface{}类型

![](.interface_images/eface.png)

结构体包含可以表示任意数据类型的_type和存储指定的数据data,data用指针来表示
```go
type eface struct {
    _type *_type  // 表示空接口所承载的具体的实体类型
    data  unsafe.Pointer
}
type _type struct {
	//  // 类型大小
    size       uintptr //占用的字节大小
    ptrdata uintptr //指针数据 size of memory prefix holding all pointers
    
    hash       uint32 //计算的hash
    tflag      tflag //额外的标记信息,和反射相关
    
    // 内存对齐相关
    align      uint8 //内存对齐系数
    fieldAlign uint8 //字段内存对齐系数
    
    // 类型的编号，有bool, slice, struct 等等等等
    kind uint8 //用于标记数据类型
    // function for comparing objects of this type
    // (ptr to object A, ptr to object B) -> ==?
    equal func(unsafe.Pointer, unsafe.Pointer) bool//用于判断当前类型多个对象是否相等
    str       nameOff //名字偏移量
    ptrToThis typeOff //指针的偏移量
}
```

Go 语言各种数据类型都是在 _type 字段的基础上，增加一些额外的字段来进行管理的：
```go
type arraytype struct {
    typ   _type
    elem  *_type
    slice *_type
    len   uintptr
}

type chantype struct {
    typ  _type
    elem *_type
    dir  uintptr
}

type slicetype struct {
    typ  _type
    elem *_type
}

type structtype struct {
    typ     _type
    pkgPath name
    fields  []structfield
}
```
这些数据类型的结构体定义，是反射实现的基础。


### 2. runtime.iface表示包含方法的接口
![](.interface_images/iface.png)
```go
type iface struct {
    tab  *itab  // tab 是接口表指针，指向类型信息  --->动态类型
    data unsafe.Pointer // 数据指针，则指向具体的数据 --> 动态值
}
```
结构体包含itab和data数据,它们分别被称为动态类型和动态值。而接口值包括动态类型和动态值。

```go
type itab struct {
    inter  *interfacetype // inter 字段则描述了接口的类型
    _type  *_type  // 描述了实体的类型，包括内存对齐方式，大小等
    link   *itab
    hash   uint32 // copy of _type.hash. Used for type switches.
    bad    bool   // type does not implement interface
    inhash bool   // has this itab been added to hash?
    unused [2]byte
    fun    [1]uintptr // variable sized, fun 字段放置和接口方法对应的具体数据类型的方法地址，实现接口调用方法的动态分派，一般在每次给接口赋值发生转换时会更新此表，或者直接拿缓存的 itab。
}
```
itab包含的是
- 接口类型interfacetype
- 装载实体的任意类型_type
- 实现接口的方法fun,fun是可变大小,go在编译期间就会对接口实现校验检查,并将对应的方法存储fun。

你可能会觉得奇怪，为什么 fun 数组的大小为 1，要是接口定义了多个方法可怎么办？实际上，这里存储的是第一个方法的函数指针，如果有更多的方法，在它之后的内存空间里继续存储。
从汇编角度来看，通过增加地址就能获取到这些函数指针，没什么影响。顺便提一句，这些方法是按照函数名称的字典序进行排列的。

Note：这里只会列出实体类型和接口相关的方法，实体类型的其他方法并不会出现在这里。如果你学过 C++ 的话，这里可以类比虚函数的概念。

```go
type interfacetype struct {
    typ     _type
    pkgpath name  // 定义了接口的包名
    mhdr    []imethod // 表示接口所定义的函数列表
}
```
interfacetype包装了 _type 类型，_type 实际上是描述 Go 语言中各种数据类型的结构体。

### 接口类型和 nil 作比较

接口值的零值是指动态类型和动态值都为 nil。当仅且当这两部分的值都为 nil 的情况下，这个接口值就才会被认为 接口值 == nil
```go
type itab struct {
    inter *interfacetype //接口类型的表示
    _type *_type
    hash  uint32 // copy of _type.hash. Used for type switches.
    _     [4]byte
    fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}
```

