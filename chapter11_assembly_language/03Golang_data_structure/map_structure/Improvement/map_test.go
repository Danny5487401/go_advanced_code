package Improvement

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"time"
)

/*
优化背景
	当 map 的 key/value 都是非指针类型的话，扫描是可以避免的，直接标记整个 map 的颜色（三色标记法）就行了，不用去扫描每个 bmap 的 overflow 指针
案例
	map[string]int -> map[[12]byte]int
原因
	因为 string 底层有指针，所以当 string 作为 map 的 key 时，GC 阶段会扫描整个 map；而数组 [12]byte 是一个值类型，不会被 GC 扫描。

*/

func MapWithPointer() {
	const N = 10000000
	m := make(map[string]string)
	for i := 0; i < N; i++ {
		n := strconv.Itoa(i)
		m[n] = n
	}
	now := time.Now()
	runtime.GC()
	fmt.Printf("With a map of strings, GC took: %s\n", time.Since(now))

	// 引用一下防止被 GC 回收掉
	_ = m["0"]
}

func MapWithoutPointer() {
	const N = 10000000
	m := make(map[int]int)
	for i := 0; i < N; i++ {
		str := strconv.Itoa(i)
		// hash string to int
		n, _ := strconv.Atoi(str)
		m[n] = n
	}
	now := time.Now()
	runtime.GC()
	fmt.Printf("With a map of int, GC took: %s\n", time.Since(now))

	_ = m[0]
}

func TestMapWithPointer(t *testing.T) {
	MapWithPointer()
}

func TestMapWithoutPointer(t *testing.T) {
	MapWithoutPointer()
}

// 验证了 string 相对于 int 这种值类型对 GC 的消耗更大
// Go语言使用 map 时尽量不要在 big map 中保存指针
