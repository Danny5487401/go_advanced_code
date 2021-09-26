package main

import "fmt"

func sum(a, b int) {
	c := a + b
	fmt.Println("sum:", c)
}

func f(a, b int) {
	defer sum(a, b)

	fmt.Printf("a: %d, b: %d\n", a, b)
}

func main() {
	a, b := 1, 2
	f(a, b)
}

//汇编命令：go tool compile -N -S -L defer_asm.go>defer_asm.s 2>&1
