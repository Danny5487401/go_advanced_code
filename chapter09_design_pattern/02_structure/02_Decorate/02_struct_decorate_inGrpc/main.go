package main

import "fmt"

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
