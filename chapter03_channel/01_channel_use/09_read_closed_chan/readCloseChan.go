package main

import "fmt"

/*
	从一个有缓冲的 channel 里读数据，当 channel 被关闭，依然能读出有效值。当返回的 ok 为 false 时，读出的数据是无效的
*/

func main() {
	ch := make(chan int64, 1)
	ch <- int64(1)
	close(ch)
	x, ok := <-ch
	if ok {
		fmt.Println("获取到值", x)
	}
	x, ok = <-ch
	if !ok {
		fmt.Printf("通道关闭获取第一次,数值是零值:%v，类型: %T \n", x, x)
	}
	x, ok = <-ch
	if !ok {
		fmt.Printf("通道关闭获取第一次,数值是零值:%v，类型: %T \n", x, x)
	}
	a := <-ch // 不会阻塞
	fmt.Printf("通道已经关闭,不阻塞,数值是零值:%v，类型: %T\n", a, a)

}
