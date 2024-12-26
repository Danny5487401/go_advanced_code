package _1_nil_size_n_addr

import (
	"fmt"
	"unsafe"
)

func main() {
	// 1. 各类型为nil时的地址
	addr()

	// 2. 各类型为nil时的大小
	size()

}

func addr() {

	var p *int
	var c chan int
	var f func()
	var m map[int]int
	var s []int
	var i interface{}
	fmt.Printf("*int地址是%p\n", p)      // 0x0
	fmt.Printf("chan int地址是%p\n", c)  // 0x0
	fmt.Printf("函数地址是%p\n", f)        // 0x0
	fmt.Printf("map地址是%p\n", m)       // 0x0
	fmt.Printf("切片地址是%p\n", s)        // 0x0
	fmt.Printf("interface地址是%p\n", i) // %!p(<nil>)
	/*
		从代码中可以直观的看到指针、管道、函数、map、切片slice为nil时输出的地址都为0x0，可以验证不同类型nil值地址都是相同的。
		而其中比较特殊的是接口，输出的是%!p(<nil>)，大致的原因是因为nil的接口经由reflect.ValueOf()函数输出的类型为<invalid reflect.Value>，
		针对于这种类型Printf函数进行了特别的拼接最终得到%!p(<nil>)


	*/
	type People struct {
		name string
		age  int
	}
	//实例化是有地址的
	var p1 = &People{}
	var p2 = People{}
	var p3 *People
	fmt.Printf("p1地址%p\n", p1)                    // 0xc00000c060
	fmt.Printf("p2地址%p\n", &p2)                   // 0xc00000c080
	fmt.Printf("p3地址%p\n", p3)                    // 0x0
	fmt.Println(p1 == nil, &p2 == nil, p3 == nil) // false false true
}

func size() {

	var p *int = nil
	fmt.Println("int: ", unsafe.Sizeof(p)) // 8
	var c chan int = nil
	fmt.Println("chan int: ", unsafe.Sizeof(c)) // 8
	var f func() = nil
	fmt.Println("func: ", unsafe.Sizeof(f)) //8
	var m map[int]int = nil
	fmt.Println("map: ", unsafe.Sizeof(m)) // 8
	var s []int = nil
	fmt.Println("slice: ", unsafe.Sizeof(s)) // 24
	var i interface{} = nil
	fmt.Println("interface: ", unsafe.Sizeof(i)) // 16
}
