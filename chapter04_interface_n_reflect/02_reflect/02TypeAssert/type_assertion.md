<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [go 类型转换](#go-%E7%B1%BB%E5%9E%8B%E8%BD%AC%E6%8D%A2)
  - [一. 断言类型的语法](#%E4%B8%80-%E6%96%AD%E8%A8%80%E7%B1%BB%E5%9E%8B%E7%9A%84%E8%AF%AD%E6%B3%95)
    - [注意](#%E6%B3%A8%E6%84%8F)
    - [另外一种方式：switch 断言方式](#%E5%8F%A6%E5%A4%96%E4%B8%80%E7%A7%8D%E6%96%B9%E5%BC%8Fswitch-%E6%96%AD%E8%A8%80%E6%96%B9%E5%BC%8F)
  - [三. 反射reflect类型断言](#%E4%B8%89-%E5%8F%8D%E5%B0%84reflect%E7%B1%BB%E5%9E%8B%E6%96%AD%E8%A8%80)
  - [四. 性能比较](#%E5%9B%9B-%E6%80%A7%E8%83%BD%E6%AF%94%E8%BE%83)
  - [参考资料](#%E5%8F%82%E8%80%83%E8%B5%84%E6%96%99)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# go 类型转换

go 存在 4 种类型转换分别为：断言、强制、显式、隐式。

通常说的类型转换是指断言，强制在日常不会使用到、显示是基本的类型转换、隐式使用到但是不会注意到。断言、强制、显式三类在 go 语法描述中均有说明，隐式是在日常使用过程中总结出来。

1. 强制类型转换
```go
var f float64
bits = *(*uint64)(unsafe.Pointer(&f))

type ptr unsafe.Pointer
bits = *(*uint64)(ptr(&f))

var p ptr = nil
```
unsafe 强制转换是指针的底层操作了，用 c 的朋友就很熟悉这样的指针类型转换，利用内存对齐才能保证转换可靠，例如 int 和 uint 存在符号位差别，利用 unsafe 转换后值可能不同，但是在内存存储二进制一模一样。

应用：接口类型检测
```go
var _ Context = (*ContextBase)(nil)
```
nil 的类型是 nil 地址值为 0，利用强制类型转换成了 * ContextBase，返回的变量就是类型为 * ContextBase 地址值为 0，然后 Context=xx 赋值如果 xx 实现了 Context 接口就没事，如果没有实现在编译时期就会报错，实现编译期间检测接口是否实现

2. 显式转换
一个显式转换的表达式 T (x) ，其中 T 是一种类型并且 x 是可转换为类型的表达式 T，例如：uint(666)。
```go
int64(222)
[]byte("ssss")

type A int
A(2)
```
在以下任何一种情况下，变量 x 都可以转换成 T 类型：
- x 可以分配成 T 类型。
- 忽略 struct 标签 x 的类型和 T 具有相同的基础类型。
- 忽略 struct 标记 x 的类型和 T 是未定义类型的指针类型，并且它们的指针基类型具有相同的基础类型。
- x 的类型和 T 都是整数或浮点类型。
- x 的类型和 T 都是复数类型。
- x 的类型是整数或 [] byte 或 [] rune，并且 T 是字符串类型。
- x 的类型是字符串，T 类型是 [] byte 或 [] rune。

3. 隐式类型转换

隐式类型转换日常使用并不会感觉到，但是运行中确实出现了类型转换

- 组合间的重新断言类型
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
type ReaderClose interface {
    Reader
    Close() error
}
var rc ReaderClose
r := rc
```
ReaderClose 接口组合了 Reader 接口，但是 r=rc 的赋值时还是类型转换了，go 使用系统内置的函数执行了类型转换。

- 相同类型间赋值
```go
type Handler func()

func NewHandler() Handler {
    return func() {}
}
```

4. 这里主要介绍类型断言

## 一. 断言类型的语法
断言通过判断变量是否可以转换成某一个类型
```go
var s = x.(T)
```
x.(T)，这里x表示一个接口的类型，T表示一个类型（也可为接口类型）。

- 如果 x 不是 nil，且 x 可以转换成 T 类型，就会断言成功，返回 T 类型的变量 s。
- 如果 T 不是接口类型，则要求 x 的类型就是 T
- 如果 T 是一个接口，要求 x 实现了 T 接口


类型断言分两种情况：

- 第一种，如果断言的类型T是一个具体类型，类型断言x.(T)就检查x的动态类型是否和T的类型相同。

  1. 如果这个检查成功了，类型断言的结果是一个类型为T的对象，该对象的值为接口变量x的动态值。换句话说，具体类型的类型断言从它的操作对象中获得具体的值。
  2. 如果检查失败，接下来这个操作会抛出panic，除非用两个变量来接收检查结果，如：f, ok := w.(*os.File)

- 第二种，如果断言的类型T是一个接口类型，类型断言x.(T)检查x的动态类型是否满足T接口。

  1. 如果这个检查成功，则检查结果的接口值的动态类型和动态值不变，但是该接口值的类型被转换为接口类型T。换句话说，对一个接口类型的类型断言改变了类型的表述方式，改变了可以获取的方法集合（通常更大），但是它保护了接口值内部的动态类型和值的部分。
  2. 如果检查失败，接下来这个操作会抛出panic，除非用两个变量来接收检查结果，如：f, ok := w.(io.ReadWriter)

### 注意

如果断言的操作对象x是一个nil接口值，那么不论被断言的类型T是什么这个类型断言都会失败。
我们几乎不需要对一个更少限制性的接口类型（更少的方法集合）做断言，因为它表现的就像赋值操作一样，除了对于nil接口值的情况。

表达式是t,ok:=i.(T)，这个表达式也是可以断言一个接口对象（i）里不是nil，并且接口对象（i）存储的值的类型是 T，如果断言成功，就会返回其类型给t，并且此时 ok 的值 为true，表示断言成功。
如果接口值的类型，并不是我们所断言的 T，就会断言失败，但和第一种表达式不同的是这个不会触发 panic，而是将 ok 的值设为false，表示断言失败，此时t为T的零值。所以推荐使用这种方式，可以保证代码的健壮性


### 另外一种方式：switch 断言方式
```go
switch i := x.(type) {
case nil:
    printString("x is nil")                // type of i is type of x (interface{})
case int:
    printInt(i)                            // type of i is int
case float64:
    printFloat64(i)                        // type of i is float64
case func(int) float64:
    printFunction(i)                       // type of i is func(int) float64
case bool, string:
    printString("type is bool or string")  // type of i is type of x (interface{})
default:
    printString("don't know the type")     // type of i is type of x (interface{})
}
```


## 三. 反射reflect类型断言

从reflect.Value中获取接口interface的信息
```go
realValue := value.Interface().(已知的类型)
```
可以理解为“强制转换”，但是需要注意的时候，转换的时候，如果转换的类型不完全符合，则直接panic


Golang 对类型要求非常严格，类型一定要完全符合,如下两个，一个是*float64，一个是float64，如果弄混，则会panic

1. 从 Value 到实例
    该方法最通用，用来将 Value 转换为空接口，该空接口内部存放具体类型实例
    可以使用接口类型查询去还原为具体的类型
```go
//func (v Value) Interface() （i interface{})
```


## 四. 性能比较


```go
var dst int64

// 空接口类型直接类型断言具体的类型
func Benchmark_efaceToType(b *testing.B) {
 b.Run("efaceToType", func(b *testing.B) {
  var ebread interface{} = int64(666)
  for i := 0; i < b.N; i++ {
   dst = ebread.(int64)
  }
 })
}

// 空接口类型使用TypeSwitch 只有部分类型
func Benchmark_efaceWithSwitchOnlyIntType(b *testing.B) {
 b.Run("efaceWithSwitchOnlyIntType", func(b *testing.B) {
  var ebread interface{} = 666
  for i := 0; i < b.N; i++ {
   OnlyInt(ebread)
  }
 })
}

// 空接口类型使用TypeSwitch 所有类型
func Benchmark_efaceWithSwitchAllType(b *testing.B) {
 b.Run("efaceWithSwitchAllType", func(b *testing.B) {
  var ebread interface{} = 666
  for i := 0; i < b.N; i++ {
   Any(ebread)
  }
 })
}

//直接使用类型转换
func Benchmark_TypeConversion(b *testing.B) {
 b.Run("typeConversion", func(b *testing.B) {
  var ebread int32 = 666

  for i := 0; i < b.N; i++ {
   dst = int64(ebread)
  }
 })
}

// 非空接口类型判断一个类型是否实现了该接口 两个方法
func Benchmark_ifaceToType(b *testing.B) {
 b.Run("ifaceToType", func(b *testing.B) {
  var iface Basic = &User{}
  for i := 0; i < b.N; i++ {
   iface.GetName()
   iface.SetName("1")
  }
 })
}

// 非空接口类型判断一个类型是否实现了该接口 12个方法
func Benchmark_ifaceToTypeWithMoreMethod(b *testing.B) {
 b.Run("ifaceToTypeWithMoreMethod", func(b *testing.B) {
  var iface MoreMethod = &More{}
  for i := 0; i < b.N; i++ {
   iface.Get()
   iface.Set()
   iface.One()
   iface.Two()
   iface.Three()
   iface.Four()
   iface.Five()
   iface.Six()
   iface.Seven()
   iface.Eight()
   iface.Nine()
   iface.Ten()
  }
 })
}

// 直接调用方法
func Benchmark_DirectlyUseMethod(b *testing.B) {
 b.Run("directlyUseMethod", func(b *testing.B) {
  m := &More{
   Name: "asong",
  }
  m.Get()
 })
}

```

```shell
/private/var/folders/sk/m49vysmj3ss_y50cv9cvcn800000gn/T/___gobench_github_com_Danny5487401_go_advanced_code_chapter04_interface_n_reflect_02_reflect_02TypeAssert_02_type_assert_performance -test.v -test.paniconexit0 -test.bench . -test.run ^$
goos: darwin
goarch: arm64
pkg: github.com/Danny5487401/go_advanced_code/chapter04_interface_n_reflect/02_reflect/02TypeAssert/02_type_assert_performance
Benchmark_efaceToType
Benchmark_efaceToType/efaceToType
Benchmark_efaceToType/efaceToType-8      	1000000000	         0.5087 ns/op
Benchmark_efaceWithSwitchOnlyIntType
Benchmark_efaceWithSwitchOnlyIntType/efaceWithSwitchOnlyIntType
Benchmark_efaceWithSwitchOnlyIntType/efaceWithSwitchOnlyIntType-8         	1000000000	         1.134 ns/op
Benchmark_efaceWithSwitchAllType
Benchmark_efaceWithSwitchAllType/efaceWithSwitchAllType
Benchmark_efaceWithSwitchAllType/efaceWithSwitchAllType-8                 	792969007	         1.512 ns/op
Benchmark_TypeConversion
Benchmark_TypeConversion/typeConversion
Benchmark_TypeConversion/typeConversion-8                                 	1000000000	         0.3209 ns/op
Benchmark_ifaceToType
Benchmark_ifaceToType/ifaceToType
Benchmark_ifaceToType/ifaceToType-8                                       	536824480	         2.234 ns/op
Benchmark_ifaceToTypeWithMoreMethod
Benchmark_ifaceToTypeWithMoreMethod/ifaceToTypeWithMoreMethod
Benchmark_ifaceToTypeWithMoreMethod/ifaceToTypeWithMoreMethod-8           	148931016	         8.136 ns/op
Benchmark_DirectlyUseMethod
Benchmark_DirectlyUseMethod/directlyUseMethod
Benchmark_DirectlyUseMethod/directlyUseMethod-8                           	1000000000	         0.0000001 ns/op
PASS
```

结论
* 空接口类型的类型断言代价并不高，与直接类型转换几乎没有性能差异
* 空接口类型使用type switch进行类型断言时，随着case的增多性能会直线下降
* 非空接口类型进行类型断言时，随着接口中方法的增多，性能会直线下降
* 直接进行方法调用要比非接口类型进行类型断言要高效很多

## 参考资料
1. [asong关于interface的类型断言是如何实现](https://segmentfault.com/a/1190000039894161)
2. [go 四种类型转换](https://learnku.com/articles/42797)
