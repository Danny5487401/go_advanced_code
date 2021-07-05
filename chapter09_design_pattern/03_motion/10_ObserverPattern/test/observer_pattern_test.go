package test

import (
	observer "go_advenced_code/chapter09_design_pattern/03_motion/10_ObserverPattern"
	"testing"
	"time"
)

func Test_ObserverPattern(t *testing.T) {
	_ = observer.NewAlarmClock("下午开会", 14, 30, false)

	_ = observer.NewAlarmClock("起床", 6, 0, true)
	_ = observer.NewAlarmClock("午饭", 12, 30, true)
	_ = observer.NewAlarmClock("午休", 13, 0, true)
	_ = observer.NewAlarmClock("晚饭", 18, 30, true)
	clock := observer.NewAlarmClock("晚安", 22, 0, true)

	for {
		if clock.Occurs() >= 2 {
			break
		}
		time.Sleep(time.Second)
	}
}
