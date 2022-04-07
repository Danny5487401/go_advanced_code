package main

import (
	"fmt"
	"github.com/smallnest/chanx"
	"sync"
)

func main() {
	ch := chanx.NewUnboundedChanSize(10, 50, 100)

	for i := 1; i < 200; i++ {
		ch.In <- int64(i)
	}

	var count int64
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		for v := range ch.Out {
			count += v.(int64)
		}
	}()

	for i := 200; i <= 1000; i++ {
		ch.In <- int64(i)
	}
	close(ch.In)

	wg.Wait()

	fmt.Println(count)
}
