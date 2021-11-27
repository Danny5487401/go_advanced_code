package main

import "fmt"



// 看example.png实例流程，FactoryPatternDemo，我们的演示类使用 ShapeFactory 来获取 Shape 对象。它将向 ShapeFactory 传递信息（CIRCLE / RECTANGLE / SQUARE），以便获取它所需对象的类型

// 步骤一： 创建一个接口
type Shape interface {
	Draw()
}

// 步骤二：创建实体接口的实体类
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

// 步骤三：生成基于给定信息的实体类的对象

type ShapeFactory struct {
}

//使用 getShape 方法获取形状类型的对象
func (s ShapeFactory) getShape(shapeType string) Shape {

	if shapeType == "" {
		return nil
	}
	if shapeType == "CIRCLE" {
		return Circle{}
	} else if shapeType == "RECTANGLE" {
		return Rectangle{}
	} else if shapeType == "SQUARE" {
		return Square{}
	}
	return nil
}

// 步骤四：使用该工厂，通过传递类型信息来获取实体类的对象。
func main() {
	factory := ShapeFactory{}
	factory.getShape("CIRCLE").Draw()
	factory.getShape("RECTANGLE").Draw()
	factory.getShape("SQUARE").Draw()
}

// 结果：
/*
Inside Circle  ::draw() method.
Inside Rectangle::draw() method.
Inside Square ::draw() method.

*/
