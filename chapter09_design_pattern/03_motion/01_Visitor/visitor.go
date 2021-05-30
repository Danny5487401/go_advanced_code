package visitor

import "fmt"

/*  访问者模式

定义：允许一个或者多个操作应用到对象上，解耦操作和对象本身

那么，对一个程序来说，具体的表现就是：

	表面：某个对象执行了一个方法
	内部：对象内部调用了多个方法，最后统一返回结果

举个例子：
	表面：调用一个查询订单的接口
	内部：先从缓存中查询，没查到再去热点数据库查询，还没查到则去归档数据库里查询

做法：
	对象只要预留访问者接口Accept则后期为对象添加功能的时候就不需要改动对象
大概的流程就是:

	1。从结构容器中取出元素
	2。创建一个访问者
	3。将访问者载入传入的元素（即让访问者访问元素）
	4。获取输出


角色组成：
	抽象访问者
	访问者
	抽象元素类
	元素类
	结构容器: (非必须) 保存元素列表，可以放置访问者
 */

/*
源码参考：k8s:--> k8s.io/cli-runtime/pkg/resource/interfaces.go
	// Visitor 即为访问者这个对象
	type Visitor interface {
		Visit(VisitorFunc) error
	}
	// VisitorFunc对应这个对象的方法，也就是定义中的“操作”
	type VisitorFunc func(*Info, error) error

1. VisitorList:封装多个Visitor为一个，出现错误就立刻中止并返回-->k8s.io/cli-runtime/pkg/resource/visitor.go
	// VisitorList定义为[]Visitor，又实现了Visit方法，也就是将多个[]Visitor封装为一个Visitor
	type VisitorList []Visitor

	// 发生error就立刻返回，不继续遍历
	func (l VisitorList) Visit(fn VisitorFunc) error {
		for i := range l {
			if err := l[i].Visit(fn); err != nil {
				return err
			}
		}
		return nil
	}
 */


// 定义元素接口
type Customer interface {
	Accept(Visitor)
}

// 定义访问者接口
type Visitor interface {
	Visit(Customer) // // 访问者的访问方法
}


// 客户列表
type CustomerCol struct {
	customers []Customer
}

func (c *CustomerCol)Add(customer Customer)  {
	c.customers = append(c.customers,customer)
	
}

func (c *CustomerCol)Accept(visitor Visitor)  {
	for _,customer := range c.customers{
		customer.Accept(visitor)
	}
}

// 企业客户
type EnterpriseCustomer struct {
	name string
}

func NewEnterpriseCustomer(name string) *EnterpriseCustomer {
	return &EnterpriseCustomer{
		name: name,
	}
}

func (c *EnterpriseCustomer) Accept(visitor Visitor) {
	visitor.Visit(c)  // 实现元素接口
}


// 个人客户
type IndividualCustomer struct {
	name string
}

func NewIndividualCustomer(name string) *IndividualCustomer {
	return &IndividualCustomer{
		name: name,
	}
}

func (c *IndividualCustomer) Accept(visitor Visitor) {
	visitor.Visit(c)
}


type ServiceRequestVisitor struct{}

func (*ServiceRequestVisitor) Visit(customer Customer) {
	switch c := customer.(type) {
	case *EnterpriseCustomer:
		fmt.Printf("serving enterprise customer %s\n", c.name)
	case *IndividualCustomer:
		fmt.Printf("serving individual customer %s\n", c.name)
	}
}

// only for enterprise
type AnalysisVisitor struct{}

func (*AnalysisVisitor) Visit(customer Customer) {
	switch c := customer.(type) {
	case *EnterpriseCustomer:
		fmt.Printf("analysis enterprise customer %s\n", c.name)
	}
}