package observer

// 闹铃的实现类, 实现ITimeObserver接口以订阅时间变化通知
import (
	"fmt"
	"sync/atomic"
	"time"
)

type AlarmClock struct {
	id         string
	name       string
	hour       time.Duration
	minute     time.Duration
	repeatable bool
	next       *time.Time
	occurs     int
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
	GlobalTimeService.Attach(it)

	return it
}

func (me *AlarmClock) NextAlarmTime() *time.Time {
	now := time.Now()
	today, _ := time.ParseInLocation("2006-01-02 15:04:05",
		fmt.Sprintf("%s 00:00:00", now.Format("2006-01-02")), time.Local)
	t := today.Add(me.hour * time.Hour).Add(me.minute * time.Minute)
	if t.Unix() < now.Unix() {
		t = t.Add(24 * time.Hour)
	}
	fmt.Printf("%s.next = %s\n", me.name, t.Format("2006-01-02 15:04:05"))
	return &t
}
func (me *AlarmClock) ID() string {
	return me.name
}

func (me *AlarmClock) TimeElapsed(now *time.Time) {
	it := me.next
	if it == nil {
		return
	}

	if now.Unix() >= it.Unix() {
		me.occurs++
		fmt.Printf("%s 时间=%s 闹铃 %s\n", time.Now().Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"), me.name)

		if me.repeatable {
			t := me.next.Add(24 * time.Hour)
			me.next = &t

		} else {
			GlobalTimeService.Detach(me.ID())
		}
	}
}

func (me *AlarmClock) Occurs() int {
	return me.occurs
}
