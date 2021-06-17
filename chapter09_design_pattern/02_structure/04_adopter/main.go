package _4_adopter

// 适配器模式
/*
4类参与者介绍：
	Target
		定义Client使用的与特定领域相关的接口。
	Client
		与符合Target接口的对象协同。
	Adaptee
		定义一个已经存在的接口，这个接口需要适配。
	Adapter
		对Adaptee的接口与Target接口进行适配。
需求：
	现在只有Type-C接口，当前有一个键盘，但它使用的是USB接口，要求实现一个适配器转接头，使得键盘也可以插上Type-C接口成功使用

应用场景：
	新旧版本
优点：
	可以让任何两个没有关联的类一起运行、提高了类的复用、增加了类的透明度、灵活性好

缺点：
	过多地使用适配器，会让系统非常凌，不易整体进行把握。比如，明明看到调用的是A接口，其实内部被适配成了B接口的实现，一个系统如果太多出现这种情况，无异于异常灾难。因为如果不是很有必要，可以不使用适配器，而是直接对系统进重构。

*/

// TypeC接口:需要被适配的接口
type TypeC interface {
	UseTypeC() string
}

// USB接口:
type USB interface {
	UseUSB() string
}

func NewUSB() USB {
	return &keyboard{}
}

type keyboard struct {
}

// keyboard实现了USB接口
func (k *keyboard) UseUSB() string {
	return "I use USB interface"
}

// 适配器
type adapter struct {
	USB // 组合内嵌
}

//UseTypeC实现了Type-C接口
func (a *adapter) UseTypeC() string {
	return a.UseUSB() + ", but now I use Type-C interface"
}

//NewAdapter 是适配器的工厂函数
func NewAdapter(keyboard USB) TypeC {
	return &adapter{
		USB: keyboard,
	}
}
