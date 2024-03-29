package observer

import (
	"sync"
	"sync/atomic"
	"time"
)

var GlobalTimeService = NewMockTimeService(1800)

// 订阅者的数据结构
// 虚拟的时间服务, 自定义时间倍率以方便时钟相关的测试
type tMockTimeService struct {
	observers map[string]ITimeObserver //// 订阅者列表，这里是一个map结构，管理观察者
	rwMutex   *sync.RWMutex
	speed     int64
	state     int64 //标记开始
}

func NewMockTimeService(speed int64) ITimeService {
	it := &tMockTimeService{
		observers: make(map[string]ITimeObserver, 0),
		rwMutex:   new(sync.RWMutex),
		speed:     speed,
		state:     0,
	}
	//开启服务
	it.Start()
	return it
}

func (me *tMockTimeService) Start() {
	if !atomic.CompareAndSwapInt64(&(me.state), 0, 1) {
		return
	}

	go func() {
		timeFrom := time.Now()
		timeOffset := timeFrom.UnixNano()

		for range time.Tick(time.Duration(100) * time.Millisecond) {
			if me.state == 0 {
				break
			}

			nanos := (time.Now().UnixNano() - timeOffset) * me.speed
			t := timeFrom.Add(time.Duration(nanos) * time.Nanosecond)

			// 通知服务
			me.NotifyAll(&t)
		}
	}()
}

func (me *tMockTimeService) NotifyAll(now *time.Time) {
	me.rwMutex.RLock()
	defer me.rwMutex.RUnlock()

	for _, it := range me.observers {
		go it.TimeElapsed(now)
	}
}

//订阅方法
func (me *tMockTimeService) Attach(it ITimeObserver) {
	me.rwMutex.Lock()
	defer me.rwMutex.Unlock()

	me.observers[it.ID()] = it
}

// 删除服务
func (me *tMockTimeService) Detach(id string) {
	me.rwMutex.Lock()
	defer me.rwMutex.Unlock()

	delete(me.observers, id) // 将订阅者从列表中删除
}
