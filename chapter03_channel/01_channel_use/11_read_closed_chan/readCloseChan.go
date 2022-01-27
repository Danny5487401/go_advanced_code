package main

import "fmt"

/*
	从一个有缓冲的 channel 里读数据，当 channel 被关闭，依然能读出有效值。只有当返回的 ok 为 false 时，读出的数据才是无效的
*/

func main() {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	close(ch)
	x, ok := <-ch
	if ok {
		fmt.Println("获取到值", x)
	}
	x, ok = <-ch
	if !ok {
		fmt.Println("通道关闭")
	}
	x, ok = <-ch
	if !ok {
		fmt.Println("通道关闭")
	}
	a := <-ch // 不会阻塞
	fmt.Printf("通道已经关闭,数值%v\n", a)
	fmt.Printf("类型%T\n", a)

}
