package main

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

func main() {
	ctx := context.TODO()
	var (
		maxWorkers = runtime.GOMAXPROCS(0)
		sem        = semaphore.NewWeighted(int64(maxWorkers))
		out        = make([]int, 3)
	)

	group, _ := errgroup.WithContext(context.Background())
	for i := range out {
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)
			break
		}
		group.Go(func() error {
			go func(i int) {
				defer sem.Release(1)
				out[i] = collatzSteps(i + 1)
			}(i)
			return nil
		})
	}

	// 这里会阻塞，直到所有goroutine都执行完毕
	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(out)
}
func collatzSteps(n int) (steps int) {
	if n <= 0 {
		panic("nonpositive input")
	}

	for ; n > 1; steps++ {
		if steps < 0 {
			panic("too many steps")
		}

		if n%2 == 0 {
			// 偶数
			n /= 2
			continue
		}

		const maxInt = int(^uint(0) >> 1)
		if n > (maxInt-1)/3 {
			panic("overflow")
		}
		n = 3*n + 1
	}

	return steps
}
