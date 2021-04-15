package main

import "fmt"

func main() {
	var a chan int
	if a == nil {
		fmt.Println("channel 是 nil 的, 不能使用，需要先创建通道。。")
		a = make(chan int)
		fmt.Printf("数据类型是： %T\n", a)
	}

	ch1 := make(chan int)
	fmt.Printf("%T,%p\n",ch1,ch1)

	test1(ch1)
}
func test1(ch chan int){
	// channel是引用类型的数据，在作为参数传递的时候，传递的是内存地址
	fmt.Printf("%T,%p\n",ch,ch)
}