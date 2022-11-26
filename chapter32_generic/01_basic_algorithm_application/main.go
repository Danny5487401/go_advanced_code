package main

import "fmt"

// 方式一：常规写法
func SumIntsOrFloats1[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

// 方式二：将类型约束声明为接口
type Number interface {
	int64 | float64
}

func SumIntsOrFloats2[K comparable, V Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func main() {
	// Initialize a map for the integer values
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}

	// Initialize a map for the float values
	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	fmt.Printf("Generic Sums: %v and %v\n",
		SumIntsOrFloats1[string, int64](ints),
		SumIntsOrFloats1[string, float64](floats))

	fmt.Printf("Generic Sums: %v and %v\n",
		SumIntsOrFloats2[string, int64](ints),
		SumIntsOrFloats2[string, float64](floats))
}
