package main

/*
一.用到反射的包：
	官方包：sort swapper,sql convertValue ,Json反序列化
	第三方包： proto reflect,sqlx scanAll
二。源码f分析
	1. reflect.Type 是以一个接口的形式存在的
	2. reflect.Value 是以一个结构体的形式存在
	接口变量，实际上都是由一 pair 对（type 和 data）组合而成，pair 对中记录着实际变量的值和类型。
		也就是说在真实世界里，type 和 value 是合并在一起组成 接口变量的。
		而在反射的世界里，type 和 data 却是分开的，他们分别由 reflect.Type 和 reflect.Value 来表现


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
