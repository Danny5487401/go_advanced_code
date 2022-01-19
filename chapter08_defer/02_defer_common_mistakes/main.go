package main

import "fmt"

// 常见错误及解决： 计划预计的输出c b a,而是输出c c c
type Test struct {
	name string
}

func (t *Test) Close() {
	fmt.Println(t.name, "Closed")
}

func Close(t Test) {
	t.Close()
}

func main() {
	ts := []Test{
		{"a"}, {"b"}, {"c"},
	}
	for _, v := range ts {
		//defer v.Close() // 错误写法

		// 正确写法一： 函数传参
		//Close(v)

		// 正确方法二：复制引用
		v2 := v
		defer v2.Close()
	}
}
