package main

import "fmt"

// 内嵌struct :embedded 的特性来删除冗余的代码。当然，代价是初始化会稍微麻烦点
type WithName struct {
	Name string
}

// 打印名字
type Country struct {
	//Name string  //嵌套简化
	WithName
}
type City struct {
	//Name string
	WithName
}

type Printer interface {
	PrintStr()
}

// 绑定接口方法
//func (c Country) PrintStr() {
//fmt.Println(c.Name) //需要简化
//}
//func (c City) PrintStr() {
//fmt.Println(c.Name) //需要简化
//}

// 给单独的WithName绑定方法
func (w WithName) PrintStr() {
	fmt.Println(w.Name)
}

func main() {
	// 初始化
	c1 := Country{WithName{"China"}}
	c2 := City{WithName{"Shanghai"}}
	var cList = []Printer{c1, c2}
	for _, v := range cList {
		v.PrintStr()
	}
}
