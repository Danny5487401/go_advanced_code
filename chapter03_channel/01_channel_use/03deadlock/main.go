package main

import "fmt"

func main() {
	length := 10
	slice := make([]int, 0, length)
	for i := 1; i <= length; i++ {
		// 1. 构建数据
		slice = append(slice, i)
	}

	sChan := make(chan int, 5)

	go func() {
		for _, v := range slice {
			// 2. 发送数据
			sChan <- v
		}
		// 演示未关闭导致死锁
		//close(sChan)
	}()

	for data := range sChan {
		// 接受数据
		go func(data int) {
			fmt.Println(data)
		}(data)
	}
}

/*
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
	E:/github.com/Danny5487401/go_advanced_code/chapter03_channel/01_basicUse/03deadlock/main_test.go:5 +0x57
*/
