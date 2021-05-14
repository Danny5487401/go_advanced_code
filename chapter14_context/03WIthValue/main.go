package main

import (
	"context"
	"fmt"
)

func main() {
	// 提供的键必须是可比性和应该不是字符串类型或任何其他内置的类型以避免包使用的上下文之间的碰撞。WithValue 用户应该定义自己的键的类型。
	type favContextKey string

	f := func(ctx context.Context, k favContextKey) {
		if v := ctx.Value(k); v != nil {
			fmt.Println("found value:", v)
			return
		}
		fmt.Println("key not found:", k)
	}

	k := favContextKey("language")
	ctx := context.WithValue(context.Background(), k, "Go")

	f(ctx, k)
	f(ctx, favContextKey("color"))
}
