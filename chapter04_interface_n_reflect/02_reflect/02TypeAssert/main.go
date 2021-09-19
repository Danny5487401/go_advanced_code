package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 1.空接口断言 eface
	var u interface{}
	u = Person{"Danny"}
	val, ok := u.(int)
	if !ok {
		fmt.Printf("空接口断言为int类型的零值:%v\n", val)
	}

	// 2.非空接口断言 iface
	var t Tester
	t = Person{"Danny"}

	//第一种情况:断言具体类型
	check(t)

	//第二种情况:断言接口类型
	// CALL runtime.assertI2I2(SB)
	if t, ok2 := t.(Tester2); ok2 { //重用变量名t（无需重新声明）
		//若类型断言为true，则新的t被转型为Tester2接口类型，但其动态类型和动态值不变
		check2(t)
	}

	//2.反射断言
	ReflectToGetType()
}

func ReflectToGetType() {

	var num float64 = 1.2345
	pointer := reflect.ValueOf(&num)
	value := reflect.ValueOf(num)
	fmt.Println("指针的值:", pointer, "float64的值:", value) // 指针的值: 0xc0000160b0 float64的值: 1.2345

	// 转换成interface进行断言
	convertPointer := pointer.Interface().(*float64)
	convertValue := value.Interface().(float64)

	fmt.Println("类型断言后，指针的值:", convertPointer, "float64的值:", convertValue) // 类型断言后，指针的值: 0xc0000160b0 float64的值: 1.2345

}

//===接口=====
type Tester interface {
	getName() string
}
type Tester2 interface {
	printName()
}

//===Person类型====
type Person struct {
	name string
}

func (p Person) getName() string {
	return p.name
}
func (p Person) printName() {
	fmt.Println("Hello" + p.name)
}

func check(t Tester) {
	//使用switch写法
	switch t.(type) {
	case Person:
		p := t.(Person)
		fmt.Printf("具体类型判断：%T,\t调用方法结果：%s\n", p, p.getName())
	default:
		fmt.Println("failed to match")
	}
}
func check2(t Tester2) {
	t.printName()
}
