# 访问者模式

将算法与操作对象的结构分离的一种方法。这种分离的实际结果是能够在不修改结构的情况下向现有对象结构添加新操作

## 需求
假设你是一个代码库的维护者， 代码库中包含不同的形状结构体， 如：
- 方形
- 圆形
- 三角形
  
上述每个形状结构体都实现了通用形状接口。

在公司员工开始使用你维护的代码库时， 你就会被各种功能请求给淹没。 
让我们来看看其中比较简单的请求： 有个团队请求你在形状结构体中添加 get Area获取面积行为。

### 解决方式
1. 将 get Area方法直接添加至形状接口， 然后在各个形状结构体中进行实现。 这似乎是比较好的解决方案， 但其代价也比较高。 
   作为代码库的管理员， 相信你也不想在每次有人要求添加另外一种行为时就去冒着风险改动自己的宝贝代码。
   不过， 你也一定想让其他团队的人还是用一用自己的代码库。
   
2. 请求功能的团队自行实现行为。 然而这并不总是可行， 因为行为可能会依赖于私有代码。

3. 使用访问者模式来解决上述问题。 首先定义一个如下访问者接口：
```go
type visitor interface {
    visitForSquare(square)
    visitForCircle(circle)
    visitForTriangle(triangle)
}
```

我们可以使用 visit For Square(square) 、  visit For Circle (circle)以及 visit ForTriangle(triangle)函数来为方形、 圆形以及三角形添加相应的功能。

你可能在想， 为什么我们不再访问者接口里面使用单一的 visit(shape)方法呢？ 这是因为 Go 语言不支持方法重载， 所以你无法以相同名称、 不同参数的方式来使用方法，但是也能采用类型断言的方式来对不同的类型做不同的操作。

第二项重要的工作是将 accept接受方法添加至形状接口中。
```go
func accept(v visitor)
```
所有形状结构体都需要定义此方法， 类似于：
```go
func (obj *square) accept(v visitor){
    v.visitForSquare(obj)
}
```

等等， 我刚才是不是提到过， 我们并不想修改现有的形状结构体？ 很不幸， 在使用访问者模式时， 我们必须要修改形状结构体。 但这样的修改只需要进行一次。

如果添加任何其他行为， 比如 get Num Sides获取边数和 get Middle Coordinates获取中点坐标 ， 我们将使用相同的 accept(v visitor)函数， 而无需对形状结构体进行进一步的修改。

最后， 形状结构体只需要修改一次， 并且所有未来针对不同行为的请求都可以使用相同的 accept 函数来进行处理。 如果团队成员请求 get Area行为， 我们只需简单地定义访问者接口的具体实现， 并在其中编写面积的计算逻辑即可

## [代码参考1：案例](https://github.com/Danny5487401/github.com/Danny5487401/go_advanced_code/tree/main/chapter09_design_pattern/03_motion/01_visitor/example1/main.go)
## [代码参考2：根据类型断言实现访问者](https://github.com/Danny5487401/github.com/Danny5487401/go_advanced_code/tree/main/chapter09_design_pattern/03_motion/01_visitor/example2)

## 示意图

```css
   ┌─────────┐       ┌───────────────────────┐
   │ Client  │─ ─ ─ >│        Visitor        │
   └─────────┘       ├───────────────────────┤
        │            │visitElementA(ElementA)│
                     │visitElementB(ElementB)│
        │            └───────────────────────┘
                                 ▲
        │                ┌───────┴───────┐
                         │               │
        │         ┌─────────────┐ ┌─────────────┐
                  │  VisitorA   │ │  VisitorB   │
        │         └─────────────┘ └─────────────┘
        ▼
┌───────────────┐        ┌───────────────┐
│ObjectStructure│─ ─ ─ ─>│    Element    │
├───────────────┤        ├───────────────┤
│handle(Visitor)│        │accept(Visitor)│
└───────────────┘        └───────────────┘
                                 ▲
                        ┌────────┴────────┐
                        │                 │
                ┌───────────────┐ ┌───────────────┐
                │   ElementA    │ │   ElementB    │
                ├───────────────┤ ├───────────────┤
                │accept(Visitor)│ │accept(Visitor)│
                │doA()          │ │doB()          │
                └───────────────┘ └───────────────┘
```

做法 : 对象只要预留访问者接口Accept则后期为对象添加功能的时候就不需要改动对象

## 大概的流程就是
1. 从结构容器中取出元素
2. 创建一个访问者
3. 将访问者载入传入的元素（即让访问者访问元素）
4. 获取输出

## 角色组成：

1. 抽象访问者:抽象类或者接口，声明访问者可以访问哪些元素，具体到程序中就是visit方法中的参数定义哪些对象是可以被访问的
2. 访问者:实现抽象访问者所声明的方法，它影响到访问者访问到一个类后该干什么，要做什么事情
3. 抽象元素类:接口或者抽象类，声明接受哪一类访问者访问，程序上是通过accept方法中的参数来定义的。抽象元素一般有两类方法，一部分是本身的业务逻辑，另外就是允许接收哪类访问者来访问。
4. 元素类:实现抽象元素类所声明的accept方法，通常都是visitor.visit(this)，基本上已经形成一种定式了
5. 结构容器: (非必须) 保存元素列表，可以放置访问者.一个元素的容器，一般包含一个容纳多个不同类、不同接口的容器，如List、Set、Map等，在项目中一般很少抽象出这个角色



## 源码参考：k8s
kubectl 的代码比较复杂，不过，简单来说，基本原理就是它从命令行和 YAML 文件中获取信息， 通过 Builder 模式并把其转成一系列的资源，最后用 Visitor 模式来迭代处理这些 Resource
### Visitor 模式定义
```go
//k8s.io/cli-runtime/pkg/resource/interfaces.go

// Visitor 即为访问者这个对象
type Visitor interface {
    Visit(VisitorFunc) error
}
// VisitorFunc对应这个对象的方法，也就是定义中的“操作”
type VisitorFunc func(*Info, error) error

type Info struct {
    Namespace   string
    Name        string
    OtherThings string
}

func (info *Info) Visit(fn VisitorFunc) error { 
	return fn(info, nil)
}

```
### 1. 具体 Name Visitor
这个 Visitor 主要是用来访问 Info 结构中的 Name 和 NameSpace 成员：
```go

type NameVisitor struct {
    visitor Visitor
}

func (v NameVisitor) Visit(fn VisitorFunc) error {
    return v.visitor.Visit(func(info *Info, err error) error {
        fmt.Println("NameVisitor() before call function")
        err = fn(info, err)
        if err == nil {
            fmt.Printf("==> Name=%s, NameSpace=%s\n", info.Name, info.Namespace)
        }
        fmt.Println("NameVisitor() after call function")
        return err
    })
}
// 在实现 Visit() 方法时，调用了自己结构体内的那个 Visitor的 Visitor() 方法，这其实是一种修饰器的模式，用另一个 Visitor 修饰了自己
```

### 2. OtherVisitor
这个 Visitor 主要用来访问 Info 结构中的 OtherThings 成员：
```go

type OtherThingsVisitor struct {
  visitor Visitor
}

func (v OtherThingsVisitor) Visit(fn VisitorFunc) error {
  return v.visitor.Visit(func(info *Info, err error) error {
    fmt.Println("OtherThingsVisitor() before call function")
    err = fn(info, err)
    if err == nil {
      fmt.Printf("==> OtherThings=%s\n", info.OtherThings)
    }
    fmt.Println("OtherThingsVisitor() after call function")
    return err
  })
}
```
### 3. LogVisitor
```go

type LogVisitor struct {
  visitor Visitor
}

func (v LogVisitor) Visit(fn VisitorFunc) error {
  return v.visitor.Visit(func(info *Info, err error) error {
    fmt.Println("LogVisitor() before call function")
    err = fn(info, err)
    fmt.Println("LogVisitor() after call function")
    return err
  })
}
```

使用时
```go
func main() {
  info := Info{}
  var v Visitor = &info
  v = LogVisitor{v}
  v = NameVisitor{v}
  v = OtherThingsVisitor{v}

  loadFile := func(info *Info, err error) error {
    info.Name = "Hao Chen"
    info.Namespace = "MegaEase"
    info.OtherThings = "We are running as remote team."
    return nil
  }
  v.Visit(loadFile)
}
```
打印结果
```shell

LogVisitor() before call function
NameVisitor() before call function
OtherThingsVisitor() before call function
==> OtherThings=We are running as remote team.
OtherThingsVisitor() after call function
==> Name=Hao Chen, NameSpace=MegaEase
NameVisitor() after call function
LogVisitor() after call function
```


Visitor 修饰器
```go

type DecoratedVisitor struct {
  visitor    Visitor
  decorators []VisitorFunc
}

func NewDecoratedVisitor(v Visitor, fn ...VisitorFunc) Visitor {
  if len(fn) == 0 {
    return v
  }
  return DecoratedVisitor{v, fn}
}

// Visit implements Visitor
func (v DecoratedVisitor) Visit(fn VisitorFunc) error {
  return v.visitor.Visit(func(info *Info, err error) error {
    if err != nil {
      return err
    }
    if err := fn(info, nil); err != nil {
      return err
    }
    for i := range v.decorators {
      if err := v.decorators[i](info, nil); err != nil {
        return err
      }
    }
    return nil
  })
}
```

上面方法简化使用
```go

info := Info{}
var v Visitor = &info
v = NewDecoratedVisitor(v, NameVisitor, OtherVisitor)

v.Visit(LoadFile)
```

### Visitor的链式处理

	1. 多个对象聚合为一个对象
		VisitorList
		EagerVisitorList
	2. 多个方法聚合为一个方法
		DecoratedVisitor
		ContinueOnErrorVisitor
	3. 将对象抽象为多个底层对象，逐个调用方法
		FlattenListVisitor
		FilteredVisitor
	
1. VisitorList:封装多个Visitor为一个，出现错误就立刻中止并返回  

```go
//-->k8s.io/cli-runtime/pkg/resource/visitor.go
	// VisitorList定义为[]Visitor，又实现了Visit方法，也就是将多个[]Visitor封装为一个Visitor
	type VisitorList []Visitor

	// 发生error就立刻返回，不继续遍历
	func (l VisitorList) Visit(fn VisitorFunc) error {
		for i := range l {
			if err := l[i].Visit(fn); err != nil {
				return err
			}
		}
		return nil
	}

```