package factorymethod

import "testing"

func compute(factory OperatorFactory, a, b int) int {
	op := factory.Create()
	op.SetA(a)
	op.SetB(b)
	return op.Result()
}

// 操作符工厂
func TestOperator(t *testing.T) {
	var (
		factory OperatorFactory
	)

	// 加法
	factory = PlusOperatorFactory{}
	if compute(factory, 1, 2) != 3 {
		t.Fatal("error with factory method pattern")
	}

	// 减法
	factory = MinusOperatorFactory{}
	if compute(factory, 4, 2) != 2 {
		t.Fatal("error with factory method pattern")
	}
}
