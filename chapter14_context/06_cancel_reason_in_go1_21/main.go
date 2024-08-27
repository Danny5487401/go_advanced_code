package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	myctxErr := fmt.Errorf("自定义ctxErr错误")

	ctx, cancel := context.WithTimeoutCause(context.Background(), 1*time.Second, myctxErr)
	defer cancel()
	time.Sleep(2 * time.Second)
	fmt.Println(context.Cause(ctx))
}
