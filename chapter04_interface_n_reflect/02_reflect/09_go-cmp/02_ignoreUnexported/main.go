package main

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type Address struct {
	Province string
	city     string
}

type User struct {
	Name    string
	Age     int
	Address Address
}

func main() {
	u1 := User{"dj", 18, Address{
		Province: "Fujian",
		city:     "Fuzhou",
	}}
	u2 := User{"dj", 18, Address{
		Province: "Fujian",
		city:     "Xiamen",
	}}

	// 默认情况下，cmp.Equal()函数不会比较未导出字段（即字段名首字母小写的字段）。遇到未导出字段，cmp.Equal()直接panic

	// cmdopts.IngoreUnexported忽略未导出字段
	fmt.Println("u1 equals u2?", cmp.Equal(u1, u2, cmpopts.IgnoreUnexported(Address{})))
	// cmdopts.AllowUnexported(User{})表示需要比较User的未导出字段：
	fmt.Println("u1 equals u2?", cmp.Equal(u1, u2, cmp.AllowUnexported(User{}.Address)))
}
