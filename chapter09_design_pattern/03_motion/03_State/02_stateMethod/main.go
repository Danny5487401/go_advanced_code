package _2_stateMethod

import "fmt"

/*
需求：
	当今社会，论坛贴吧很多，我们也会加入感兴趣的论坛，偶尔进行发言，但有时却会发现不能发帖了，原来是昨天的某个帖子引发了口水战，被举报了。
	假设有三种状态，normal(正常），restricted(受限)，closed(封号)，判断依据是一个健康值（这里只是假设）。不同的状态有不同的权限
实现方式二：
	使用状态模式
*/

type Account struct {
	State       ActionState
	HealthValue int
}

// 定义接口拥有不同的功能
type ActionState interface {
	View()
	Comment()
	Post()
}

// 初始化
func NewAccount(health int) *Account {
	a := &Account{
		HealthValue: health,
	}
	a.changeState()
	return a
}

// 改变状态
func (a *Account) changeState() {
	if a.HealthValue <= -10 {
		a.State = &CloseState{}
	} else if a.HealthValue > -10 && a.HealthValue <= 0 {
		a.State = &RestrictedState{}
	} else if a.HealthValue > 0 {
		a.State = &NormalState{}
	}
}

// 定义状态--正常-限制-封号

type NormalState struct {
}

func (n *NormalState) View() {
	fmt.Println("正常看帖")
}

func (n *NormalState) Comment() {
	fmt.Println("正常评论")
}
func (n *NormalState) Post() {
	fmt.Println("正常发帖")
}

type RestrictedState struct {
}

func (r *RestrictedState) View() {
	fmt.Println("正常看帖")
}

func (r *RestrictedState) Comment() {
	fmt.Println("正常评论")
}
func (r *RestrictedState) Post() {
	fmt.Println("抱歉，你的健康值小于0，不能发帖")
}

type CloseState struct {
}

func (c *CloseState) View() {
	fmt.Println("账号被封，无法看帖")
}

func (c *CloseState) Comment() {
	fmt.Println("抱歉，你的健康值小于-10，不能评论")
}
func (c *CloseState) Post() {
	fmt.Println("抱歉，你的健康值小于0，不能发帖")
}

//给账户设定健康值
func (a *Account) SetHealth(value int) {
	a.HealthValue = value
	a.changeState()
}

/*
优点
	1。封装了状态的转换规则，在状态模式中可以将状态的转换代码封装在环境类或者具体状态类中，可以对状态转换代码进行集中管理，而不是分散在一个个业务方法中。
	2。将所有与某个状态有关的行为放到一个类中，只需要注入一个不同的状态对象即可使环境对象拥有不同的行为。
	3。允许状态转换逻辑与状态对象合成一体，而不是提供一个巨大的条件语句块，状态模式可以避免使用庞大的条件语句来将业务方法和状态转换代码交织在一起。
	4。可以让多个环境对象共享一个状态对象，从而减少系统中对象的个数。
缺点
	1。状态模式的使用必然会增加系统中类和对象的个数，导致系统运行开销增大。
	2。状态模式的结构与实现都较为复杂，如果使用不当将导致程序结构和代码的混乱，增加系统设计的难度。
	3。状态模式对“开闭原则”的支持并不太好，增加新的状态类需要修改那些负责状态转换的源代码，否则无法转换到新增状态；
		而且修改某个状态类的行为也需修改对应类的源代码。
适用环境
	对象的行为依赖于它的状态（属性）并且可以根据它的状态改变而改变它的相关行为。
	代码中包含大量与对象状态有关的条件语句，这些条件语句的出现，会导致代码的可维护性和灵活性变差，不能方便地增加和删除状态，使客户类与类库之间的耦合增强
模式应用：
	状态模式在工作流或游戏等类型的软件中得以广泛使用，甚至可以用于这些系统的核心功能设计，
	如在政府OA办公系统中，一个批文的状态有多种：尚未办理；正在办理；正在批示；正在审核；已经完成等各种状态，而且批文状态不同时对批文的操作也有所差异。
	使用状态模式可以描述工作流对象（如批文）的状态转换以及不同状态下它所具有的行为
*/
