package main

import "fmt"

/*
	1.当指针为nil时不能对指针进行解引用,  *p 称为 解引用 或者 间接引用
	2.结构体指针为nil时，能够调用结构体指针类型实现的方法，且该方法不能包含结构体的属性
	3.结构体指针为nil时，不能调用结构体类型实现的方法，也不能调用使用结构体属性的方法
*/

type People struct {
	Name string
}

func (p *People) Born() {
	fmt.Println("Hello World1~~~")
}

func (p *People) Born2() {
	fmt.Println("Hello World2~~~", p.Name)
}

func (p People) Born3() {
	fmt.Println("Hello World3~~~")
}

func main() {
	var pPeople *People

	// pass example
	pPeople.Born()

	// fail example
	// _ = *pPeople    // panic: runtime error: invalid memory address or nil pointer dereference
	// pPeople.Born2() // panic: runtime error: invalid memory address or nil pointer dereference
	// pPeople.Born3() // panic: runtime error: invalid memory address or nil pointer dereference
}
