package main

import "fmt"

/*
需求：
	国家，城市，地区打印他们自己的名字
方法一：
	 绑定接口方法，经常这样写
缺点：
	如果我要实现N个Printer，就要定义N个strcut+N个PrintStr()方法。
*/

type Country struct {
	Name string
}
type City struct {
	Name string
}

type Printer interface {
	PrintStr()
}

// 绑定接口方法，这样写存在大量冗余写法
func (c Country) PrintStr() {
	fmt.Println(c.Name) //
}
func (c City) PrintStr() {
	fmt.Println(c.Name)
}

func main() {
	// 初始化
	c1 := Country{"China"}
	c2 := City{"Shanghai"}
	var cList = []Printer{c1, c2}
	for _, v := range cList {
		v.PrintStr()
	}
}
