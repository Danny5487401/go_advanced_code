package main

import "fmt"

// 1.创建接口
type Manager interface {
	HaveRight(money int) bool                     //有权做
	HandleFeeRequest(name string, money int) bool //处理费用请求
}

// 2, 创建请求链
type RequestChain struct {
	Manager
	successor *RequestChain //后续者(链式结构，链表)
}

func (rc *RequestChain) HaveRight(money int) bool {
	return true
}
func (rc *RequestChain) HandleFeeRequest(name string, money int) bool {
	//先判断金额在对应的部门是否有权限
	if rc.Manager.HaveRight(money) {
		return rc.Manager.HandleFeeRequest(name, money) //如果有权限，处理请求
	}
	//继任者不为空,传递到继任者处理
	if rc.successor != nil {
		return rc.successor.HandleFeeRequest(name, money)
	}
	return false
}

//赋值下一个继任者
func (rc *RequestChain) SetSuccessor(m *RequestChain) {
	rc.successor = m
}

// 3. 项目经理
type ProjectManager struct{}

func (pm *ProjectManager) HaveRight(money int) bool {
	return money < 500 // 金额在500元以下有权限
}
func (pm *ProjectManager) HandleFeeRequest(name string, money int) bool {
	if name == "浦东新区" {
		fmt.Printf("工程管理拥有权限，%v：%v\n", name, money)
		return true
	}
	fmt.Printf("工程管理没有权限，%v：%v\n", name, money)
	return false
}
func NewProjectManager() *RequestChain {
	return &RequestChain{
		Manager:   &ProjectManager{},
		successor: nil,
	}
}

// 3. 地区经理
type DepManager struct{}

func (dm *DepManager) HaveRight(money int) bool {
	return money < 5000 // 金额在5000元以下有权限
}
func (dm *DepManager) HandleFeeRequest(name string, money int) bool {
	if name == "上海市" {
		fmt.Printf("部门管理授权通过，%v:%v\n", name, money)
		return true
	}
	fmt.Printf("部门管理未授权，%v:%v\n", name, money)
	return false
}
func NewDepManager() *RequestChain {
	return &RequestChain{
		Manager:   &DepManager{},
		successor: nil,
	}
}

// 3. 总经理
type GeneralManager struct{}

func (gm *GeneralManager) HaveRight(money int) bool {
	return true // 至高无上的权限
}
func (gm *GeneralManager) HandleFeeRequest(name string, money int) bool {
	if name == "中央" {
		fmt.Printf("全体管理授权通过，%v:%v\n", name, money)
		return true
	}
	fmt.Printf("全体管理未授权，%v:%v\n", name, money)
	return false
}
func NewGeneralManager() *RequestChain {
	return &RequestChain{
		Manager:   &GeneralManager{},
		successor: nil,
	}
}

// 主体逻辑
func main() {
	c1 := NewProjectManager()
	c2 := NewDepManager()
	c3 := NewGeneralManager()
	//责任传递方向： c1---->c2----->c3
	c1.SetSuccessor(c2)
	c2.SetSuccessor(c3)

	var c Manager = c1
	//c.HandleFeeRequest("浦东新区", 400)
	//c.HandleFeeRequest("上海市", 1400)
	//c.HandleFeeRequest("中央", 10000)

	c.HandleFeeRequest("天津市", 4000)

}
