package main

import "fmt"

// 更好的方法
/*
方法三：
	嵌套+将具体的实现Country和City私有化，不对外暴露实现细节

*/

// 内嵌struct :embedded 的特性来删除冗余的代码。当然，代价是初始化会稍微麻烦点
//type WithName struct {
//	Name string
//}
type WithTypeName struct {
	Type string // 多一个tYpe字段，代表是Country还是City,还可以其他
	Name string
}

// 打印名字
type Country struct {
	// Name string
	// WithName //嵌套简化
	WithTypeName
}
type City struct {
	// Name string
	// WithName  //嵌套简化
	WithTypeName
}

type Printer interface {
	PrintStr()
}

// 给单独的WithName绑定方法
//func (w WithName) PrintStr() {
//	fmt.Println(w.Name)
//}
func (c WithTypeName) PrintStr() {
	fmt.Printf("%s:%s\n", c.Type, c.Name)
}

// 内部初始化复杂，返回接口类型
func NewCountry(name string) Printer {
	return Country{WithTypeName{Type: "Country", Name: name}}
}
func NewCity(name string) Printer {
	return City{WithTypeName{Type: "City", Name: name}}
}
func main() {
	// 外部初始化要求简单
	//c1 := Country{WithName{"China"}} // 初始化复杂
	//c2 := City{WithName{"Shanghai"}} // 初始化复杂

	// 初始化简单
	c1 := NewCountry("China")
	c2 := NewCity("Shanghai")
	var cList = []Printer{c1, c2}
	for _, v := range cList {
		v.PrintStr()
	}
}
