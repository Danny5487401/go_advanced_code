package bench

import (
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
