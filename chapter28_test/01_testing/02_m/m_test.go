package _2_m

import (
	"fmt"
	"testing"
)

// testing.M函数可以在测试函数执行之前做一些其他操作

func TestMain(m *testing.M) {
	fmt.Println("测试执行第一步: main开始测试,是在测试之前执行的")
	m.Run()
}

func TestUser(t *testing.T) {
	fmt.Println("main 测试")
	t.Run("测试执行第二步:开始测试第一个子测试函数:", testAddUser)
}

func testAddUser(t *testing.T) {
	fmt.Println("测试执行第三步:sub 1 测试")
}
