package main

import (
	"fmt"
	"time"
)

/*
双向：
chan T

单向：

chan <- T,

	只支持写 send-only

<- chan T,

	只读 receive-only
*/
func main() {
	c := make(chan int) // an unbuffered channel
	go func(ch chan<- int, x int) {
		time.Sleep(time.Second)
		// <-ch    // fails to compile
		// Send the value and block until the result is received.
		ch <- x * x // 9 is sent
	}(c, 3)
	done := make(chan struct{})
	go func(ch <-chan int) {
		// Block until 9 is received.
		n := <-ch
		fmt.Println(n) // 9
		// ch <- 123   // fails to compile
		time.Sleep(time.Second)
		done <- struct{}{}
	}(c)
	// Block here until a value is received by
	// the channel "done".
	<-done
	fmt.Println("bye")
}
