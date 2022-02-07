package main

import "fmt"

// 方法二：用组合代替继承
/*
组合模式：
	将对象组合成树形结构以表示“部分整体”的层次结构，组合模式使得用户对单个对象和组合对象的使用具有一致性
实现分析：
	本来属于基类的name和description不能放到基类中实现。其实只要转换一下思路，这个问题是很容易用组合解决的。
	如果我们认为Menu和MenuItem本质上是两个不同的事物，只是恰巧有（has-a）一些相同的属性，那么将相同的属性抽离出来，
	再分别组合进两者，问题就迎刃而解了
实现：
	Go语言中善用组合有助于表达数据结构的意图。特别是当一个比较复杂的对象同时处理几方面的事情时，将对象拆成独立的几个部分再组合到一起，
	会非常清晰优雅。例如上面的MenuItem就是描述+价格，Menu就是描述+子菜单

*/

type MenuComponent interface {
	Price() float32
	Print()
}

type MenuDesc struct {
	name        string
	description string
}

func (m *MenuDesc) Name() string {
	return m.name
}

func (m *MenuDesc) Description() string {
	return m.description
}

type MenuItem struct {
	MenuDesc
	price float32
}

func NewMenuItem(name, description string, price float32) *MenuItem {
	return &MenuItem{
		MenuDesc: MenuDesc{
			name:        name,
			description: description,
		},
		price: price,
	}
}
func (m *MenuItem) Price() float32 {
	return m.price
}

func (m *MenuItem) Print() {
	fmt.Printf("  %s, ￥%.2f\n", m.name, m.price)
	fmt.Printf("    -- %s\n", m.description)
}

type MenuGroup struct {
	children []MenuComponent
}

func (m *Menu) Add(c MenuComponent) {
	m.children = append(m.children, c)
}

func (m *Menu) Remove(idx int) {
	m.children = append(m.children[:idx], m.children[idx+1:]...)
}

func (m *Menu) Child(idx int) MenuComponent {
	return m.children[idx]
}

type Menu struct {
	MenuDesc
	MenuGroup
}

func NewMenu(name, description string) *Menu {
	return &Menu{
		MenuDesc: MenuDesc{
			name:        name,
			description: description,
		},
	}
}

func (m *Menu) Price() (price float32) {
	for _, v := range m.children {
		price += v.Price()
	}
	return
}

func (m *Menu) Print() {
	fmt.Printf("%s, %s, ￥%.2f\n", m.name, m.description, m.Price())
	fmt.Println("------------------------")
	for _, v := range m.children {
		v.Print()
	}
	fmt.Println()
}

type Group interface {
	Add(c MenuComponent)
	Remove(idx int)
	Child(idx int) MenuComponent
}

/*
比较与思考
前后两份代码差异其实很小：
	1. 第二份实现的接口简单一些，只有两个函数。
	2. New函数返回值的类型不一样。

从思路上看，差异很大却也有些微妙：
	1. 第一份实现中接口是模板，是struct的蓝图，其属性来源于事先对系统组件的综合分析归纳；第二份实现中接口是一份约束声明，其属性来源于使用者对被使用者的要求。
	2. 第一份实现认为children中的MenuComponent是一种具体对象，这个对象具有一系列方法可以调用，只是其方法的功能会由于子类覆盖而表现不同；第二份实现则认为children中的MenuComponent可以是任意无关的对象，唯一的要求是他们“恰巧”实现了接口所指定的约束条件。
	注意第一份实现中，MenuComponent中有Add()、Remove()、Child()三个方法，但却不一定是可用的，能不能使用由具体对象的类型决定；第二份实现中则不存在这些不安全的方法，因为New函数返回的是具体类型，所以可以调用的方法都是安全的

*/
