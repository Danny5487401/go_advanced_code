//性能测试
package main

import "testing"

func array() [1024]int {
	var x [1024]int
	for i := 0; i < len(x); i++ {
		x[i] = i
	}
	return x
}

func slice() []int {
	x := make([]int, 1024)
	for i := 0; i < len(x); i++ {
		x[i] = i
	}
	return x
}

func BenchmarkArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		array()
	}
}

func BenchmarkSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slice()
	}
}

// 命令    go test -bench . -benchmem -gcflags "-N -l"
/*
goos: windows
goarch: amd64
pkg: go_advenced_code/chapter05_slice_n_array/03Array_n_slice_performance
BenchmarkArray-4          307869              3538 ns/op               0 B/op          0 allocs/op
BenchmarkSlice-4          203505              5235 ns/op            8192 B/op          1 allocs/op
PASS
ok      go_advenced_code/chapter05_slice_n_array/03Array_n_slice_performance    2.311s

 */

/*
结论：
	并非所有时候都适合用切片代替数组，因为切片底层数组可能会在堆上分配内存，而且小数组在栈上拷贝的消耗也未必比 make 消耗大
 */
