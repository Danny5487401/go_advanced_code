package main

import "testing"

func TestFactory(t *testing.T) {

	// 创建工厂
	factory := ShapeFactory{}

	// 使用该工厂，通过传递类型信息来获取实体类的对象。
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