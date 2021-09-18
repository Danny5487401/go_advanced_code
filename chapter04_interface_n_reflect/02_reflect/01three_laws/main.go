package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 定律1：从接口值到反射对象  普通变量 -> 接口interface变量 -> 反射对象
	//反射操作：通过反射，可以获取一个接口类型变量的 类型和数值
	Law1()

	// 定律2：反射对象->interface变量
	Law2()

	// 定律3：当反射对象所存的值是可设置时，反射对象才可修改
	//Go 语言里的函数都是值传递，只要你传递的不是变量的指针，你在函数内部对变量的修改是不会影响到原始的变量的。
	//回到反射上来，当你使用 reflect.Typeof 和 reflect.Valueof 的时候，如果传递的不是接口变量的指针，
	//反射世界里的变量值始终将只是真实世界里的一个拷贝，你对该反射对象进行修改，并不能反映到真实世界里。
	// 可以使用Value.CanSet()检测是否可以通过此Value修改原始变量的值,可寻址的
	Law3()
}
func Law1() {
	// 1。内置类型
	var x float64 = 3.4

	fmt.Println("type:", reflect.TypeOf(x))   //type: float64 获取某个变量的静态类型
	fmt.Println("value:", reflect.ValueOf(x)) //value: 3.4  获取某个变量的值

	fmt.Println("-------------------")
	//根据反射的值，来获取对应的类型和数值
	v := reflect.ValueOf(x)
	// reflect.Value.Kind() Kind获取变量值的底层类型（类别），注意不是类型，是Int、Float，还是Struct，还是Slice
	fmt.Println("kind is float64: ", v.Kind() == reflect.Float64) // kind is float64:  true

	// 获取变量值的类型，效果等同于reflect.TypeOf
	fmt.Println("type : ", v.Type())   // type :  float64
	fmt.Println("value : ", v.Float()) // value :  3.4

	// 2。自定义类型
	type MyInt int
	var y MyInt = 7
	v1 := reflect.ValueOf(y)                         // 接口变量存储了实际的值和值的类型
	fmt.Println("kind 底层类型 is float64: ", v1.Kind()) // kind 底层类型 is float64:  int
	// 获取变量值的类型，效果等同于reflect.TypeOf
	fmt.Println("type 定义类型 : ", v1.Type()) // type 定义类型 :  main.MyInt

	//3。结构体类型
	coder := &Coder{Name: "danny"}
	typ := reflect.TypeOf(coder)
	val := reflect.ValueOf(coder)
	typeofStringer := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	fmt.Println("kind of coder : ", typ.Kind()) // kind of coder :  ptr
	fmt.Println("type of coder : ", typ)        // type of coder :  *main.Coder
	fmt.Println("value of coder : ", val)       // value of coder :  danny

	// 确认coder是否实现了String() string接口
	fmt.Println("implements of stringer:", typ.Implements(typeofStringer)) // implements of stringer: true
}
func Law2() {
	coder := &Coder{Name: "danny"}
	val := reflect.ValueOf(coder)
	c, ok := val.Interface().(*Coder)
	if ok {
		fmt.Println(c.Name) // danny
	} else {
		panic("type assert to *Coder error")
	}
}
func Law3() {
	z := 10 // 不是指针
	v2 := reflect.ValueOf(z)
	fmt.Println("settable:", v2.CanSet()) // settable: false
	p := reflect.ValueOf(&z)
	fmt.Println("settable:", p.CanSet())  // settable: false
	v3 := p.Elem()                        // 指针指向的元素
	fmt.Println("settable:", v3.CanSet()) // settable: true
	// Value.SetXXX()系列函数可设置Value中原始对象的值。 根据Value.Kind()的结果去获得变量的底层类型，然后选用该类别的Set函数。
	switch v3.Kind() {
	case reflect.Int:
		v3.SetInt(100) //设置值   源码会检测是否可设置mustBeAssignable,然后检测可导出
	default:
		fmt.Println("未知类型")
	}

	fmt.Printf("z新值是%d", z) // 100
}

type Coder struct {
	Name string
}

func (c *Coder) String() string {
	return c.Name
}
