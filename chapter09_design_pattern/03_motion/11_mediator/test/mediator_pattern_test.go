package test

import (
	mediator "go_advanced_code/chapter09_design_pattern/03_motion/11_Mediator"
	"testing"
)

func Test_MediatorPattern(t *testing.T) {
	// 设备注册中心
	center := mediator.DefaultCloudCenter

	// 创建设备号为1的 智能灯
	light := mediator.NewMockSmartLight(1)

	// 注册设备
	center.Register(light)

	// 创建客户端app
	app := mediator.NewMockPhoneApp(mediator.DefaultCloudMediator)

	// 定义动作
	fnCallAndLog := func(fn func() error) {
		e := fn()
		if e != nil {
			t.Log(e)
		}
	}

	// 开始执行设备控制测试
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
