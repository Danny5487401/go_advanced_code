package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/smallnest/chanx"
)

// github.com/smallnest/chanx@v1.2.0/unbounded_chan_test.go
func main() {
	// 初始化
	ch := chanx.NewUnboundedChanSize[int64](context.Background(), 10, 50, 100)

	for i := 1; i < 200; i++ {
		ch.In <- int64(i)
	}

	var count int64
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range ch.Out {
			// 开始读
			count += v
		}
	}()

	for i := 200; i <= 1000; i++ {
		// 开始写
		ch.In <- int64(i)
	}
	close(ch.In)

	wg.Wait()

	fmt.Println(count == 500500)
}
