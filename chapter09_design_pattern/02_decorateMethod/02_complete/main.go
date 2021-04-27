package main

/*
装饰模式：一种动态地往一个类中添加新的行为的设计模式.
从process2中得知主要包含四个角色，抽象构件 Component，具体构件 ConcreteComponent，抽象装饰类 Decorator，具体装饰类 ConcreteComponent
优点：
1. 可以通过一种动态的方式来扩展一个对象的功能
2. 可以使用多个具体装饰类来装饰同一对象，增加其功能
3. 具体组件类与具体装饰类可以独立变化，符合“开闭原则”

缺点：
1. 对于多次装饰的对象，易于出错，排错也很困难
2. 对于产生很多具体装饰类 ，增加系统的复杂度以及理解成本

使用场景：
1.需要给一个对象增加功能，这些功能可以动态地撤销，例如：在不影响其他对象的情况下，动态、透明的方式给单个对象添加职责，处理那些可以撤销的职责
2.需要给一批兄弟类增加或者改装功能
 */

import (
	"fmt"
)
//定义一个抽象组件
type Company interface {
	Showing()
}

//实现Company的一个组件
type BaseCompany struct {
}

func (pB *BaseCompany) Showing() {
	fmt.Println("基础公司有老板，有前台，有人事...")
}

//实现Company的一个组件
type DevelopingCompany struct {
	Company
}

func (pD *DevelopingCompany) AddWorker() {
	fmt.Println("发展中公司还有开发、测试、财务人员")
}

func (pD *DevelopingCompany) Showing() {
	fmt.Println("发展中公司中：")
	pD.Company.Showing()
	pD.AddWorker()
}

//实现Company的一个组件
type BigCompany struct {
	Company
}

func (pD *BigCompany) AddWorker() {
	fmt.Println("大公司除此之外，个职能人员应有尽有")
}

func (pD *BigCompany) Showing() {
	fmt.Println("大型公司中：")
	pD.Company.Showing()
	pD.AddWorker()
}

func main() {

	company := &BaseCompany{}
	developingCompany := &DevelopingCompany{Company: company}
	developingCompany.Showing()

	bigCompany := &BigCompany{Company: developingCompany}
	bigCompany.Showing()
	return
}


