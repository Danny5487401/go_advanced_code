<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [makefunc](#makefunc)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## makefunc

提供了函数反射调用的能力

```go
// MakeFunc返回一个给定类型的新函数，该函数封装了函数fn。调用时，该新函数执行以下操作：
// - 将其参数转换为一个切片Slice的值。
// - 运行结果：=fn（args）。
// - 将结果作为一个切片Slice的值返回，每个正式结果一个值。

// 实现fn可以假设参数值切片具有typ给定的参数数量和类型。
// 如果typ描述可变函数，则最终值本身是表示可变参数的切片，如在可变函数体中。
// fn返回的结果值切片必须具有typ给定的结果数量和类型。

// Value.Call方法允许调用方根据值调用类型化函数；相反，MakeFunc允许调用方根据值实现类型化函数。

// 文档的示例部分演示了如何使用MakeFunc为不同类型构建交换函数
func MakeFunc(typ Type, fn func(args []Value) (results []Value)) Value {
   if typ.Kind() != Func {
      panic("reflect: call of MakeFunc with non-Func type")
   }

   t := typ.common()
   ftyp := (*funcType)(unsafe.Pointer(t))

   // Go func的间接值（虚拟）以获取实际代码地址。
   // Go func值是指向C函数指针的指针。https://golang.org/s/go11func
   dummy := makeFuncStub
   code := **(**uintptr)(unsafe.Pointer(&dummy))

   // makeFuncImpl包含一个堆栈映射，供运行时使用
   _, argLen, _, stack, _ := funcLayout(ftyp, nil)

   impl := &makeFuncImpl{code: code, stack: stack, argLen: argLen, ftyp: ftyp, fn: fn}

   return Value{t, unsafe.Pointer(impl), flag(Func)}
}
```