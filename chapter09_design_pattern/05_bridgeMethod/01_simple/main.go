/*
  桥接模式：将事务的多个纬度都抽象出来以解耦抽象与实际之间的绑定关系，使抽象和实际向着不同纬度改变；

如 process.png 所示：
1. Abstraction（抽象类）：用于定义抽象类的接口，它一般是抽象类而不是接口，其中定义了一个Implementor（实现类接口）类型的对象并可以维护该对象，
	它与Implementor之间具有关联关系，它既可以包含抽象业务方法，也可以包含具体业务方法。

2. RefinedAbstraction（扩充抽象类）：扩充由Abstraction定义的接口，通常情况下它不再是抽象类而是具体类，它实现了在Abstraction中声明的抽象业务方法，
	在RefinedAbstraction中可以调用在Implementor中定义的业务方法。

3. Implementor（实现类接口）：定义实现类的接口，这个接口不一定要与Abstraction的接口完全一致，事实上这两个接口可以完全不同，
	一般而言，Implementor接口仅提供基本操作，而Abstraction定义的接口可能会做更多更复杂的操作。Implementor接口对这些基本操作进行了声明，而具体实现交给其子类。
	通过关联关系，在Abstraction中不仅拥有自己的方法，还可以调用到Implementor中定义的方法，使用关联关系来替代继承关系。

4. ConcreteImplementor（具体实现类）：具体实现Implementor接口，在不同的ConcreteImplementor中提供基本操作的不同实现，
	在程序运行时，ConcreteImplementor对象将替换其父类对象，提供给抽象类具体的业务操作方法
意图：将抽象部分与实现部分分离，使他们都可以独立的变化。

主要解决：在有多种可能会变化的情况下，用继承会造成类爆炸问题，扩展起来不灵活。

何时使用：实现系统可能有多个角度分类，每一种角度都可能变化。

如何解决：把这种多角度分类分离出来，让他们独立变化，减少他们之间的耦合。

关键代码：抽象类依赖实现类。

优点：抽象和实现的分离、优秀的扩展能力、实现细节与客户透明

缺点：桥接模式的引用会增加系统的理解与设计难度，由于聚合关联关系建立在抽象层，要求开发者针对抽象进行设计与编程。

使用场景：对于两个独立变化的纬度，使用桥接模式在适合不过了。

 */
/*做法： 见process2.png
在使用桥接模式时，首先应该识别出一个类所具有的两个独立变化的维度，将它们设计为两个独立的继承等级结构，为两个维度都提供抽象层，并建立抽象耦合。
	通常情况下，将具有两个独立变化维度的类的一些普通业务方法和与之关系最密切的维度设计为“抽象类”层次结构（抽象部分），而将另一个维度设计为“实现类”层次结构（实现部分）

 */


package main

import "fmt"

// 发送信息的具体实现（操作）
type MessageImplementer interface {
	send(test, to string)
}
// 1.发送SMS
type MessageSMS struct{}

func (*MessageSMS) send(test, to string) {
	fmt.Printf("SMS信息：[%v]；发送到：[%v]\n", test, to)
}
func ViaSMS() *MessageSMS { //创建
	return &MessageSMS{}
}
// 2.发送email
type MessageEmail struct{}

func (*MessageEmail) send(test, to string) {
	fmt.Printf("Email信息：[%v]；发送到：[%v]\n", test, to)
}
func ViaEmail()  *MessageEmail{//创建
	return &MessageEmail{}
}


//发送信息的二次封装（发送普通信息、紧急信息）
type AbstractMessage interface {
	SendMessage(text, to string)
}

// 1。普通信息
type CommonMessage struct {
	method MessageImplementer
}

func (c *CommonMessage) SendMessage(text, to string) {
	c.method.send(text, to)
}
func NewCommonMessage(mi MessageImplementer) *CommonMessage {
	return &CommonMessage{method: mi}
}

// 2.紧急信息
type UrgencyMessage struct {
	method MessageImplementer
}

func (u *UrgencyMessage) SendMessage(text, to string) {
	u.method.send(fmt.Sprintf("【紧急信息】%v",text), to)
}
func NewUrgencyMessage(mi MessageImplementer) *UrgencyMessage {
	return &UrgencyMessage{method: mi}
}

// 主体逻辑
func main() {
	var m AbstractMessage
	//返回值为：*CommonMessage类型，而CommonMessage实现了AbstractMessage，因此需要定义的类型为；接口AbstractMessage

	m=NewCommonMessage(ViaSMS())
	m.SendMessage("你需要喝一杯咖啡吗","美女")

	m=NewCommonMessage(ViaEmail())
	m.SendMessage("好滴","帅哥")

	m=NewUrgencyMessage(ViaSMS())
	m.SendMessage("晚上桥头见","美女")

	m=NewUrgencyMessage(ViaEmail())
	m.SendMessage("不见不散","帅哥")
}

