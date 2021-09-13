package visitor

import "testing"

func TestExampleRequestVisitor(t *testing.T) {
	// 创建一个对象容器
	c := &CustomerCol{}
	// 添加元素
	c.Add(NewEnterpriseCustomer("A company"))
	c.Add(NewEnterpriseCustomer("B company"))
	c.Add(NewIndividualCustomer("bob"))
	// 开始接受 访问器1
	c.Accept(&ServiceRequestVisitor{})
	t.Log("---------------")
	// 开始接受 访问器2
	c.Accept(&AnalysisVisitor{})

}

/*
Output:
	serving enterprise customer A company
	serving enterprise customer B company
	serving individual customer bob
		visitor_test.go:14: ---------------
	analysis enterprise customer A company
	analysis enterprise customer B company
*/
