package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	d := time.Now().Add(5 * time.Millisecond)

	ctx, cancel := context.WithDeadline(context.Background(), d)

	// Even though ctx will be expired, it is good practice to call its
	// cancelation function in any case. Failure to do so may keep the
	// context and its parent alive longer than necessary.
	// 尽管ctx会过期，但在任何情况下调用它的cancel函数都是很好的实践。 如果不这样做，可能会使上下文及其父类存活的时间超过必要的时间
	defer cancel()

	select {
	case t := <-time.After(1 * time.Second):
		fmt.Println("overslept:", t)
	case <-ctx.Done():
		t, _ := ctx.Deadline()
		err := ctx.Err()
		fmt.Printf("deadline:%s, err:%v", t.String(), err) // cdeadline:2025-12-16 17:17:15.21407 +0800 CST m=+0.005130417 ,err:context deadline exceeded
	}
}
