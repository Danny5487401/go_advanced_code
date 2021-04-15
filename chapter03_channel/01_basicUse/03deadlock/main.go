package main

func main() {
	ch := make(chan int)
	ch <- 5
}

/*
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
	E:/go_advanced_code/chapter03_channel/01_basicUse/03deadlock/main.go:5 +0x57
 */