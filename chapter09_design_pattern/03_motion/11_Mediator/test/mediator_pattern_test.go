package test

import (
	mediator "go_advenced_code/chapter09_design_pattern/03_motion/11_Mediator"
	"testing"
)

func Test_MediatorPattern(t *testing.T) {
	// 设备注册
	center := mediator.DefaultCloudCenter
	light := mediator.NewMockSmartLight(1)
	center.Register(light)

	fnCallAndLog := func(fn func() error) {
		e := fn()
		if e != nil {
			t.Log(e)
		}
	}

	// 创建app
	app := mediator.NewMockPhoneApp(mediator.DefaultCloudMediator)

	// 设备控制测试
	fnCallAndLog(func() error {
		return app.LightOpen(1)
	})
	fnCallAndLog(func() error {
		return app.LightSwitchMode(1, 1)
	})
	fnCallAndLog(func() error {
		return app.LightSwitchMode(1, 2)
	})
	fnCallAndLog(func() error {
		return app.LightClose(1)
	})
}
