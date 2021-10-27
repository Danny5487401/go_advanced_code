package main

import "fmt"

// Component 步骤一 ：实现Component抽象类以及ConcreteComponent具体构建
type Component interface {
	Display()
}

type Window struct{}

func (w Window) Display() {
	fmt.Println("显示窗体")
}

type TextBox struct{}

func (t TextBox) Display() {
	fmt.Println("显示文本框")
}

type ListBox struct{}

func (l ListBox) Display() {
	fmt.Println("显示列表框")
}

// ScrollBarDecorator 步骤二：实现ConcreteDecorator具体装饰类
type ScrollBarDecorator struct {
	Component
}

func (sbd ScrollBarDecorator) Display() {
	fmt.Println("为构建增加滚动条")
	sbd.Component.Display()
}

type BlackBorderDecorator struct {
	Component
}

func (bbd BlackBorderDecorator) Display() {
	fmt.Println("为构建增加黑色边框")
	bbd.Component.Display()
}

// NewDecorator 步骤三：  定义工厂函数生产出具体装饰类
func NewDecorator(t string, decorator Component) Component {
	switch t {
	case "sbd":
		return ScrollBarDecorator{
			Component: decorator,
		}
	case "bbd":
		return BlackBorderDecorator{
			Component: decorator,
		}
	default:
		return nil
	}
}

// 测试
func main() {
	// 1. 选择Window
	component := Window{}
	tScrollBarDecorator := NewDecorator("sbd", component)
	tScrollBarDecorator.Display()
	fmt.Println("==============================")
	//tBlackBorderDecorator := NewDecorator("bbd", component)
	tBlackBorderDecorator := NewDecorator("bbd", tScrollBarDecorator)
	tBlackBorderDecorator.Display()

	// 2. 选择TextBox
	//component := TextBox{}
	//tScrollBarDecorator := NewDecorator("sbd", component)
	//tScrollBarDecorator.Display()
	//fmt.Println("==============================")
	//tBlackBorderDecorator := NewDecorator("bbd", tScrollBarDecorator)
	//tBlackBorderDecorator.Display()
}
