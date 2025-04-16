<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Generic泛型](#generic%E6%B3%9B%E5%9E%8B)
  - [演进过程](#%E6%BC%94%E8%BF%9B%E8%BF%87%E7%A8%8B)
  - [泛型使用场景](#%E6%B3%9B%E5%9E%8B%E4%BD%BF%E7%94%A8%E5%9C%BA%E6%99%AF)
  - [Go 泛型方案](#go-%E6%B3%9B%E5%9E%8B%E6%96%B9%E6%A1%88)
    - [1. 字典(dictionaries)：单份代码实例，以字典传递类型参数信息](#1-%E5%AD%97%E5%85%B8dictionaries%E5%8D%95%E4%BB%BD%E4%BB%A3%E7%A0%81%E5%AE%9E%E4%BE%8B%E4%BB%A5%E5%AD%97%E5%85%B8%E4%BC%A0%E9%80%92%E7%B1%BB%E5%9E%8B%E5%8F%82%E6%95%B0%E4%BF%A1%E6%81%AF)
    - [2. 模版(stenciling)：为每次调用生成代码实例，即便类型参数相同](#2-%E6%A8%A1%E7%89%88stenciling%E4%B8%BA%E6%AF%8F%E6%AC%A1%E8%B0%83%E7%94%A8%E7%94%9F%E6%88%90%E4%BB%A3%E7%A0%81%E5%AE%9E%E4%BE%8B%E5%8D%B3%E4%BE%BF%E7%B1%BB%E5%9E%8B%E5%8F%82%E6%95%B0%E7%9B%B8%E5%90%8C)
    - [3. 混合方案（GC Shape Stenciling）](#3-%E6%B7%B7%E5%90%88%E6%96%B9%E6%A1%88gc-shape-stenciling)
  - [概念](#%E6%A6%82%E5%BF%B5)
    - [类型形参(Type Parameters)](#%E7%B1%BB%E5%9E%8B%E5%BD%A2%E5%8F%82type-parameters)
    - [类型约束(Type Constraint)](#%E7%B1%BB%E5%9E%8B%E7%BA%A6%E6%9D%9Ftype-constraint)
    - [类型具化（instantiation）与类型推导（type inference）](#%E7%B1%BB%E5%9E%8B%E5%85%B7%E5%8C%96instantiation%E4%B8%8E%E7%B1%BB%E5%9E%8B%E6%8E%A8%E5%AF%BCtype-inference)
    - [从方法集(Method set)到类型集(Type set)](#%E4%BB%8E%E6%96%B9%E6%B3%95%E9%9B%86method-set%E5%88%B0%E7%B1%BB%E5%9E%8B%E9%9B%86type-set)
  - [类型的交集](#%E7%B1%BB%E5%9E%8B%E7%9A%84%E4%BA%A4%E9%9B%86)
    - [案例](#%E6%A1%88%E4%BE%8B)
  - [接口两种类型](#%E6%8E%A5%E5%8F%A3%E4%B8%A4%E7%A7%8D%E7%B1%BB%E5%9E%8B)
    - [基本接口(Basic interface)](#%E5%9F%BA%E6%9C%AC%E6%8E%A5%E5%8F%A3basic-interface)
    - [一般接口(General interface)](#%E4%B8%80%E8%88%AC%E6%8E%A5%E5%8F%A3general-interface)
  - [常见错误](#%E5%B8%B8%E8%A7%81%E9%94%99%E8%AF%AF)
  - [推荐写法](#%E6%8E%A8%E8%8D%90%E5%86%99%E6%B3%95)
  - [参考资料](#%E5%8F%82%E8%80%83%E8%B5%84%E6%96%99)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Generic泛型
Go的泛型指的是在Go的类型声明和函数声明中增加可选的类型形参(type parameters)。

类型形参 (type parameters)要受到类型约束（constraint），Go使用嵌入额外元素的接口类型来定义类型约束，类型约束定义了一组满足约束的类型集(Type Set)

形参 (type parameters)只是类似占位符的东西并没有具体的值，只有我们调用函数传入实参(type arguments)之后才有具体的值。

Go 1.18版本增加了对泛型的支持, 除了语法，外加两个预定义类型：

1. comparable: 是Go语言内置的类型约束，它表示类型的值可以使用==和!=比较大小，这也是map类型的key要求的.
```go
// go1.21.5/src/builtin/builtin.go

// comparable is an interface that is implemented by all comparable types
// (booleans, numbers, strings, pointers, channels, arrays of comparable types,
// structs whose fields are all comparable types).
// The comparable interface may only be used as a type parameter constraint,
// not as the type of a variable.
type comparable interface{ comparable }
```
Note: 结构体中的所有成员变量都是可比较类型才行。类型==判断详情可参考：https://go.dev/ref/spec#Comparison_operators

2. any 任意类型 


增加了两个操作符: 
1. ~ : ~int 这种写法的话，就代表着不光是 int ，以 int 为底层类型的类型也都可用于实例化
2. | : 表示允许两者中的任何一个

```go
func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}
```
泛型在Go语言中增加了三个新的重要内容

- 函数和类型新增对 *类型形参(type parameters)* 的支持。 
- 将接口类型定义为类型集合，包括没有方法的接口类型。 
- 支持类型推导，大多数情况下，调用泛型函数时可省略类型实参(type arguments)。



## 演进过程
- Go官方团队在Go 1.18 Beta 1版本的标准库里因为泛型设计而引入了contraints包，Go官方团队的技术负责人Russ Cox在2022.01.25提议将constraints包从Go标准库里移除，放到x/exp项目下，
  该提议也同Go语言发明者Rob Pike, Robert Griesemer和Ian Lance Taylor做过讨论，得到了他们的同意。
- golang.org/x下所有package的源码独立于Go源码的主干分支，也不在Go的二进制安装包里。如果需要使用golang.org/x下的package，可以使用go get来安装。
- golang.org/x/exp下的所有package都属于实验性质或者被废弃的package，不建议使用。

## 泛型使用场景
在 Ian Lance Taylor 的 When To Use Generics 中列出了泛型的典型使用场景，归结为三种主要情况：

- 使用内置的容器类型，如 slices、maps 和 channels
- 实现通用的数据结构，如 linked list 或 tree
- 编写一个函数，其实现对许多类型来说都是一样的，比如一个排序函数

## Go 泛型方案
Go 的泛型代码是在编译时生成的，而不是在运行时进行类型断言。这意味着泛型代码在编译时就能够获得类型信息，从而保证类型安全性。生成的代码针对具体的类型进行了优化，避免了运行时的性能开销。

### 1. 字典(dictionaries)：单份代码实例，以字典传递类型参数信息
在编译时生成一组实例化的字典，在实例化一个泛型函数的时候会使用字典进行蜡印(stencile).

当为泛型函数生成代码的时候，会生成唯一的一块代码，并且会在参数列表中增加一个字典做参数，就像方法会把receiver当成一个参数传入。字典包含为类型参数实例化的类型信息。

字典在编译时生成，存放在只读的data section中。

当然字段可以当成第一个参数，或者最后一个参数，或者放入一个独占的寄存器。

当然这种方案还有依赖问题，比如字典递归的问题，更重要的是，它对性能可能有比较大的影响，比如一个实例化类型int, x=y可能通过寄存器复制就可以了，但是泛型必须通过memmove

### 2. 模版(stenciling)：为每次调用生成代码实例，即便类型参数相同
同一个泛型函数，为每一个实例化的类型参数生成一套独立的代码，感觉和rust的泛型特化一样。

![](.generic_images/stencile.gif)


比如下面一个泛型方法:
```go
func f[T1, T2 any](x int, y T1) T2 {
    //...
}
```
如果有两个不同的类型实例化的调用：
```go
var a float64 = f[int, float64](7, 8.0)
var b struct{f int} = f[complex128, struct{f int}](3, 1+1i)
```

那么这个方案会生成两套代码：
```go
func f1(x int, y int) float64 {
    //... identical bodies ...
}
func f2(x int, y complex128) struct{f int} {
    //... identical bodies ...
}
```

因为编译f时是不知道它的实例化类型的，只有在调用它时才知道它的实例化的类型，所以需要在调用时编译f。
对于相同实例化类型的多个调用，同一个package下编译器可以识别出来是一样的，只生成一个代码就可以了，但是不同的package就不简单了，这些函数表标记为DUPOK,所以链接器会丢掉重复的函数实现

这种策略需要更多的编译时间，因为需要编译泛型函数多次。因为对于同一个泛型函数，每种类型需要单独的一份编译的代码，如果类型非常多，编译的文件可能非常大，而且性能也比较差。


### 3. 混合方案（GC Shape Stenciling）

啥叫shape?

类型的shape是它对内存分配器/垃圾回收器呈现的方式，包括它的大小、所需的对齐方式、以及类型哪些部分包含指针.

每一个唯一的shape会产生一份代码，每份代码携带一个字典，包含了实例化类型的信息

```go
type a int
type b int
type c = int
```
任何指针类型，或具有相同底层类型(underlying type)的类型，属于同一GCShape组。

这两种方法中哪一种最适合 Go?快速编译很重要，但运行时性能也很重要。为了满足这些要求，Go 团队决定在实现泛型时混合两种方法。

Go 使用单态化，但试图减少需要生成的函数副本的数量。它不是为每个类型创建一个副本，而是为内存中的每个布局生成一个副本：int、float64、Node 和其他所谓的 “值类型” 在内存中看起来都不一样，因此编译器将为所有这些类型生成不同的副本。

与值类型相反，指针和接口在内存中总是有相同的布局。编译器将为指针和接口的调用生成同一个泛型函数的副本。就像虚函数表一样，泛型函数接收指针，因此需要一个表来动态地查找方法地址。在 Go 实现中的字典与虚拟方法表的性能特点相同。

这种混合方法的好处是，你在使用值类型的调用中获得了 Monomorphization 的性能优势，而只在使用指针或接口的调用中付出了 Virtual Method Table 的成本

## 概念
- 类型形参 (Type parameter)
- 类型实参(Type argument)
- 类型形参列表( Type parameter list)
- 类型约束(Type constraint)
- 实例化(Instantiations)
- 泛型类型(Generic type)
- 泛型接收器(Generic receiver)
- 泛型函数(Generic function)

### 类型形参(Type Parameters)

类型形参是在函数声明、方法声明的receiver部分或类型定义的类型参数列表中，声明的（非限定）类型名称。
类型形参在声明中充当了一个未知类型的占位符（placeholder），在泛型函数或泛型类型实例化(instantiation)时，类型形参(Type Parameters)会被一个类型实参（type argument）替换

函数和类型被允许拥有类型形参(Type Parameters)。一个类型形参列表看起来和普通的函数形参列表一样，只是它使用的是中括号[方括号]而不是小括号()。

```go
func GenericFoo[P aConstraint, Q anotherConstraint](x,y P, z Q "P aConstraint, Q anotherConstraint")

```

P，Q是类型形参的名字，也就是类型，aConstraint，anotherConstraint代表类型参数的约束（constraint），我们可以理解为对类型参数可选值的一种限定 

```go
func F[T any](p T) { ... }
```
声明的类型参数可以在函数的参数和函数体中使用。 在这个例子中，T是类型参数的名字，也就是类型，any是类型参数的约束，是对类型参数可选类型的约束。但是T的类型要等到泛型函数具化时才能确定




### 类型约束(Type Constraint)

约束（constraint）规定了一个类型实参（type argument）必须满足的条件要求。
如果某个类型满足了某个约束规定的所有条件要求，那么它就是这个约束修饰的类型形参的一个合法的类型实参。 

在Go泛型中，我们使用interface类型来定义约束。为此，Go接口类型的定义也进行了扩展，我们既可以声明接口的方法集合，也可以声明可用作类型实参的类型列表。
```go
[T any]             // 任意类型
[T int]             // 只能是 int
[T ~int]            // 是 int 或底层类型是 int 的类型。(type I int)
[T int | string]    // 只能是 int 或 string。(interface{ int | string})
[T io.Reader]       // 任何实现io.Reader 接口的类型
```

- [参考代码](./02_typeParam_n_typeArgument/main.go)


为了支持使用接口类型来定义Go泛型类型参数的类型约束，Go 1.18对接口定义语法进行了扩展。 在接口定义中既可以定义接口的方法集(Method Set)，也可以声明可以被用作泛型类型参数的类型实参的类型集(Type Set)




T 就是上面介绍过的类型形参(Type parameter)，在定义Slice类型的时候 T 代表的具体类型并不确定，类似一个占位符

* int|float32|float64 这部分被称为类型约束(Type constraint)，中间的 | 的意思是告诉编译器，类型形参 T 只可以接收 int 或 float32 或 float64 这三种类型的实参
* 中括号里的 T int|float32|float64 这一整串因为定义了所有的类型形参(在这个例子里只有一个类型形参T），所以我们称其为 类型形参列表(type parameter list)
* 这里新定义的类型名称叫 Slice[T]



### 类型具化（instantiation）与类型推导（type inference）

Go编译器会根据传入的实参的类型，进行类型实参（type argument）的自动推导。自动类型推导使得人们在编写调用泛型函数的代码时可以使用一种更为自然的风格。


### 从方法集(Method set)到类型集(Type set)

> 在Go1.18之前，Go官方对接口(interface) 的定义是：接口是一个方法集(method set)


> Go1.18开始, An interface type defines a type set (一个接口类型定义了一个类型集)


## 类型的交集

```go
type AllInt interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint32
}

type Uint interface {
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// 接口 A 代表的是 AllInt 与 Uint 的 交集，即 ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
type A interface { // 接口A代表的类型集是 AllInt 和 Uint 的交集
    AllInt
    Uint
}

//  接口 B 代表的则是 AllInt 和 ~int 的交集，即 ~int
type B interface { // 接口B代表的类型集是 AllInt 和 ~int 的交集
    AllInt
    ~int
}


type Bad interface {
  int
  float32
} // 类型 int 和 float32 没有相交的类型，所以接口 Bad 代表的类型集为空
```

### 案例
```go
// 使用前
type StringSlice []string
type Float32Slie []float32
type Float64Slice []float64

// 使用后
type Slice[T int|float32|float64 ] []T

// 这里传入了类型实参int，泛型类型Slice[T]被实例化为具体的类型 Slice[int]
var a Slice[int] = []int{1, 2, 3}
fmt.Printf("Type Name: %T",a)  //输出：Type Name: Slice[int]
```



## 接口两种类型


Go1.18开始将接口分为了两种类型

- 基本接口(Basic interface)
- 一般接口(General interface)

### 基本接口(Basic interface)
接口定义中如果只有方法的话，那么这种接口被称为基本接口(Basic interface)。这种接口就是Go1.18之前的接口，用法也基本和Go1.18之前保持一致。

```go
type MyError interface { // 接口中只有方法，所以是基本接口
    Error() string
}

// 用法和 Go1.18之前保持一致
var err MyError = fmt.Errorf("hello world")
```

### 一般接口(General interface)
如果接口内不光只有方法，还有类型的话，这种接口被称为 一般接口(General interface)

```go
type ReadWriter interface {
    ~string | ~[]rune

    Read(p []byte) (n int, err error)
    Write(p []byte) (n int, err error)
}
```

般接口类型不能用来定义变量，只能用于泛型的类型约束中。所以以下的用法是错误的：
```go
type Uint interface {
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

var uintInf Uint // 错误。Uint是一般接口，只能用于类型约束，不得用于变量定义

```

## 常见错误

1. 定义泛型类型的时候，基础类型不能只有类型形参
```go

type CommonType[T int | string | float32] T
```

2. 当类型约束的一些写法会被编译器误认为是表达式时会报错

```go
//✗ 错误。T *int会被编译器误认为是表达式 T乘以int，而不是int指针
type NewType[T *int] []T
// 上面代码再编译器眼中：它认为你要定义一个存放切片的数组，数组长度由 T 乘以 int 计算得到
type NewType [T * int][]T 

//✗ 错误。和上面一样，这里不光*被会认为是乘号，| 还会被认为是按位或操作
type NewType2[T *int|*float64] []T 

//✗ 错误
type NewType2 [T (int)] []T 
```
为了避免这种误解，解决办法就是给类型约束包上 interface{} 或加上逗号消除歧义（

```go
type NewType[T interface{*int}] []T
type NewType2[T interface{*int|*float64}] []T 

// 如果类型约束中只有一个类型，可以添加个逗号消除歧义
type NewType3[T *int,] []T


//✗ 错误。如果类型约束不止一个类型，加逗号是不行的
type NewType4[T *int|*float32,] []T 
```

3. ~ 限制
```go
type MyInt int

type _ interface {
    ~[]byte  // 正确
    ~MyInt   // 错误，~后的类型必须为基本类型
    ~error   // 错误，~后的类型不能为接口
}
```

4. 接口定义的种种限制规则

用 | 连接多个类型的时候，类型之间不能有相交的部分(即必须是不交集):
```go
type MyInt int

// 错误，MyInt的底层类型是int,和 ~int 有相交的部分
type _ interface {
    ~int | MyInt
}
```

但是相交的类型中是接口的话，则不受这一限制：
```go
type MyInt int

type _ interface {
    ~int | interface{ MyInt }  // 正确
}

type _ interface {
    interface{ ~int } | MyInt // 也正确
}

type _ interface {
    interface{ ~int } | interface{ MyInt }  // 也正确
}

```


类型的并集中不能有类型形参
```go
type MyInf[T ~int | ~string] interface {
    ~float32 | T  // 错误。T是类型形参
}

type MyInf2[T ~int | ~string] interface {
    T  // 错误
}
```

接口不能直接或间接地并入自己

```go
type Bad interface {
    Bad // 错误，接口不能直接并入自己
}

type Bad2 interface {
    Bad1
}
type Bad1 interface {
    Bad2 // 错误，接口Bad1通过Bad2间接并入了自己
}

type Bad3 interface {
    ~int | ~string | Bad3 // 错误，通过类型的并集并入了自己
}
```

接口的并集成员个数大于一的时候不能直接或间接并入 comparable 接口
```go
type OK interface {
    comparable // 正确。只有一个类型的时候可以使用 comparable
}

type Bad1 interface {
    []int | comparable // 错误，类型并集不能直接并入 comparable 接口
}

type CmpInf interface {
    comparable
}
type Bad2 interface {
    chan int | CmpInf  // 错误，类型并集通过 CmpInf 间接并入了comparable
}
type Bad3 interface {
    chan int | interface{comparable}  // 理所当然，这样也是不行的
}
```

带方法的接口(无论是基本接口还是一般接口)，都不能写入接口的并集中：
```go
type _ interface {
    ~int | ~string | error // 错误，error是带方法的接口(一般接口) 不能写入并集中
}

type DataProcessor[T any] interface {
    ~string | ~[]byte

    Process(data T) (newData T)
    Save(data T) error
}

// 错误，实例化之后的 DataProcessor[string] 是带方法的一般接口，不能写入类型并集
type _ interface {
    ~int | ~string | DataProcessor[string] 
}

type Bad[T any] interface {
    ~int | ~string | DataProcessor[T]  // 也不行
}
```

5. Method cannot have type parameters 不支持泛型方法
解决方式: 实现泛型函数来实现泛型方法，把方法的receiver当成第一个参数传递过去,参考 https://github.com/marwan-at-work/singleflight

## 推荐写法

阶段一
```go
// 一个可以容纳所有int,uint以及浮点类型的泛型切片
type Slice[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64] []T
```

阶段二
```go
type IntUintFloat interface {
    int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type Slice[T IntUintFloat] []T
```

阶段三
```go
type Int interface {
    int | int8 | int16 | int32 | int64
}

type Uint interface {
    uint | uint8 | uint16 | uint32
}

type Float interface {
    float32 | float64
}

type Slice[T Int | Uint | Float] []T  // 使用 '|' 将多个接口类型组合
```



## 参考资料
- https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md
- https://dominikbraun.io/blog/a-gentle-introduction-to-generics-in-go/
- [鸟窝关于 Go 泛形实现](https://colobu.com/2021/08/30/how-is-go-generic-implemented/)
- [Go泛型不支持泛型方法，这是一个悲伤的故事](https://colobu.com/2021/12/22/no-parameterized-methods/)
- [Go 泛型全面讲解](https://www.cnblogs.com/insipid/p/17772581.html)
