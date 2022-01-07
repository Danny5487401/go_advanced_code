package main

import "fmt"


func main() {
	var ch1 chan bool       //声明，没有创建
	fmt.Println(ch1)        //<nil>
	fmt.Printf("%T\n", ch1) //chan bool
	ch1 = make(chan bool)   //0xc0000a4000,是引用类型的数据
	fmt.Println(ch1)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println("子goroutine中，i：", i)
		}
		// 循环结束后，向通道中写数据，表示要结束了。。
		ch1 <- true

		fmt.Println("结束。。")

	}()
	// 一个通道发送和接收数据，默认是阻塞的
	data := <-ch1 // 从ch1通道中读取数据
	fmt.Println("data-->", data)
	fmt.Println("main。。over。。。。")
}