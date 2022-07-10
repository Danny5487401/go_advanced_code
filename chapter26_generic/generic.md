# Generic泛型
Go 1.18版本增加了对泛型的支持

泛型在Go语言中增加了三个新的重要内容

- 函数和类型新增对 *类型形参(type parameters)* 的支持。 
- 将接口类型定义为类型集合，包括没有方法的接口类型。 
- 支持类型推导，大多数情况下，调用泛型函数时可省略类型实参(type arguments)。

## 使用背景
- Go官方团队在Go 1.18 Beta 1版本的标准库里因为泛型设计而引入了contraints包，Go官方团队的技术负责人Russ Cox在2022.01.25提议将constraints包从Go标准库里移除，放到x/exp项目下，
  该提议也同Go语言发明者Rob Pike, Robert Griesemer和Ian Lance Taylor做过讨论，得到了他们的同意。
- golang.org/x下所有package的源码独立于Go源码的主干分支，也不在Go的二进制安装包里。如果需要使用golang.org/x下的package，可以使用go get来安装。
- golang.org/x/exp下的所有package都属于实验性质或者被废弃的package，不建议使用。

## 类型形参(Type Parameters)

函数和类型被允许拥有类型形参(Type Parameters)。一个类型形参列表看起来和普通的函数形参列表一样，只是它使用的是中括号而不是小括号。


看一个用于浮点值的基本的、非泛型的Min函数：
```go
func Min(x, y float64) float64 {
    if x < y {
        return x
    }
    return y
}
```
通过添加一个类型形参列表来使这个函数泛型化，以使其适用于不同的类型。
```go
func GMin[T constraints.Ordered](x, y T) T {
    if x < y {
        return x
    }
    return y
}
```
现在我们可以像下面代码那样，用一个类型实参(Type argument)来调用这个函数了：
```go
x := GMin[int](2, 3)


```
向GMin函数提供类型实参，在本例中是int，称为实例化(instantiation)。实例化分两步进行。
1. 首先，编译器在整个泛型函数或泛型类型中把所有的类型形参替换成它们各自的类型实参。
2. 第二，编译器验证每个类型实参是否满足各自的约束条件。我们很快就会知道这意味着什么，但是如果第二步失败，实例化就会失败，程序也会无效。

```go
fmin := GMin[float64]
m := fmin(2.71, 3.14)
```
GMin[float64]的实例化产生了一个与Min函数等效的函数，我们可以在函数调用中使用它。


## 类型集合(Type sets)
一个普通函数的每个值形参(译注：value parameter，相对于类型形参type parameter)都有一个对应的类型；该类型定义了一组值.

类型形参列表中的每个类型形参都有一个类型。因为类型形参本身就是一个类型，所以类型形参的类型定义了类型的集合。
这种元类型(meta-type)被称为类型约束(type constraint).
