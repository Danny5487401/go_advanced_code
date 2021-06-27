package main

/*
semaphore的作用
	信号量是在并发编程中比较常见的一种同步机制，它会保证持有的计数器在0到初始化的权重之间，
	每次获取资源时都会将信号量中的计数器减去对应的数值，在释放时重新加回来，当遇到计数器大于信号量大小时就会进入休眠等待其他进程释放信号。
	go中的semaphore，提供sleep和wakeup原语，使其能够在其它同步原语中的竞争情况下使用。
	当一个goroutine需要休眠时，将其进行集中存放，当需要wakeup时，再将其取出，重新放入调度器中。
需求：
	考拉兹猜想，是指对于每一个正整数，如果它是奇数，则对它乘3再加1，如果它是偶数，则对它除以2，如此循环，最终都能够得到1。
	通过信号量实现并发对 考拉兹猜想的示例，对1-32之间的数字进行计算，并打印32个符合结果的值
*/

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"golang.org/x/sync/semaphore"
)

// Example_workerPool demonstrates how to use a semaphore to limit the number of
// goroutines working on parallel tasks.
//
// This use of a semaphore mimics a typical “worker pool” pattern, but without
// the need to explicitly shut down idle workers when the work is done.
func main() {
	ctx := context.TODO()

	// 权重值为逻辑cpu个数
	var (
		maxWorkers = runtime.GOMAXPROCS(0)
		sem        = semaphore.NewWeighted(int64(maxWorkers))
		out        = make([]int, 32)
	)

	// Compute the output using up to maxWorkers goroutines at a time.
	for i := range out {
		// When maxWorkers goroutines are in flight, Acquire blocks until one of the
		// workers finishes.
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)
			break
		}

		go func(i int) {
			defer sem.Release(1)
			out[i] = collatzSteps(i + 1)
		}(i)
	}

	// 为了保证每次for循环都会正常结束，最后调用了 sem.Acquire(ctx, int64(maxWorkers)) ，
	//	表示最后一次执行必须获取的权重值为 maxWorkers。当然如果使用 errgroup 同步原语的话，这一步可以省略掉
	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
		log.Printf("Failed to acquire semaphore: %v", err)
	}

	fmt.Println(out)

}

// collatzSteps computes the number of steps to reach 1 under the Collatz
// conjecture. (See https://en.wikipedia.org/wiki/Collatz_conjecture.)
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

		const maxInt = int(^uint(0) >> 1) // ^   ^作一元运算符表示是按位取反
		if n > (maxInt-1)/3 {
			panic("overflow")
		}
		n = 3*n + 1
	}

	return steps
}
