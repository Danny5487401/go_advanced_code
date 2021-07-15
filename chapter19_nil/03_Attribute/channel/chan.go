package main

/*
	当一个chan为nil时，向chan中发送数据则会永远阻塞
	当一个chan为nil时，接收chan数据则会永远阻塞
	close一个为nil的chan则会panic
	当一个chan为nil时，则会屏蔽select中的case
*/

func main() {
	//var c chan int
	// c <- 3   // fatal error: all goroutines are asleep - deadlock!
	// <-c      // fatal error: all goroutines are asleep - deadlock!
	//close(c) // panic: close of nil channel
}
