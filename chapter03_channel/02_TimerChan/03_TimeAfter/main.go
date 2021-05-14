package main

import (
	"errors"
	"fmt"
	"time"
)

func main() {

	/*
		func After(d Duration) <-chan Time
			返回一个通道：chan，存储的是d时间间隔后的当前时间。
	*/
	// 简单理解
	//ch1 := time.After(3 * time.Second) //3s后
	//fmt.Printf("%T\n", ch1)            // <-chan time.Time
	//fmt.Println(time.Now())            //2021-04-15 10:24:06.8008246 +0800 CST m=+0.031982001
	//time2 := <-ch1
	//fmt.Println(time2) //2021-04-15 10:24:09.8000272 +0800 CST m=+3.031184601

	/*为操作加上超时
	场景：需要超时控制的操作
	原理：使用select和time.After，看操作和定时器哪个先返回，处理先完成的，就达到了超时控制的效果
	*/
	rsp, err := doWithTimeOut(3 * time.Second)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("成功返回%+v", rsp)
}

func doWithTimeOut(timeout time.Duration) (int, error) {
	fmt.Println(time.Now())
	select {
	case ret := <-do():
		return ret, nil
	case time2 := <-time.After(timeout):
		fmt.Println(time2)
		return 0, errors.New("timeout")
	}
}

func do() <-chan int {
	outCh := make(chan int)
	go func() {
		//outCh <- 2 // 看是否返回
	}()
	return outCh
}
