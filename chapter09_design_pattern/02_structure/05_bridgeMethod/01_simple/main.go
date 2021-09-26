package main

import "fmt"

// 一.定义Implementor实现类接口及具体实现类

// 发送信息的具体实现（操作）
type MessageImplementer interface {
	send(test, to string)
}

// 1.发送SMS
type MessageSMS struct{}

func (*MessageSMS) send(test, to string) {
	fmt.Printf("SMS信息：[%v]；发送到：[%v]\n", test, to)
}
func ViaSMS() *MessageSMS { //创建
	return &MessageSMS{}
}

// 2.发送email
type MessageEmail struct{}

func (*MessageEmail) send(test, to string) {
	fmt.Printf("Email信息：[%v]；发送到：[%v]\n", test, to)
}
func ViaEmail() *MessageEmail { //创建
	return &MessageEmail{}
}

// 二. 定义Abstraction抽象类及RefinedAbstraction扩充抽象类

//发送信息的二次封装（发送普通信息、紧急信息）
type AbstractMessage interface {
	SendMessage(text, to string)
}

// 1。普通信息
type CommonMessage struct {
	method MessageImplementer
}

func (c *CommonMessage) SendMessage(text, to string) {
	c.method.send(text, to)
}
func NewCommonMessage(mi MessageImplementer) *CommonMessage {
	return &CommonMessage{method: mi}
}

// 2.紧急信息
type UrgencyMessage struct {
	method MessageImplementer
}

func (u *UrgencyMessage) SendMessage(text, to string) {
	u.method.send(fmt.Sprintf("【紧急信息】%v", text), to)
}
func NewUrgencyMessage(mi MessageImplementer) *UrgencyMessage {
	return &UrgencyMessage{method: mi}
}

// 三. 定义工厂方法生产具体的消息

func NewMsgSender(sendWay MessageImplementer, MessageType string) AbstractMessage {
	switch MessageType {
	case "common":
		return NewCommonMessage(sendWay)
	case "urgency":
		return NewUrgencyMessage(sendWay)
	default:
		return nil
	}
}

// 主体逻辑
func main() {
	//var m AbstractMessage
	//返回值为：*CommonMessage类型，而CommonMessage实现了AbstractMessage，因此需要定义的类型为；接口AbstractMessage

	//m = NewCommonMessage(ViaSMS())
	//m.SendMessage("你需要喝一杯咖啡吗", "美女")
	//
	//m = NewCommonMessage(ViaEmail())
	//m.SendMessage("好滴", "帅哥")
	//
	//m = NewUrgencyMessage(ViaSMS())
	//m.SendMessage("晚上桥头见", "美女")
	//
	//m = NewUrgencyMessage(ViaEmail())
	//m.SendMessage("不见不散", "帅哥")

	// 1.邮件
	sms := ViaSMS()
	msgSender := NewMsgSender(sms, "common")
	msgSender.SendMessage("你需要喝一杯咖啡吗", "美女")

}
