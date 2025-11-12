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

type Robot struct {
}

func main() {
	// 1. 实现接口
	var u User
	t := reflect.TypeOf(u)

	// 将nil转成Tester接口指针，然后再通过反射,Elem()方法获取指针对应的接口类型
	ele := reflect.TypeOf((*Tester)(nil)).Elem()
	fmt.Println("打印包含路径的接口名称:", ele.String()) // 打印包含路径的接口名称: main.Tester

	if t.Implements(ele) {
		fmt.Printf("%v 实现了 %v 接口 !!!\n", t.Name(), ele.Name()) // User 实现了 Tester 接口 !!!
	}

	// 2. 未实现接口
	var r Robot
	t2 := reflect.TypeOf(r)
	if !t2.Implements(ele) {
		fmt.Printf("%v 没有实现 %v 接口 !!!\n", t2.Name(), ele.Name()) // Robot 没有实现 Tester 接口 !!!
	}

}
