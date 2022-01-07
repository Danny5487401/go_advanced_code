/* 使用channel传递channel
场景：使用场景有点多，通常是用来获取结果。
原理：channel可以用来传递变量，channel自身也是变量，可以传递自己。
用法：下面示例展示了有序展示请求的结果
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// handle 处理请求，耗时随机模拟
func handle(wg *sync.WaitGroup, a int) chan int {
	out := make(chan int)
	go func() {
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
		out <- a
		wg.Done()
	}()
	return out
}

func main() {
	reqs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	// 存放结果的channel的channel
	outs := make(chan chan int, len(reqs))
	var wg sync.WaitGroup
	wg.Add(len(reqs))
	for _, x := range reqs {
		o := handle(&wg, x)
		outs <- o
	}

	go func() {
		wg.Wait()
		close(outs)
	}()

	// 读取结果，结果有序
	for o := range outs {
		fmt.Println(<-o)
	}
}


