package main

import "fmt"

/*
用到反射的包：
	官方包：sort swapper,sql convertValue ,Json反序列化
	第三方包： proto reflect,sqlx scanAll
源码
	reflect/type.go
	type定义了接口，rtype实现了接口

	func TypeOf(i interface{}) Type {
		eface := *(*emptyInterface)(unsafe.Pointer(&i))
		return toType(eface.typ)
	}

	// emptyInterface is the header for an interface{} value.
	//跟eface一样，不过eface用于运行时,emptyInterface用于反射
	type emptyInterface struct {
		typ  *rtype
		word unsafe.Pointer  //数据
	}

	// reflect/value.go
	type Value struct {
		typ *rtype
		ptr unsafe.Pointer
		flag  //元信息
	}

*/

type Person struct {
	age int
}

func (p Person) howOld() int {
	return p.age
}

func (p *Person) growUp() {
	p.age += 1
}

func main() {
	// qcrao 是值类型
	qcrao := Person{age: 18}

	// 值类型 调用接收者也是值类型的方法
	fmt.Println(qcrao.howOld())

	// 值类型 调用接收者是指针类型的方法
	qcrao.growUp()
	fmt.Println(qcrao.howOld())

	// ----------------------

	// stefno 是指针类型
	stefno := &Person{age: 100}

	// 指针类型 调用接收者是值类型的方法
	fmt.Println(stefno.howOld())

	// 指针类型 调用接收者也是指针类型的方法
	stefno.growUp()
	fmt.Println(stefno.howOld())
}
