package mediator

import (
	"fmt"
	"strconv"
	"strings"
)

// 虚拟的智能灯设备, 实现ISmartDevice接口
type tMockSmartLight struct {
	id int
}

func NewMockSmartLight(id int) ISmartDevice {
	return &tMockSmartLight{
		id,
	}
}

// 设备唯一号
func (me *tMockSmartLight) ID() int {
	return me.id
}

// 执行的操作
func (me *tMockSmartLight) Command(cmd string) string {
	if cmd == "light open" {
		e := me.open()
		if e != nil {
			return e.Error()
		}
	} else if cmd == "light close" {
		e := me.close()
		if e != nil {
			return e.Error()
		}
	} else if strings.HasPrefix(cmd, "light switch_mode") {
		args := strings.Split(cmd, " ")
		if len(args) != 3 {
			return "invalid switch command"
		}

		n, e := strconv.Atoi(args[2])
		if e != nil {
			return "invalid mode number"
		}

		e = me.switchMode(n)
		if e != nil {
			return e.Error()
		}

	} else {
		return "unrecognized command"
	}

	return "OK"
}

func (me *tMockSmartLight) open() error {
	fmt.Printf("tMockSmartLight.open, id=%v\n", me.id)
	return nil
}

func (me *tMockSmartLight) close() error {
	fmt.Printf("tMockSmartLight.close, id=%v\n", me.id)
	return nil
}

func (me *tMockSmartLight) switchMode(mode int) error {
	fmt.Printf("tMockSmartLight.switchMode, id=%v, mode=%v\n", me.id, mode)
	return nil
}
