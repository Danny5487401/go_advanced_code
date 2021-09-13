package mediator

// 智能设备接口
type ISmartDevice interface {
	ID() int
	Command(cmd string) string
}
