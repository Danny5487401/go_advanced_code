package observer

// 闹铃的实现类, 实现ITimeObserver接口以订阅时间变化通知
import (
	"fmt"
	"sync/atomic"
	"time"
)

type AlarmClock struct {
	id         string
	name       string        //响铃的名称
	hour       time.Duration //具体时间
	minute     time.Duration //具体时间
	repeatable bool
	next       *time.Time //下次响铃时间
	occurs     int        //响铃次数
}

var gClockID int64 = 0

func newClockID() string {
	id := atomic.AddInt64(&gClockID, 1)
	return fmt.Sprintf("AlarmClock-%d", id)
}

func NewAlarmClock(name string, hour int, minute int, repeatable bool) *AlarmClock {
	it := &AlarmClock{
		id:         newClockID(),
		name:       name,
		hour:       time.Duration(hour),
		minute:     time.Duration(minute),
		repeatable: repeatable,
		next:       nil,
		occurs:     0,
	}
	it.next = it.NextAlarmTime()

	// 注册这个闹钟
	GlobalTimeService.Attach(it)

	return it
}

func (ac *AlarmClock) NextAlarmTime() *time.Time {
	now := time.Now()
	today, _ := time.ParseInLocation("2006-01-02 15:04:05",
		fmt.Sprintf("%s 00:00:00", now.Format("2006-01-02")), time.Local)
	t := today.Add(ac.hour * time.Hour).Add(ac.minute * time.Minute)
	if t.Unix() < now.Unix() {
		// 代表过了一天
		t = t.Add(24 * time.Hour)
	}
	fmt.Printf("%s.next = %s\n", ac.name, t.Format("2006-01-02 15:04:05"))
	return &t
}
func (ac *AlarmClock) ID() string {
	return ac.name
}

func (ac *AlarmClock) TimeElapsed(now *time.Time) {
	it := ac.next
	if it == nil {
		return
	}

	if now.Unix() >= it.Unix() {
		// 时间过了就发生次数加一
		ac.occurs++
		fmt.Printf("%s 时间=%s 闹铃 %s\n", time.Now().Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"), ac.name)

		if ac.repeatable {
			t := ac.next.Add(24 * time.Hour)
			ac.next = &t

		} else {
			// 不允许多次，开始注销服务
			GlobalTimeService.Detach(ac.ID())
		}
	}
}

func (ac *AlarmClock) Occurs() int {
	return ac.occurs
}
