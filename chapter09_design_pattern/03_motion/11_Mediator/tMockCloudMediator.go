package mediator

import "sync"

type tMockCloudMediator struct {
	mDevices map[int]ISmartDevice // 获取不同的设备
	mRWMutex *sync.RWMutex
}

func newMockCloudMediator() ICloudMediator {
	return &tMockCloudMediator{
		make(map[int]ISmartDevice),
		new(sync.RWMutex),
	}
}

func (me *tMockCloudMediator) Register(it ISmartDevice) {
	me.mRWMutex.Lock()
	defer me.mRWMutex.Unlock()

	me.mDevices[it.ID()] = it
}

func (me *tMockCloudMediator) Command(id int, cmd string) string {
	me.mRWMutex.RLock()
	defer me.mRWMutex.RUnlock()

	it, ok := me.mDevices[id]
	if !ok {
		return "device not found"
	}
	return it.Command(cmd)
}

var DefaultCloudMediator = newMockCloudMediator()
var DefaultCloudCenter = DefaultCloudMediator.(ICloudCenter)
