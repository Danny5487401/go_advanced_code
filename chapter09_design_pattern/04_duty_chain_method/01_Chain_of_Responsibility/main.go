package main

/*
责任链模式：一种处理请求的模式，它让多个处理器都有机会处理该请求，直到其中某个处理成功为止。责任链模式把多个处理器串成链，然后让请求在链上传递

意图：避免请求发送者与接受者耦合在一起，让多个对象都有可能接收请求，将这些对象连接成一条链，并且沿着这条链传递请求，直到有对象处理它为止。

主要解决：职责链上的处理者负责处理请求，客户只需要将请求发送到职责链上即可，无须关心请求的处理细节和请求的传递，所以职责链将请求的发送者和请求的处理者解偶了。

何时使用：在处理消息的时候以过滤很多道。

如何解决：拦截的类都实现统一接口。

优点：

1、降低耦合度。它将请求的发送者和接收者解耦。
2、简化了对象。使得对象不需要知道链的结构。
3、增强给对象指派职责的灵活性。通过改变链内的成员或者调动它们的次序，允许动态地新增或者删除责任。
4、增加新的请求处理类很方便。
缺点：

1、不能保证请求一定被接收。
2、系统性能将受到一定影响，而且在进行代码调试时不太方便，可能会造成循环调用。
3、可能不容易观察运行时的特征，有碍于除错。
使用场景：

1、有多个对象可以处理同一个请求，具体哪个对象处理该请求由运行时刻自动确定。
2、在不明确指定接受者的情况下，向多个对象中的一个提交请求。
3、可动态指定一组对象处理请求。

应用实例： 1、红楼梦中的"击鼓传花"。 2、JS 中的事件冒泡。 3、JAVA WEB 中 Apache Tomcat 对 Encoding 的处理，Struts2 的拦截器，jsp servlet 的 Filter。
 */



import "fmt"
// 1.创建接口
type Manager interface {
	HaveRight(money int)bool //有权做
	HandleFeeRequest(name string,money int)bool //处理费用请求
}

// 2, 创建请求链
type RequestChain struct {
	Manager
	successor *RequestChain //后续者(链式结构，链表)
}

func (rc *RequestChain)HaveRight(money int)bool {
	return true
}
func (rc *RequestChain)HandleFeeRequest(name string,money int)bool {
	//先判断金额在对应的部门是否有权限
	if rc.Manager.HaveRight(money) {
		return rc.Manager.HandleFeeRequest(name,money) //如果有权限，处理请求
	}
	//继任者不为空,传递到继任者处理
	if rc.successor!=nil {
		return rc.successor.HandleFeeRequest(name,money)
	}
	return false
}
//赋值下一个继任者
func (rc *RequestChain)SetSuccessor(m *RequestChain) {
	rc.successor=m
}


// 3. 项目经理
type ProjectManager struct {}

func (pm *ProjectManager)HaveRight(money int)bool {
	return  money<500  // 金额在500元以下有权限
}
func (pm *ProjectManager)HandleFeeRequest(name string,money int)bool {
	if name=="浦东新区" {
		fmt.Printf("工程管理拥有权限，%v：%v\n",name,money)
		return true
	}
	fmt.Printf("工程管理没有权限，%v：%v\n",name,money)
	return false
}
func NewProjectManager() *RequestChain {
	return &RequestChain{
		Manager:   &ProjectManager{},
		successor: nil,
	}
}

// 3. 地区经理
type DepManager struct {}

func (dm *DepManager)HaveRight(money int)bool  {
	return money<5000  // 金额在5000元以下有权限
}
func (dm *DepManager)HandleFeeRequest(name string,money int)bool  {
	if name=="上海市" {
		fmt.Printf("部门管理授权通过，%v:%v\n",name,money)
		return true
	}
	fmt.Printf("部门管理未授权，%v:%v\n",name,money)
	return false
}
func NewDepManager() *RequestChain {
	return &RequestChain{
		Manager:   &DepManager{},
		successor: nil,
	}
}

// 3. 总经理
type GeneralManager struct {}
func (gm *GeneralManager)HaveRight(money int)bool  {
	return true  // 至高无上的权限
}
func (gm *GeneralManager)HandleFeeRequest(name string,money int)bool  {
	if name=="中央" {
		fmt.Printf("全体管理授权通过，%v:%v\n",name,money)
		return true
	}
	fmt.Printf("全体管理未授权，%v:%v\n",name,money)
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