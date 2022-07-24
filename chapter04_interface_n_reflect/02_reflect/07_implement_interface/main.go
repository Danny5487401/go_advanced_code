package main

import (
	"fmt"
	"reflect"
)

type Tester interface {
	Do() string
}

type User struct {
	Name string
	Age  int
}

func (u User) Do() string {
	return "do it"
}

func main() {
	var u User

	t := reflect.TypeOf(u)

	//将nil转成Tester接口指针，然后再通过反射,Elem()方法获取指针对应的接口类型
	ele := reflect.TypeOf((*Tester)(nil)).Elem()
	fmt.Println("打印接口名称", ele.String()) // 获得接口类型

	if t.Implements(ele) {
		fmt.Println("实现了Tester接口 !!!")
	}
}
