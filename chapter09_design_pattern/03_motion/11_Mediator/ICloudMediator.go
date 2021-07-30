package mediator

// 云中心面向手机app的接口
type ICloudMediator interface {
	Command(id int, cmd string) string
}
