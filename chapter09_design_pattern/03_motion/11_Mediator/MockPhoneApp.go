package mediator

import (
	"errors"
	"fmt"
)

// 虚拟的手机app, 用于跟云中心通信, 控制智能设备

type MockPhoneApp struct {
	mediator ICloudMediator
}

func NewMockPhoneApp(mediator ICloudMediator) *MockPhoneApp {
	return &MockPhoneApp{
		mediator,
	}
}

func (me *MockPhoneApp) LightOpen(id int) error {
	return me.lightCommand(id, "light open")
}

func (me *MockPhoneApp) LightClose(id int) error {
	return me.lightCommand(id, "light close")
}

func (me *MockPhoneApp) LightSwitchMode(id int, mode int) error {
	return me.lightCommand(id, fmt.Sprintf("light switch_mode %v", mode))
}

func (me *MockPhoneApp) lightCommand(id int, cmd string) error {
	res := me.mediator.Command(id, cmd)
	if res != "OK" {
		return errors.New(res)
	}
	return nil
}
