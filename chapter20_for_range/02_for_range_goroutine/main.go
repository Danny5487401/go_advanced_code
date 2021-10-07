package main

import (
	"fmt"
	"time"
)

func main() {
	Goroutine3()
}

func Goroutine1() {
	var m = []int{1, 2, 3}
	for i, v := range m {
		go func() {
			fmt.Println(i, v)
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
		if i == 0 {
			time.Sleep(time.Second * 1)
		}
	}
	time.Sleep(time.Second * 3)
}

/*
第二种情况
	第一次遍历后 sleep 了1秒,所以第一次循环中的协程有时间执行了，开始执行时当前上下文中 i 和 v 的值还是第一次遍历的0和1，
	后面的没 sleep 就是最后循环结束时的2和3了。

	这里只是为了讲明白环境上下文，其实我们平时不会这么用的，协程本来就是为了提升并发特性的，如果每次都 sleep 那还有什么意义呐
*/

//正确方式
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
