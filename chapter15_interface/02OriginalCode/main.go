package main

import "fmt"

/*
背景：
	接口的定义可以说是一种规范，是一组方法的集合，通常在代码设计的层面，对多个组件有共性的方法进行抽象(共性可以分为横向和纵向)引入一层中间层，
		解除上游与下游的耦合关系，让代码可读性更高并不用关心方法的具体实现，同时借助接口也可以实现多态。
	共性可以分为横向和纵向的
	纵向：
		例如动物这个对象可以向下细分为狗和猫，它们有共同的行为可以跑。
	横向
		再或者数据库的连接可以抽象为接口，可以支持mysql、oracle等
源码：
	interface的定义在1.15.3源码包runtime中,interface的定义分为两种，一种是不带方法的runtime.eface和带方法的runtime.iface。
	1。runtime.eface表示不含方法的interface{}类型,结构体包含可以表示任意数据类型的_type和存储指定的数据data,data用指针来表示
	type eface struct {
		_type *_type
		data  unsafe.Pointer
	}
	type _type struct {
		size       uintptr //占用的字节大小
		ptrdata uintptr //指针数据 size of memory prefix holding all pointers
		hash       uint32 //计算的hash
		tflag      tflag //额外的标记信息
		align      uint8 //内存对齐系数
		fieldAlign uint8 //字段内存对齐系数
		kind uint8 //用于标记数据类型
		// function for comparing objects of this type
		// (ptr to object A, ptr to object B) -> ==?
		equal func(unsafe.Pointer, unsafe.Pointer) bool//用于判断当前类型多个对象是否相等
		str       nameOff //名字偏移量
		ptrToThis typeOff //指针的偏移量
	}
	2。runtime.iface表示包含方法的接口,结构体包含itab和data数据,itab包含的是接口类型interfacetype
		和装载实体的任意类型_type以及实现接口的方法fun,fun是可变大小,go在编译期间就会对接口实现校验检查,并将对应的方法存储fun。
		type iface struct {
			tab  *itab
			data unsafe.Pointer
		}

		type itab struct {
			inter *interfacetype //接口类型的表示
			_type *_type
			hash  uint32 // copy of _type.hash. Used for type switches.
			_     [4]byte
			fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
		}

*/
// 带方法的interface iface
type Print interface {
	Print()
}

type User struct {
	Name string
	Age  int
}

func (r User) Print() {
	fmt.Printf("hello %v,Age %v", r.Name, r.Age)
}

func main() {
	//eface 不带方法
	var y interface{} = User{}
	fmt.Printf("%+v\n", y) // {Name: Age:0}

	//iface 带方法
	var p Print = User{Name: "小明", Age: 18}
	fmt.Printf("%+v\n", p) // {Name:小明 Age:18}
	p.Print()
}
