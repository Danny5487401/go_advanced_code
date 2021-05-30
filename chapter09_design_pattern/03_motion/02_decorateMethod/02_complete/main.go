package main

import "fmt"

/*
装饰模式：一种动态地往一个类中添加新的行为的设计模式.
从process2中得知主要包含四个角色，抽象构件 Component，具体构件 ConcreteComponent，抽象装饰类 Decorator，具体装饰类 ConcreteComponent
优点：
1. 可以通过一种动态的方式来扩展一个对象的功能
2. 可以使用多个具体装饰类来装饰同一对象，增加其功能
3. 具体组件类与具体装饰类可以独立变化，符合“开闭原则”

缺点：
1. 对于多次装饰的对象，易于出错，排错也很困难
2. 对于产生很多具体装饰类 ，增加系统的复杂度以及理解成本

使用场景：
1.需要给一个对象增加功能，这些功能可以动态地撤销，例如：在不影响其他对象的情况下，动态、透明的方式给单个对象添加职责，处理那些可以撤销的职责
2.需要给一批兄弟类增加或者改装功能
 */

/*
从example.png得知，
1. Component(抽象构建)：具体构建和抽象装饰类的基类，声明了在具体构建中实现的业务方法，UML类图中的Component
2. ConcreteComponent(具体构建)：抽象构建的子类，用于定义具体的构建对象，实现了在抽象构建中声明的方法，装饰器可以给它增加额外的职责(方法)，UML类图中的Window、TextBox、ListBox
3. Decorator(抽象装饰类)：也是抽象构建类的子类，用于给具体构建增加职责，但是具体职责在其子类中实现，UML类图中的ComponentDecorator
4. ConcreteDecorator(具体装饰类)：抽象装饰类的子类，负责向构建添加新的职责，UML类图中的ScrollBarDecorator、BlackBorderDecorator

 */

// 步骤一 ：实现Component抽象类以及ConcreteComponent具体构建

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
//  步骤二：实现ConcretDecorator具体装饰类
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
// 步骤三：  定义工厂函数生产出具体装饰类
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




