package main

import (
	"fmt"
	"reflect"
)

//很多情况下，我们可能并不知道其具体类型，那么这个时候，该如何做呢？需要我们进行遍历探测其Filed来得知

// 定义结构体及tag
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  string `json:"sex"`
}

func (p Person) Say(msg string) {
	fmt.Println("hello，", msg)
}
func (p Person) PrintInfo() {
	fmt.Printf("姓名：%s,年龄：%d，性别：%s\n", p.Name, p.Age, p.Sex)
}

func main() {
	p1 := Person{"danny", 30, "男"}

	// 传实体
	doFiledAndMethod(p1)

	// 传指针
	doFiledAndMethod(&p1)

}

// 通过接口来获取任意参数
func doFiledAndMethod(input interface{}) {

	getType := reflect.TypeOf(input) //先获取input的类型
	if getType.Kind() == reflect.Ptr {
		getType = getType.Elem()
	}
	fmt.Println("get Type is :", getType.Name()) // Person
	fmt.Println("get Kind is :", getType.Kind()) // struct

	getValue := reflect.ValueOf(input)
	fmt.Println("get all Fields is:", getValue) //{danny 30 男}

	if getValue.Kind() == reflect.Ptr {
		getValue = getValue.Elem()
	}

	if getValue.Kind() != reflect.Struct { // 非结构体返回错误提示
		fmt.Printf("ToMap only accepts struct or struct pointer; got %T", getValue)
		return
	}

	// 获取结构体字段
	// 1. 先获取interface的reflect.Type，然后通过NumField进行遍历
	// 2. 再通过reflect.Type的Field获取其Field
	// 3. 最后通过Field的Interface()得到对应的value
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface() //获取第i个值
		// 字段名称:Name, 字段类型:string, 字段数值:danny,tag标签:json:"name"
		tagName := "json"
		tagValue := field.Tag.Get(tagName)
		fmt.Printf("字段名称:%s, 字段类型:%s, 字段数值:%v，tag标签:%v ,json标签对应的值:%v\n", field.Name, field.Type, value, field.Tag, tagValue)
	}

	// 通过反射，操作方法
	// 1. 先获取interface的reflect.Type，然后通过.NumMethod进行遍历
	// 2. 再公国reflect.Type的Method获取其Method
	for i := 0; i < getType.NumMethod(); i++ {
		method := getType.Method(i)
		// 方法名称:PrintInfo, 方法类型:func(main.Person)
		// 方法名称:Say, 方法类型:func(main.Person, string)
		fmt.Printf("方法名称:%s, 方法类型:%v \n", method.Name, method.Type)
	}
	fmt.Println("---------------")
}
