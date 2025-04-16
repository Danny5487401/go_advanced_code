package bench

type Addable interface {
	int
}

func add[T Addable](a, b T) T {
	return a + b
}
func addInt(a, b int) int {
	return a + b
}
