package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {

	var ch = make(chan bool, 20000)
	var begin = make(chan bool)

	go func() {
		// 注释掉runtime.LockOSThread()，有2-60倍的时间差
		runtime.LockOSThread()
		<-begin
		fmt.Println("begin")
		tm := time.Now()
		for i := 0; i < 10000000; i++ {
			<-ch
		}
		fmt.Println(time.Now().Sub(tm))
		os.Exit(0)
	}()

	for i := 0; i < 50000; i++ {
		// 负载
		go func() {
			var count int
			load := 100000
			for {
				count++
				if count >= load {
					count = 0
					runtime.Gosched()
				}
			}
		}()
	}

	for i := 0; i < 20; i++ {
		go func() {
			for {
				ch <- true
			}
		}()
	}

	fmt.Println("all start")
	begin <- true

	select {}
}
