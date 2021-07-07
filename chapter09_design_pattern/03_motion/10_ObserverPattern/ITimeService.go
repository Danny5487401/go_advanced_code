package observer

// 定义时间服务的接口, 接受观察者的注册和注销
type ITimeService interface {
	Attach(observer ITimeObserver) //注册
	Detach(id string)              // 注销
}
