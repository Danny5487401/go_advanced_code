package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculate(t *testing.T) {
	// 方式一 if
	//if Calculate(2) != 4 {
	//	t.Error("Expected 2 + 2 to equal 4")
	//}
	// 方式二 简短
	assert.Equal(t, Calculate(2), 3, "应该相等")
	assert.Equalf(t, Calculate(4), 7, "应该%d相等7", Calculate(4))

	// 表驱动
	// 初始化Assertions 对象
	assert := assert.New(t)

	var tests = []struct {
		input    int
		expected int
	}{
		{2, 4},
		{-1, 1},
		{0, 2},
		{-5, -3},
		{99999, 100001},
	}

	for _, test := range tests {
		// 不用传参数t
		assert.Equal(Calculate(test.input), test.expected)
	}
}
