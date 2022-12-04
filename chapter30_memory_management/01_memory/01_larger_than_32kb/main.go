package main

import (
	"fmt"
)

func main() {
	var m = make([]int, 10024)
	fmt.Println(m[0])
}

// arm64  架构
// go tool compile -S main.go
