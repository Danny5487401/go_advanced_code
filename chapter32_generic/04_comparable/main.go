package main

import (
	"fmt"
)

func Equal[T comparable](param1, param2 T) bool {
	return param1 == param2
}

// Person 成员变量都是可比较的，所以编译通过并且能得到正确结果
type Person struct {
	Name string
	Age  int
}

// NewPerson 有一个 Address 是 切片，切片是不可比较的
type NewPerson struct {
	Person
	Address []string
}

func main() {
	fmt.Println(Equal[Person](
		Person{
			Name: "Danny",
			Age:  21,
		},
		Person{
			Name: "Joy",
			Age:  22,
		},
	))

	// 编译器提示: NewPerson does not satisfy comparable
	//fmt.Println(Equal[NewPerson](
	//	NewPerson{
	//		Person: Person{
	//			Name: "Danny",
	//			Age:  21,
	//		},
	//		Address: []string{"厦门"},
	//	},
	//	NewPerson{
	//		Person: Person{
	//			Name: "Joy",
	//			Age:  22,
	//		},
	//		Address: []string{"福州"},
	//	},
	//))
}
