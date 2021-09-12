#访问者模式

定义：

    表示一个作用于某对象结构中的各元素的操作。它使你可以在不改变各元素的类的前提下定义作用于这些元素的新操作
示意图

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

做法：

	对象只要预留访问者接口Accept则后期为对象添加功能的时候就不需要改动对象
大概的流程就是:

	1。从结构容器中取出元素
	2。创建一个访问者
	3。将访问者载入传入的元素（即让访问者访问元素）
	4。获取输出

##角色组成：

	1.抽象访问者:抽象类或者接口，声明访问者可以访问哪些元素，具体到程序中就是visit方法中的参数定义哪些对象是可以被访问的
	2.访问者:实现抽象访问者所声明的方法，它影响到访问者访问到一个类后该干什么，要做什么事情
	3.抽象元素类:接口或者抽象类，声明接受哪一类访问者访问，程序上是通过accept方法中的参数来定义的。抽象元素一般有两类方法，一部分是本身的业务逻辑，另外就是允许接收哪类访问者来访问。
	4.元素类:实现抽象元素类所声明的accept方法，通常都是visitor.visit(this)，基本上已经形成一种定式了
	5.结构容器: (非必须) 保存元素列表，可以放置访问者.一个元素的容器，一般包含一个容纳多个不同类、不同接口的容器，如List、Set、Map等，在项目中一般很少抽象出这个角色

##源码参考：k8s
```go
//k8s.io/cli-runtime/pkg/resource/interfaces.go

// Visitor 即为访问者这个对象
type Visitor interface {
    Visit(VisitorFunc) error
}
// VisitorFunc对应这个对象的方法，也就是定义中的“操作”
type VisitorFunc func(*Info, error) error
```

###Visitor的链式处理

	1.多个对象聚合为一个对象
		VisitorList
		EagerVisitorList
	2.多个方法聚合为一个方法
		DecoratedVisitor
		ContinueOnErrorVisitor
	3.将对象抽象为多个底层对象，逐个调用方法
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