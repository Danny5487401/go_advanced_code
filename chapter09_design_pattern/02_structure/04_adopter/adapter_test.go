package _4_adopter

import "testing"

func TestAdapter(t *testing.T) {
	// 生成需要被适配的接口
	usb := NewUSB()
	// 生成适配器
	adapter := NewAdapter(usb)
	t.Log(usb.UseUSB())
	t.Log(adapter.UseTypeC())
}
