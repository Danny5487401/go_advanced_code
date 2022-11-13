package hello

import (
	"context"
	"fmt"
)

// greet函数
func greet() {
	var msg = "Hello World!"
	fmt.Println(msg)
}

// Foo 结构体
type Foo struct {
	i int
}

// Bar 接口
type Bar interface {
	Do(ctx context.Context) error
}
