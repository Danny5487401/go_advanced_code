package main

import (
	"fmt"
	"time"
)

func main() {

	/*
		func After(d Duration) <-chan Time
			返回一个通道：chan，存储的是d时间间隔后的当前时间。
	*/
	ch1 := time.After(3 * time.Second) //3s后
	fmt.Printf("%T\n", ch1) // <-chan time.Time
	fmt.Println(time.Now()) //2021-04-15 10:24:06.8008246 +0800 CST m=+0.031982001
	time2 := <-ch1
	fmt.Println(time2) //2021-04-15 10:24:09.8000272 +0800 CST m=+3.031184601


}