package main

import (
	"errors"
	"fmt"
	"time"
)

/*
为操作加上超时：
	func After(d Duration) <-chan Time
		返回一个通道：chan，存储的是d时间间隔后的当前时间。

场景：需要超时控制的操作
原理：使用select和time.After，看操作和定时器哪个先返回，处理先完成的，就达到了超时控制的效果
*/

func main() {

	// 为操作加超时
	rsp, err := doWithTimeOut(3 * time.Second)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("成功返回%+v", rsp)
}

func doWithTimeOut(timeout time.Duration) (int, error) {
	select {
	case ret := <-do():
		return ret, nil
	case time2 := <-time.After(timeout):
		fmt.Println("超时了,时间是", time2)
		return 0, errors.New("timeout")
	}
}

// 执行的业务逻辑
func do() <-chan int {
	outCh := make(chan int)
	go func() {
		fmt.Println("执行业务逻辑")

		// 场景：
		// 1. 注释: 会阻塞，返回没有值，会超时
		// 2. 不注释：返回有值，不超时
		//outCh <- 2 // 看是否返回
	}()
	return outCh
}
