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
	u1 := User{"Danny", 18, Address{
		Province: "Fujian",
		city:     "Fuzhou",
	}}
	u2 := User{"Danny", 18, Address{
		Province: "Fujian",
		city:     "Xiamen",
	}}

	// 默认情况下，cmp.Equal()函数不会比较未导出字段（即字段名首字母小写的字段）。遇到未导出字段，cmp.Equal()直接panic

	// cmdopts.IngoreUnexported忽略未导出字段
	fmt.Println("u1 equals u2?", cmp.Equal(u1, u2, cmpopts.IgnoreUnexported(Address{}))) // u1 equals u2? true

	// cmdopts.AllowUnexported(User{})表示需要比较User的未导出字段：
	fmt.Println("u1 equals u2?", cmp.Equal(u1, u2, cmp.AllowUnexported(User{}.Address))) // u1 equals u2? false

	// cmpopts.IgnoreFields 忽略子段
	opts := cmp.Options{
		cmpopts.IgnoreFields(Address{}, "city"), // 忽略 Address 结构体的 city 字段
	}
	fmt.Println("u1 equals u2?", cmp.Equal(u1, u2, opts)) // u1 equals u2? true

}
