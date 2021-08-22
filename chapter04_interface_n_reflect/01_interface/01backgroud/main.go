package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

/*
为什么要用反射?
	1. 有时你需要编写一个函数，但是并不知道传给你的参数类型是什么，可能是没约定好；也可能是传入的类型很多，这些类型并不能统一表示。这时反射就会用的上了
	2. 有时候需要根据某些条件决定调用哪个函数，比如根据用户的输入来决定。这时就需要对函数和函数的参数进行反射，在运行期间动态地执行函数。

不建议使用反射的原因：
	1。与反射相关的代码，经常是难以阅读的。在软件工程中，代码可读性也是一个非常重要的指标
	2. Go 语言作为一门静态语言，编码过程中，编译器能提前发现一些类型错误，但是对于反射代码是无能为力的。
		所以包含反射相关的代码，很可能会运行很久，才会出错，这时候经常是直接 panic，可能会造成严重的后果
	3. 反射对性能影响还是比较大的，比正常代码运行速度慢一到两个数量级。所以，对于一个项目中处于运行效率关键位置的代码，尽量避免使用反射特性

反射实现原理：
	1.当向接口变量赋予一个实体类型的时候，接口会存储实体的类型信息，反射就是通过接口的类型信息实现的，反射建立在类型的基础上。
	2. Go 语言在 reflect 包里定义了各种类型，实现了反射的各种函数，通过它们可以在运行时检测类型的信息、改变类型的值

Go语言的类型:
	1.变量包括（type, value）两部分,这一点就知道为什么nil != nil了
	2. type 包括 static type和concrete type. 简单来说 static type是你在编码是看见的类型(如int、string_test)，
		concrete type是runtime系统看见的类型
	3. 类型断言能否成功，取决于变量的concrete type，而不是static type。
		因此，一个 reader变量如果它的concrete type也实现了write方法的话，它也可以被类型断言为writer

在反射的概念中， 编译时就知道变量类型的是静态类型；运行时才知道一个变量类型的叫做动态类型。
	静态类型
		//静态类型就是变量声明时的赋予的类型。比如：
		type MyInt int // int 就是静态类型

		type A struct{
			Name string_test  // string就是静态
		}
		var i *int  // *int就是静态类型
	动态类型
		//动态类型：运行时给这个变量赋值时，这个值的类型(如果值为nil的时候没有动态类型)。
		//一个变量的动态类型在运行时可能改变，这主要依赖于它的赋值（前提是这个变量是接口类型）
		var A interface{} // 静态类型interface{}
		A = 10            // 静态类型为interface{}  动态为int
		A = "String"      // 静态类型为interface{}  动态为string
		var M *int
		A = M             // A的值可以改变

		Noted:
		Go语言的反射就是建立在类型之上的，Golang的指定类型的变量的类型是静态的（也就是指定int、string这些的变量，它的type是static type），
		在创建变量的时候就已经确定，反射主要与Golang的interface类型相关（它的type是concrete type），只有interface类型才有反射一说
在Golang的实现中，每个interface变量都有一个对应pair，pair中记录了实际变量的值和类型:(value, type)
	value是实际变量值，type是实际变量的类型。一个interface{}类型的变量包含了2个指针，
	1.一个指针指向值的类型【对应concrete type】，2.另外一个指针指向实际的值【对应value】。
*/
func main() {

	// 一.带函数的interface
	// 类型为*os.File的变量,然后将其赋给一个接口变量r

	var r io.Reader

	tty, err := os.OpenFile("chapter04_reflect/danny_reflect.txt", os.O_RDWR, 0)
	if err != nil {
		fmt.Println("出现错误", err.Error())
	}
	fmt.Printf("tty是%+v\n", tty) // tty是&{file:0xc000058180}

	r = tty
	// 首先声明 r 的类型是 io.Reader，注意，这是 r 的静态类型，此时它的动态类型为 nil，并且它的动态值也是 nil。
	//之后，r = tty 这一语句，将 r 的动态类型变成 *os.File，动态值则变成非空，表示打开的文件对象。
	//	这时，r 可以用<value, type>对来表示为： <tty, *os.File>
	// 接口变量r的pair中将记录如下信息：(tty, *os.File)，这个pair在接口变量的连续赋值过程中是不变的，将接口变量r赋给另一个接口变量w:

	//var w io.Writer
	//w = r.(io.Writer)
	// 接口变量w的pair与r的pair相同，都是:(tty, *os.File)，即使w是空接口类型，pair也是不变的。
	// 和 r 相比，仅仅是 fun 对应的函数变了：Read -> Write

	//interface及其pair的存在，是Golang中实现反射的前提，理解了pair，就更容易理解反射。
	//反射就是用来检测存储在接口变量内部(值value；类型concrete type) pair对的一种机制

	// reflect.TypeOf()是获取pair中的type，reflect.ValueOf()获取pair中的value。
	v := reflect.ValueOf(r) //得到实际的值，通过v我们获取存储在里面的值，还可以去改变值
	t := reflect.TypeOf(r)  //得到类型的元数据,通过t我们能获取类型定义里面的所有元素

	fmt.Printf("reflect.ValueOf(r)==%+v,reflect.TypeOf(r)==%+v\n", v, t) // &{0xc0000ce780} *os.File
	fmt.Printf("t的kind是%+v\n", v.Kind())                                 // t的kind是ptr

	//不带函数的interface
	var empty interface{}
	empty = tty
	fmt.Printf("%T", empty)

}

/*
	反射简单来说就是取得对象的类型(Type)，类别(Kind)，值(Value)，对元素（Element）的字段（Field）进行遍历和操作（读写）。

	对于类型(Type)和类别(Kind)需要注意一下。Type可以认为是Kind的子集
	对于基本类型来说Type和Kind是一致的。例如int的Type和Kind一样都是int
	对于Struct来说，Type是你定义的结构体， Kind为Struct
*/
