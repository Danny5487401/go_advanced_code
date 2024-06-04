package main

import "testing"

// TestSub 内部调用sub1、sub2和sub3三个子测试
func TestSub(t *testing.T) {
	// setup code

	t.Run("A=1", sub1)
	t.Run("A=2", sub2)
	t.Run("B=1", sub3)

	// tear-down code
}

/*
=== RUN   TestSub
=== RUN   TestSub/A=1
=== RUN   TestSub/A=2
=== RUN   TestSub/B=1
--- PASS: TestSub (0.00s)
    --- PASS: TestSub/A=1 (0.00s)
    --- PASS: TestSub/A=2 (0.00s)
    --- PASS: TestSub/B=1 (0.00s)
PASS

*/

/*
通过上面的例子我们知道Run()方法第一个参数为子测试的名字，而实际上子测试的内部命名规则为："<父测试名字>/<传递给Run的名字>"。比如，传递给Run()的名字是“A=1”，那么子测试名字为“TestSub/A=1”。这个在上面的命令行输出中也可以看出
*/
