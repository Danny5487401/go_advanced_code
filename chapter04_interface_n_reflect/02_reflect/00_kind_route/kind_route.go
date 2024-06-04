package kind_route

import (
	"fmt"
	"reflect"
)

func CollectUserInfo(param interface{}) {
	val := reflect.ValueOf(param)
	switch val.Kind() {
	case reflect.String:
		fmt.Println("姓名:", val.String())
	case reflect.Struct:
		fmt.Println("姓名:", val.FieldByName("Name"))
		fmt.Println("年龄:", val.FieldByName("Age"))
		fmt.Println("住址:", val.FieldByName("Address"))
		fmt.Println("电话:", val.FieldByName("Phone"))
	case reflect.Ptr:
		fmt.Println("姓名:", val.Elem().FieldByName("Name"))
		fmt.Println("年龄:", val.Elem().FieldByName("Age"))
		fmt.Println("住址:", val.Elem().FieldByName("Address"))
		fmt.Println("电话:", val.Elem().FieldByName("Phone"))
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			fmt.Println("姓名:", val.Index(i).FieldByName("Name"))
			fmt.Println("年龄:", val.Index(i).FieldByName("Age"))
			fmt.Println("住址:", val.Index(i).FieldByName("Address"))
			fmt.Println("电话:", val.Index(i).FieldByName("Phone"))
		}
	case reflect.Map:
		itr := val.MapRange()
		for itr.Next() {
			fmt.Println("用户ID:", itr.Key())
			fmt.Println("姓名:", itr.Value().Elem().FieldByName("Name"))
			fmt.Println("年龄:", itr.Value().Elem().FieldByName("Age"))
			fmt.Println("住址:", itr.Value().Elem().FieldByName("Address"))
			fmt.Println("电话:", itr.Value().Elem().FieldByName("Phone"))
		}
	default:
		panic("unsupport type !!!")
	}
}

type User struct {
	Name    string
	Age     int64
	Address string
	Phone   int64
}

func array() {
	arr := [5]int{1, 2, 3, 4, 5}
	typ := reflect.TypeOf(arr)
	fmt.Println(typ.Kind())       // array
	fmt.Println(typ.Len())        // 5
	fmt.Println(typ.Comparable()) // true

	elemTyp := typ.Elem()
	fmt.Println(elemTyp.Kind())       // int
	fmt.Println(elemTyp.Comparable()) // true
}

func slice() {
	s := make([]int, 5, 10)
	typ := reflect.TypeOf(s)
	fmt.Println(typ.Kind()) // slice
	fmt.Println(typ.Elem()) // int
}

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

func (p Person) unexportedMethod() {
}

func structInfo() {
	p1 := Person{"danny", 30, "男"}

	// 传指针
	doFiledAndMethod(&p1)

}

// 通过接口来获取任意参数
func doFiledAndMethod(input interface{}) {

	getType := reflect.TypeOf(input) // 先获取input的类型
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
		value := getValue.Field(i).Interface() //获取第i个值,后面需要断言获取不同的具体类型
		// 字段名称:Name, 字段类型:string, 字段数值:danny,tag标签:json:"name"
		tagName := "json"
		tagValue := field.Tag.Get(tagName)
		fmt.Printf("字段名称:%s, 字段类型:%s, 字段数值:%v，tag标签:%v, json标签对应的值:%v\n", field.Name, field.Type, value, field.Tag, tagValue)
	}

	// 通过反射，操作方法
	// 1. 先获取interface的reflect.Type，然后通过.NumMethod进行遍历
	// 2. 再通过 reflect.Type的Method获取其Method
	for i := 0; i < getType.NumMethod(); i++ {
		method := getType.Method(i)
		// 方法名称:PrintInfo, 方法类型:func(main.Person)
		// 方法名称:Say, 方法类型:func(main.Person, string)
		fmt.Printf("方法名称:%s, 方法类型:%v \n", method.Name, method.Type)
	}

}

func channelInfo() {
	ch := make(chan<- int, 10)
	ch <- 1
	ch <- 2
	typ := reflect.TypeOf(ch)
	fmt.Println(typ.Kind())    // chan
	fmt.Println(typ.Elem())    // int
	fmt.Println(typ.ChanDir()) // chan<-

	fmt.Println(reflect.ValueOf(ch).Len()) // 2
	fmt.Println(reflect.ValueOf(ch).Cap()) // 10
}

func mapInfo() {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	typ := reflect.TypeOf(m)
	fmt.Println(typ.Kind()) // map
	fmt.Println(typ.Key())  // string
	fmt.Println(typ.Elem()) // int

	fmt.Println(reflect.ValueOf(m).Len()) // 3
}

func pointerInfo() {
	i := 10
	p := &i
	typ := reflect.TypeOf(p)
	fmt.Println(typ.Kind()) // ptr
	fmt.Println(typ.Elem()) // int
}

func interfaceInfo() {
	var a Animal = Cat{}
	typ := reflect.TypeOf(a)
	fmt.Println(typ.Kind())         // interface
	fmt.Println(typ.NumMethod())    // 1
	fmt.Println(typ.Method(0).Name) // Speak
	fmt.Println(typ.Method(0).Type) // func(main.Animal) string
}

type Animal interface {
	Speak() string
}

type Cat struct{}

func (c Cat) Speak() string {
	return "Meow"
}

func funcInfo() {
	typ := reflect.TypeOf(foo)
	fmt.Println(typ.Kind())                      // func
	fmt.Println(typ.NumIn())                     // 3
	fmt.Println(typ.In(0), typ.In(1), typ.In(2)) // int int *int
	fmt.Println(typ.NumOut())                    // 2
	fmt.Println(typ.Out(0))                      // int
	fmt.Println(typ.Out(1))                      // bool
}

func foo(a, b int, c *int) (int, bool) {
	*c = a + b
	return *c, true
}
