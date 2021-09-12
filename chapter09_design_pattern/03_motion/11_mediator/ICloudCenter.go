package mediator

// 云中心面向智能设备的注册接口
type ICloudCenter interface {
	Register(dev ISmartDevice) //注册设备
}
