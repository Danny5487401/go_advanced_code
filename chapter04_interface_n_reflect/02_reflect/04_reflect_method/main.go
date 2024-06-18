package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 一. 使用反射修改值
	changeValue()

	// 二. 通过反射，调用结构体方法
	getStructFuncAndCall()

	// 三. 反射调用普通函数
	getFuncTypeAndCallByParams()

}

type myFloat64 float64

func changeValue() {
	var num myFloat64 = 1.2345

	// 通过reflect.ValueOf获取num中的reflect.Value，注意，参数必须是指针才能修改其值
	pointer := reflect.ValueOf(&num)
	if pointer.Kind() != reflect.Ptr {
		// 如果reflect.ValueOf的参数不是指针，会如何？
		//03PointerSetPrivateValue = reflect.ValueOf(num)
		//newValue = 03PointerSetPrivateValue.Elem() // 如果非指针，这里直接panic，“panic: reflect: call of reflect.Value.Elem on float64 Value”
		fmt.Println("不可以修改值")
		return

	}
	newValue := pointer.Elem()

	fmt.Println("type of num:", newValue.Type())          // float64
	fmt.Println("settability of num:", newValue.CanSet()) // true

	// 重新赋值
	newValue.SetFloat(77)
	fmt.Println("修改后的值", num) // 77
}

// 如何通过反射来进行方法的调用？
// 本来可以用结构体对象.方法名称()直接调用的，
// 但是如果要通过反射， 那么首先要将方法注册，也就是MethodByName，然后通过反射调动mv.Call

func getStructFuncAndCall() {
	p2 := Person{"Danny", 30, "男"}
	// 1. 要通过反射来调用起对应的方法，必须要先通过reflect.ValueOf(interface)来获取到reflect.Value，
	// 得到“反射类型对象”后才能做下一步处理
	getValue := reflect.ValueOf(p2)

	// 2.一定要指定参数为正确的方法名

	// 先看看没有参数的调用方法
	methodValue1 := getValue.MethodByName("PrintInfo")
	fmt.Printf("Kind : %s, Type : %s\n", methodValue1.Kind(), methodValue1.Type()) // Kind : func, Type : func()
	methodValue1.Call(nil)                                                         // 没有参数，直接写nil
	args1 := make([]reflect.Value, 0)                                              // 或者创建一个空的切片也可以
	methodValue1.Call(args1)

	// 有参数的方法调用
	methodValue2 := getValue.MethodByName("Say")
	fmt.Printf("Kind : %s, Type : %s\n", methodValue2.Kind(), methodValue2.Type()) // Kind : func, Type :  func(string)
	args2 := []reflect.Value{reflect.ValueOf("反射机制")}
	methodValue2.Call(args2)

	// 多个不同类型的参数
	methodValue3 := getValue.MethodByName("Test")
	fmt.Printf("Kind : %s, Type : %s ， 入参个数：%d，出参个数：%d \n",
		methodValue3.Kind(), methodValue3.Type(), methodValue3.Type().NumIn(), methodValue3.Type().NumOut()) // Kind : func, Type : func(int, int, string) error
	args3 := []reflect.Value{reflect.ValueOf(100), reflect.ValueOf(200), reflect.ValueOf("Hello")}
	rsp := methodValue3.Call(args3)

	// 获取返回值数组、验证数组的长度以及类型并打印其中的数据；
	fmt.Println("获取结果的 type ：", rsp[0].Type()) //error
	if len(rsp) != 1 || rsp[0].Kind() != reflect.Invalid {
		return
	}
}

// 通过反射，调用普通函数
func getFuncTypeAndCallByParams() {
	value := reflect.ValueOf(fun1)
	fmt.Printf("Kind : %s , Type : %s\n", value.Kind(), value.Type()) //Kind : func , Type : func()

	value2 := reflect.ValueOf(fun2)
	fmt.Printf("Kind : %s , Type : %s\n", value2.Kind(), value2.Type()) //Kind : func , Type : func(int, string)

	// 通过反射调用函数
	value.Call(nil)

	value2.Call([]reflect.Value{reflect.ValueOf(100), reflect.ValueOf("hello")})
}

func fun1() {
	fmt.Println("我是函数fun1()，无参的。。")
}

func fun2(i int, s string) {
	fmt.Println("我是函数fun2()，有参数。。", i, s)
}

type Person struct {
	Name string
	Age  int
	Sex  string
}

func (p Person) Say(msg string) {
	fmt.Println("hello，", msg)
}
func (p Person) PrintInfo() {
	fmt.Printf("姓名：%s,年龄：%d，性别：%s \n", p.Name, p.Age, p.Sex)
}

func (p Person) Test(i, j int, s string) error {
	fmt.Println(i, j, s)
	return nil
}
