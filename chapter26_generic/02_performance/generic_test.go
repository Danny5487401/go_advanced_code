package bench_test

import (
	"fmt"
	"testing"
)

func BenchmarkAdd_Generic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		add(i, i)
	}
}
func BenchmarkAdd_NonGeneric(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addInt(i, i)
	}
}

type Addable interface {
	int
}

func add[T Addable](a, b T) T {
	return a + b
}
func addInt(a, b int) int {
	return a + b
}
func main() {
	fmt.Println(add(1, 2))
	fmt.Println(addInt(1, 2))
}
