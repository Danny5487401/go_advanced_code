package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

func wrapper(i int, wg *sync.WaitGroup) func() {
	// 任务是无参数的，使用闭包传入
	return func() {
		fmt.Printf("hello from task:%d\n", i)
		time.Sleep(1 * time.Second)
		wg.Done()
	}
}

func main() {
	p, _ := ants.NewPool(4, ants.WithMaxBlockingTasks(2)) // 最大等待队列长度
	defer p.Release()

	var wg sync.WaitGroup
	wg.Add(8)
	for i := 1; i <= 8; i++ {
		go func(i int) {
			err := p.Submit(wrapper(i, &wg))
			if err != nil {
				// 超过最大等待队列长度这个长度，提交任务直接返回错误
				fmt.Printf("task:%d err:%v\n", i, err)
				wg.Done()
			}
		}(i)
	}

	wg.Wait()
}
