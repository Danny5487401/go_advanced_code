package main

import "fmt"

/*
需求：
	把KFC里的食物认为是菜单项，一份套餐是菜单。菜单和菜单项有一些公有属性：名字、描述、价格、都能被购买等
实现
	面向对象实现
*/

// MenuComponent 第一个版本,提取所有的抽象功能
type MenuComponent interface {
	// 菜单具体项目的方法
	Name() string
	Description() string
	Price() float32
	Print()

	// 菜单额外需要的方法
	Add(MenuComponent)
	Remove(int)
	Child(int) MenuComponent
}

// MenuItem 菜单项的实现
type MenuItem struct {
	name        string
	description string
	price       float32
}

func NewMenuItem(name, description string, price float32) MenuComponent {
	return &MenuItem{
		name:        name,
		description: description,
		price:       price,
	}
}

// Name 实现接口
func (m *MenuItem) Name() string {
	return m.name
}

func (m *MenuItem) Description() string {
	return m.description
}

func (m *MenuItem) Price() float32 {
	return m.price
}

func (m *MenuItem) Print() {
	fmt.Printf("  %s, ￥%.2f\n", m.name, m.price)
	fmt.Printf("    -- %s\n", m.description)
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

/*

有两点请留意一下。

	1。NewMenuItem()创建的是MenuItem，但返回的是抽象的接口MenuComponent。（面向对象中的多态）
	2。因为MenuItem是叶节点，无法提供Add() Remove() Child()这三个方法的实现，所以若被调用会panic
*/

// 菜单的实现
type Menu struct {
	name        string
	description string
	children    []MenuComponent
}

func NewMenu(name, description string) MenuComponent {
	return &Menu{
		name:        name,
		description: description,
	}
}

func (m *Menu) Name() string {
	return m.name
}

func (m *Menu) Description() string {
	return m.description
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

func (m *Menu) Add(c MenuComponent) {
	m.children = append(m.children, c)
}

func (m *Menu) Remove(idx int) {
	m.children = append(m.children[:idx], m.children[idx+1:]...)
}

func (m *Menu) Child(idx int) MenuComponent {
	return m.children[idx]
}

/*
问题：
	1，代码重复。MenuItem和Menu中都有name、description这两个属性和方法，重复写两遍明显冗余。如果使用其它任何面向对象语言，这两个属性和方法都应该移到基类中实现。可是Go没有继承
	func (m *Menu) Name() string {
		return m.name
	}

	func (m *Menu) Description() string {
		return m.description
	}
	2.添加产品menuComponent
		使用者拿到一个MenuComponent后，依然要知道其类型后才能正确使用，假如不加判断在MenuItem使用Add()等未实现的方法就会产生panic
*/
