package main

import "fmt"

type SomeStruct struct{}

func main() {
	var h *SomeStruct
	var wrapper interface{} = h
	// Go预置的 nil 是没有类型的，为了让 wrapper 和 nil 进行比较，编译器会首先将 nil 转化为一个 interface{}，然后进行比较。但是注意，因为 nil 无默认类型，即便转为 interface{}，它也是没有对应的动态类型的，跟 wrapper 的动态类型*SomeStruct 不匹配，所以会返回 false。
	fmt.Println(h == nil, wrapper == nil, wrapper == h) // true false true
	fmt.Println((interface{})(nil) == (*int)(nil))      // false
}
