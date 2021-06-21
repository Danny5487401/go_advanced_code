package main

import "fmt"

// 方法二：用组合代替继承
/*
分析：
	本来属于基类的name和description不能放到基类中实现。其实只要转换一下思路，这个问题是很容易用组合解决的。
	如果我们认为Menu和MenuItem本质上是两个不同的事物，只是恰巧有（has-a）一些相同的属性，那么将相同的属性抽离出来，
	再分别组合进两者，问题就迎刃而解了
*/

// 抽离出来的属性
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

// 改写MenuItem
type MenuItem struct {
	MenuDesc
	price float32
}

func NewMenuItem(name, description string, price float32) MenuComponent {
	return &MenuItem{
		MenuDesc: MenuDesc{
			name:        name,
			description: description,
		},
		price: price,
	}
}

// 实际不需要的功能
func (m *MenuItem) Add(MenuComponent) {
	panic("not implement")
}

func (m *MenuItem) Remove(int) {
	panic("not implement")
}

func (m *MenuItem) Child(int) MenuComponent {
	panic("not implement")
}

// 改写 Menu方式一：
//type Menu struct {
//	MenuDesc
//	children []MenuComponent
//}
//
//func NewMenu(name, description string) MenuComponent {
//	return &Menu{
//		MenuDesc: MenuDesc{
//			name:        name,
//			description: description,
//		},
//	}
//}

// 改写 Menu方式二：
// 对于Menu，更好的做法是把children和Add() Remove() Child()也提取封装后再进行组合，这样Menu的功能一目了然
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

func NewMenu(name, description string) MenuComponent {
	return &Menu{
		MenuDesc: MenuDesc{
			name:        name,
			description: description,
		},
	}
}

// MenuItem的实现
func (m *MenuItem) Price() float32 {
	return m.price
}

func (m *MenuItem) Print() {
	fmt.Printf("  %s, ￥%.2f\n", m.name, m.price)
	fmt.Printf("    -- %s\n", m.description)
}

// Menu的实现
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

// 共有抽象方法
type MenuComponent interface {
	Price() float32
	Print()

	Add(MenuComponent)
	Remove(int)
	Child(int) MenuComponent
}

// 开始调用
func main() {
	menu1 := NewMenu("培根鸡腿燕麦堡套餐", "供应时间：09:15--22:44")
	menu1.Add(NewMenuItem("主食", "培根鸡腿燕麦堡1个", 11.5))
	menu1.Add(NewMenuItem("小吃", "玉米沙拉1份", 5.0))
	menu1.Add(NewMenuItem("饮料", "九珍果汁饮料1杯", 6.5))

	menu2 := NewMenu("奥尔良烤鸡腿饭套餐", "供应时间：09:15--22:44")
	menu2.Add(NewMenuItem("主食", "新奥尔良烤鸡腿饭1份", 15.0))
	menu2.Add(NewMenuItem("小吃", "新奥尔良烤翅2块", 11.0))
	menu2.Add(NewMenuItem("饮料", "芙蓉荟蔬汤1份", 4.5))

	all := NewMenu("超值午餐", "周一至周五有售")
	all.Add(menu1)
	all.Add(menu2)

	all.Print()
	// 添加子项 修改前
	if m, ok := all.Child(1).(*Menu); ok {
		m.Add(NewMenuItem("玩具", "Hello Kitty", 5.0))
	}
	all.Print()
	// 这里我们对类型的要求其实并没有那么强，并不需要它一定要是Menu，只是需要其提供组合MenuComponent的功能，所以可以提炼出这样一个接口
	if m, ok := all.Child(1).(Group); ok {
		m.Add(NewMenuItem("飞机", "Hello Danny2", 6.0))
	}
	all.Print()
}

type Group interface {
	Add(c MenuComponent)
	Remove(idx int)
	Child(idx int) MenuComponent
}
