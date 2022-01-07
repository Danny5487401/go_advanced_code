package main

import "fmt"

// 双向通道

/*
	双向：
		chan T -->
			chan <- data,写出数据，写
			data <- chan,获取数据，读
	单向：定向
		chan <- T,
			只支持写，
		<- chan T,
			只读
*/
func main() {

	ch1 := make(chan string) // 双向，可读，可写
	done := make(chan bool)
	go sendData(ch1, done)
	data := <-ch1 //阻塞
	fmt.Println("子goroutine传来：", data)
	ch1 <- "我是main。。" // 阻塞

	<-done
	fmt.Println("main...over....")
}

//子goroutine-->写数据到ch1通道中
//main goroutine-->从ch1通道中取
func sendData(ch1 chan string, done chan bool) {
	ch1 <- "我是小明" // 阻塞
	data := <-ch1 // 阻塞
	fmt.Println("main goroutine传来：", data)

	done <- true
}

/*
子goroutine传来： 我是小明
main goroutine传来： 我是main。。
main...over....
*/
