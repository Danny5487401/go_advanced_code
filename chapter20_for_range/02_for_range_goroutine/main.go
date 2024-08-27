package main

import (
	"fmt"
	"time"
)

func main() {
	Goroutine1()
	//Goroutine2()
	//Goroutine3()
	//Goroutine4()
}

func Goroutine1() {
	var m = []int{1, 2, 3}
	for index, value := range m {
		go func() {
			fmt.Printf("%v:%v \n", index, value)
		}()
	}
	time.Sleep(time.Second * 3)
}

/*
第一种情况
	各个 goroutine 中输出的 i、v 值都是 for-range 循环结束后的 i、v 最终值，而不是各个 goroutine 启动时的 i, v值。
	因为 goroutine 执行是在后面的某一个时间，使用的是执行时上下文环境的变量值，i，v又相当于一个全局变量，
	协程执行时 for-range 循环已结束，i 和 v 都是最后一次循环的值2和3，所以最后输出都是2和3
*/

func Goroutine2() {
	var m = []int{1, 2, 3}
	for i, v := range m {
		go func() {
			fmt.Println(i, v)
		}()

		time.Sleep(time.Microsecond * 50)

	}
	time.Sleep(time.Second * 3)
}

/*
第二种情况
	一次goroutine的启动准备时间在数十微秒左右。当然该值在不同的操作系统和硬件设备上肯定会存在一些差异。

	这里只是为了讲明白环境上下文，其实我们平时不会这么用的，协程本来就是为了提升并发特性的，如果每次都 sleep 那还有什么意义呐
*/

// 正确方式
func Goroutine3() {

	var m = []int{1, 2, 3}
	for i, v := range m {
		index := i // 这里的 := 会新声明变量，而不是重用
		value := v
		go func() {
			fmt.Println(index, value)
		}()

	}
	time.Sleep(time.Second * 3)

}

func Goroutine4() {

	var m = []int{1, 2, 3}

	for i, v := range m {
		go func(i, v int) {
			fmt.Println(i, v)
		}(i, v)
	}

	time.Sleep(time.Second * 3)

}
