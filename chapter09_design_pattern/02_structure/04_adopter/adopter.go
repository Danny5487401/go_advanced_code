package _4_adopter

// TypeC TypeC接口:需要被适配的接口
type TypeC interface {
	UseTypeC() string
}

// USB接口
type USB interface {
	UseUSB() string
}

func NewUSB() USB {
	return &keyboard{}
}

type keyboard struct {
}

// keyboard实现了USB接口
func (k *keyboard) UseUSB() string {
	return "I user USB interface"
}

// 适配器
type adapter struct {
	USB // 组合内嵌
}

// UseTypeC实现了Type-C接口
func (a *adapter) UseTypeC() string {
	return a.UseUSB() + ", but now I user Type-C interface"
}

// NewAdapter 是适配器的工厂函数
func NewAdapter(keyboard USB) TypeC {
	return &adapter{
		USB: keyboard,
	}
}
