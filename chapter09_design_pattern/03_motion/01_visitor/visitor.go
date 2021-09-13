package visitor

import "fmt"

// 1.定义元素接口
type Customer interface {
	Accept(Visitor)
}

// 2.定义访问者接口
type Visitor interface {
	Visit(Customer) // // 访问者的访问方法
}

//3。对象容器

//客户群
type CustomerCol struct {
	customers []Customer
}

func (c *CustomerCol) Add(customer Customer) {
	c.customers = append(c.customers, customer)

}

func (c *CustomerCol) Accept(visitor Visitor) {
	for _, customer := range c.customers {
		customer.Accept(visitor)
	}
}

//企业级客户
type EnterpriseCustomer struct {
	name string
}

func NewEnterpriseCustomer(name string) *EnterpriseCustomer {
	return &EnterpriseCustomer{
		name: name,
	}
}

func (c *EnterpriseCustomer) Accept(visitor Visitor) {
	visitor.Visit(c) // 实现元素接口
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

//访问器1
type ServiceRequestVisitor struct{}

func (*ServiceRequestVisitor) Visit(customer Customer) {
	switch c := customer.(type) {
	case *EnterpriseCustomer:
		fmt.Printf("serving enterprise customer %s\n", c.name)
	case *IndividualCustomer:
		fmt.Printf("serving individual customer %s\n", c.name)
	}
}

//访问器2
type AnalysisVisitor struct{}

func (*AnalysisVisitor) Visit(customer Customer) {
	switch c := customer.(type) {
	case *EnterpriseCustomer:
		fmt.Printf("analysis enterprise customer %s\n", c.name)
	}
}
