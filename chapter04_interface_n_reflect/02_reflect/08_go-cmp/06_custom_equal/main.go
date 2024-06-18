package main

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"strings"
)

// 如果在类型上定义的Equal方法不合适，可以动态地转换类型，不执行Equal方法(或任何方法)。
type otherString string
type myString otherString

// Equal 自定义类型的Equal方法
func (sf otherString) Equal(s otherString) bool {
	return strings.EqualFold(string(sf), string(s))
}

func main() {
	// NOTE: 如果不想让指定类型进行执行equal方法,可以进行类型转换.
	trans := cmp.Transformer("type_change", func(in otherString) myString {
		return myString(in)
	})

	x := []otherString{"foo", "bar", "baz"}
	y := []otherString{"fOO", "bAr", "Baz"}

	fmt.Println(cmp.Equal(x, y))        // true 相等,并且不区分大小写
	fmt.Println(cmp.Equal(x, y, trans)) // false 不相等,类型转换过,不会执行equal方法

	a1 := NetAddr{"127.0.0.1", 5000}
	a2 := NetAddr{"localhost", 5000}

	fmt.Println("a1 equals a2?", cmp.Equal(a1, a2, cmp.Transformer("localhost", transformLocalhost))) // a1 equals a2? true

}

type NetAddr struct {
	IP   string
	Port int
}

// 函数格式 trFunc  // func(T) R
func transformLocalhost(a NetAddr) NetAddr {
	if a.IP == "localhost" {
		return NetAddr{IP: "127.0.0.1", Port: a.Port}
	}

	return a
}
