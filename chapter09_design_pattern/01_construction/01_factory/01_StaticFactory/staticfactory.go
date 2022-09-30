package main

import "fmt"

// 看example.png实例流程，FactoryPatternDemo，我们的演示类使用 ShapeFactory 来获取 Shape 对象。它将向 ShapeFactory 传递信息（CIRCLE / RECTANGLE / SQUARE），以便获取它所需对象的类型

// Shape 步骤一： 创建一个接口
type Shape interface {
	Draw()
}

// Rectangle 步骤二：创建实体接口的实体类
type Rectangle struct {
}

func (r Rectangle) Draw() {
	fmt.Println("Inside Rectangle::draw() method.")
}

type Square struct {
}

func (s Square) Draw() {
	fmt.Println("Inside Square ::draw() method.")
}

type Circle struct {
}

func (c Circle) Draw() {
	fmt.Println("Inside Circle  ::draw() method.")
}

// ShapeFactory 步骤三：生成基于给定信息的实体类的对象
type ShapeFactory struct {
}

// 使用 getShape 方法获取形状类型的对象
func (s ShapeFactory) getShape(shapeType string) Shape {

	switch shapeType {
	case "CIRCLE":
		return Circle{}
	case "RECTANGLE":
		return Rectangle{}
	case "SQUARE":

		return Square{}
	default:
		return nil
	}

}
