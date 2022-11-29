package main

import (
	"testing"
)

func TestVisitor(t *testing.T) {
	square := NewSquare(2)
	circle := &circle{radius: 3}
	rectangle := &rectangle{l: 2, b: 3}
	// 打印信息
	t.Log(square.getType())

	// 计算面积
	areaCalculator := &areaCalculator{}
	square.accept(areaCalculator)
	circle.accept(areaCalculator)
	rectangle.accept(areaCalculator)

	// 计算中点
	middleCoordinates := &middleCoordinates{}
	square.accept(middleCoordinates)
	circle.accept(middleCoordinates)
	rectangle.accept(middleCoordinates)
}
