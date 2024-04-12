package main

import (
	"fmt"
)

// 下次 GC 的时机: 可以通过一个环境变量 GOGC 来控制，默认是 100 ，即增长 100% 的堆内存才会触发 GC。

// 官方的解释是，如果当前使用了 4M 内存，那么下次 GC 将会在内存达到 8M 的时候

func main() {
	fmt.Println("start.")

	container := make([]int, 8)
	fmt.Println("> loop.")
	for i := 0; i < 32*1000*1000; i++ {
		container = append(container, i)
	}

	fmt.Println("< loop.")
}

// ➜  01_next_gc_stage git:(feature/gc) ✗ go build myenumstr.go
// ➜  01_next_gc_stage git:(feature/gc) ✗ GODEBUG=gctrace=1 ./main

/*
> loop.
gc 1 @0.002s 2%: 0.010+1.1+0.017 ms clock, 0.081+0/0.52/0.14+0.14 ms cpu, 4->6->2 MB, 5 MB goal, 8 P
gc 2 @0.004s 22%: 0.006+5.0+2.5 ms clock, 0.050+0/0.067/1.3+20 ms cpu, 5->14->10 MB, 6 MB goal, 8 P
gc 3 @0.016s 9%: 0.043+13+0.034 ms clock, 0.35+0/0.17/6.8+0.27 ms cpu, 20->27->12 MB, 21 MB goal, 8 P
gc 4 @0.029s 8%: 0.028+3.2+0.013 ms clock, 0.22+0.040/0.25/3.1+0.10 ms cpu, 21->21->8 MB, 25 MB goal, 8 P
gc 5 @0.033s 8%: 0.027+1.6+0.003 ms clock, 0.22+0/0.25/1.5+0.029 ms cpu, 19->19->10 MB, 20 MB goal, 8 P
gc 6 @0.035s 7%: 0.023+3.0+0.004 ms clock, 0.19+0/0.10/2.9+0.032 ms cpu, 24->24->13 MB, 25 MB goal, 8 P
gc 7 @0.039s 4%: 0.021+24+0.020 ms clock, 0.17+0/1.8/0.11+0.16 ms cpu, 30->51->38 MB, 31 MB goal, 8 P
gc 8 @0.064s 3%: 0.028+21+0.012 ms clock, 0.23+0/0.10/21+0.10 ms cpu, 64->64->26 MB, 76 MB goal, 8 P
gc 9 @0.087s 3%: 0.049+6.0+0.003 ms clock, 0.39+0/0.13/6.1+0.030 ms cpu, 59->59->33 MB, 60 MB goal, 8 P
gc 10 @0.093s 3%: 0.026+4.6+0.012 ms clock, 0.21+0.005/0.079/4.6+0.098 ms cpu, 74->74->41 MB, 75 MB goal, 8 P
gc 11 @0.099s 2%: 0.12+41+0.003 ms clock, 1.0+0/0.12/41+0.029 ms cpu, 93->93->51 MB, 94 MB goal, 8 P
gc 12 @0.142s 2%: 0.023+10+0.003 ms clock, 0.18+0/0.11/10+0.029 ms cpu, 116->116->64 MB, 117 MB goal, 8 P
gc 13 @0.153s 1%: 0.025+24+0.004 ms clock, 0.20+0/0.10/24+0.032 ms cpu, 145->145->80 MB, 146 MB goal, 8 P
gc 14 @0.180s 5%: 0.043+80+0.005 ms clock, 0.35+0/80/0.097+0.042 ms cpu, 182->182->101 MB, 183 MB goal, 8 P
gc 15 @0.261s 4%: 0.025+27+0.020 ms clock, 0.20+0/0.12/27+0.16 ms cpu, 227->227->227 MB, 228 MB goal, 8 P
gc 16 @0.403s 2%: 0.24+144+0.003 ms clock, 1.9+0/0.10/144+0.031 ms cpu, 582->582->197 MB, 583 MB goal, 8 P
gc 17 @0.551s 2%: 0.053+35+0.035 ms clock, 0.42+0/0.10/35+0.28 ms cpu, 444->444->444 MB, 445 MB goal, 8 P
< loop.

*/
