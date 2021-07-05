package observer

import "time"

// 定义时间观察者接口, 接收时间变化事件的通知

type ITimeObserver interface {
	ID() string
	TimeElapsed(now *time.Time)
}
