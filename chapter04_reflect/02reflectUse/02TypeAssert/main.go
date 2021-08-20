package main

import (
	"fmt"
	"reflect"
)

/*


一。断言类型的语法：x.(T)，这里x表示一个接口的类型，T表示一个类型（也可为接口类型）。
一个类型断言检查一个接口对象x的动态类型是否和断言的类型T匹配。

类型断言分两种情况：
	第一种，如果断言的类型T是一个具体类型，类型断言x.(T)就检查x的动态类型是否和T的类型相同。

		1。如果这个检查成功了，类型断言的结果是一个类型为T的对象，该对象的值为接口变量x的动态值。换句话说，具体类型的类型断言从它的操作对象中获得具体的值。
		2。如果检查失败，接下来这个操作会抛出panic，除非用两个变量来接收检查结果，如：f, ok := w.(*os.File)
	第二种，如果断言的类型T是一个接口类型，类型断言x.(T)检查x的动态类型是否满足T接口。

		1。如果这个检查成功，则检查结果的接口值的动态类型和动态值不变，但是该接口值的类型被转换为接口类型T。换句话说，对一个接口类型的类型断言改变了类型的表述方式，改变了可以获取的方法集合（通常更大），但是它保护了接口值内部的动态类型和值的部分。
		2。如果检查失败，接下来这个操作会抛出panic，除非用两个变量来接收检查结果，如：f, ok := w.(io.ReadWriter)
注意：

如果断言的操作对象x是一个nil接口值，那么不论被断言的类型T是什么这个类型断言都会失败。
我们几乎不需要对一个更少限制性的接口类型（更少的方法集合）做断言，因为它表现的就像赋值操作一样，除了对于nil接口值的情况。

二。反射的方式
reflect类型断言
	// 从reflect.Value中获取接口interface的信息
	// realValue := value.Interface().(已知的类型)
	// 可以理解为“强制转换”，但是需要注意的时候，转换的时候，如果转换的类型不完全符合，则直接panic
	// Golang 对类型要求非常严格，类型一定要完全符合
	// 如下两个，一个是*float64，一个是float64，如果弄混，则会panic

	1. 从 Value 到实例
		该方法最通用，用来将 Value 转换为空接口，该空接口内部存放具体类型实例
		可以使用接口类型查询去还原为具体的类型
		//func (v Value) Interface() （i interface{})

*/

func main() {
	// 1.接口断言
	var t Tester
	t = Person{"Danny"}
	check(t)

	//2.反射断言
}

func ReflectToGetType() {

	var num float64 = 1.2345
	pointer := reflect.ValueOf(&num)
	value := reflect.ValueOf(num)
	fmt.Println(pointer) // 0xc00009c058
	fmt.Println(value)   // 1.2345

	// 转换成interface进行断言
	convertPointer := pointer.Interface().(*float64)
	convertValue := value.Interface().(float64)

	fmt.Println(convertPointer) // 0xc00009c058
	fmt.Println(convertValue)   // 1.2345
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
	//第一种情况:具体类型
	if f, ok1 := t.(Person); ok1 {
		fmt.Printf("%T\n%s\n", f, f.getName())
	}
	//第二种情况:扩展到其他类型
	if t, ok2 := t.(Tester2); ok2 { //重用变量名t（无需重新声明）
		check2(t) //若类型断言为true，则新的t被转型为Tester2接口类型，但其动态类型和动态值不变
	}
}
func check2(t Tester2) {
	t.printName()
}
