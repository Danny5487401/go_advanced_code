package observer

import (
	"sync"
	"sync/atomic"
	"time"
)

var GlobalTimeService = NewMockTimeService(1800)

// // 订阅者的数据结构
// 虚拟的时间服务, 自定义时间倍率以方便时钟相关的测试
type tMockTimeService struct {
	observers map[string]ITimeObserver //// 订阅者列表，这里是一个map结构，管理观察者
	rwmutex   *sync.RWMutex
	speed     int64
	state     int64
}

func NewMockTimeService(speed int64) ITimeService {
	it := &tMockTimeService{
		observers: make(map[string]ITimeObserver, 0),
		rwmutex:   new(sync.RWMutex),
		speed:     speed,
		state:     0,
	}
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

			me.NotifyAll(&t)
		}
	}()
}

func (me *tMockTimeService) NotifyAll(now *time.Time) {
	me.rwmutex.RLock()
	defer me.rwmutex.RUnlock()

	for _, it := range me.observers {
		go it.TimeElapsed(now)
	}
}

//订阅方法
func (me *tMockTimeService) Attach(it ITimeObserver) {
	me.rwmutex.Lock()
	defer me.rwmutex.Unlock()

	me.observers[it.ID()] = it
}

func (me *tMockTimeService) Detach(id string) {
	me.rwmutex.Lock()
	defer me.rwmutex.Unlock()

	delete(me.observers, id) // 将订阅者从列表中删除
}
