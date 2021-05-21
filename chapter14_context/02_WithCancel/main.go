package main

// 简化 对于处理单个请求的多个 goroutine 之间与请求域的数据、取消信号、截止时间等相关操

import (
	"context"
	"fmt"
)

/*
Context 接口
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}


1. Deadline方法需要返回当前Context被取消的时间，也就是完成工作的截止时间（deadline）；
2. Done方法需要返回一个Channel，这个Channel会在当前工作完成或者上下文被取消之后关闭，多次调用Done方法会返回同一个Channel；
3. Err方法会返回当前Context结束的原因，它只会在Done返回的Channel被关闭时才会返回非空的值；
	如果当前Context被取消就会返回Canceled错误；
	如果当前Context超时就会返回DeadlineExceeded错误；
4. Value方法会从Context中返回键对应的值，对于同一个上下文来说，多次调用Value 并传入相同的Key会返回相同的结果，该方法仅用于传递跨API和进程间跟请求域的数据
*/

func main() {
	// gen generates integers in a separate goroutine and
	// sends them to the returned channel.
	// The callers of gen need to cancel the context once
	// they are done consuming generated integers not to leak
	// the internal goroutine started by gen.

	/*
		gen函数在单独的goroutine中生成整数并将它们发送到返回的通道。 gen的调用者在使用生成的整数之后需要取消上下文，以免gen启动的内部goroutine发生泄漏。
	*/

	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					return //  return结束该goroutine，防止泄露
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 当我们取完需要的整数后调用cancel

	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
}
