package main

import "fmt"

var done = make(chan bool)
var msg string

//func main() {
//	ch := make(chan int)
//	ch <- 5
//}

func main() {
	go aGoroutine()
	done <- true
	fmt.Println(msg)

}
func aGoroutine() {
	msg = "hello world"
	<-done
}

/*
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
	E:/go_advanced_code/chapter03_channel/01_basicUse/03deadlock/main_test.go:5 +0x57
*/
