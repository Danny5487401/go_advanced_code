package main

import "fmt"

/*
	从一个有缓冲的 channel 里读数据，当 channel 被关闭，依然能读出有效值。只有当返回的 ok 为 false 时，读出的数据才是无效的
*/

func main() {
	ch := make(chan int, 5)
	ch <- 18
	close(ch)
	x, ok := <-ch
	if ok {
		fmt.Println("获取到值", x)
	}
	x, ok = <-ch
	if !ok {
		fmt.Println("通道关闭")
	}
}
