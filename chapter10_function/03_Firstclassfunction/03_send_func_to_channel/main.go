package _3_send_func_to_channel

import "fmt"

// Peer代表一个节点
type Peer struct {
	ID string
	// peer其他信息，比如网络连接，地址，协议等等
}

func (p *Peer) WriteMsg(msg string) {
	fmt.Printf("send to:%v,msg:%v\n", p.ID, msg)

}

type Operation func(peers map[string]*Peer)
type Host struct {
	opCh chan Operation
	stop chan struct{}
}

func NewHost() *Host {
	h := &Host{
		opCh: make(chan Operation),
		stop: make(chan struct{}),
	}
	return h

}
func (h *Host) loop() {
	peers := make(map[string]*Peer)
	for {
		select {
		case op := <-h.opCh:
			op(peers)
		case <-h.stop:
			return
		}
	}
}
func (h *Host) AddPeer(p *Peer) {
	add := func(peers map[string]*Peer) {
		peers[p.ID] = p
	}
	h.opCh <- add
}

func (h *Host) Remove(pid string) {
	rm := func(peers map[string]*Peer) {
		delete(peers, pid)
	}
	h.opCh <- rm
}
func (h *Host) Broadcast(msg string) {
	broadcast := func(peers map[string]*Peer) {
		for _, p := range peers {
			p.WriteMsg(msg)
		}
	}
	h.opCh <- broadcast
}

// 获取
func (h *Host) GetPeer(pid string) *Peer {
	// retCh用于接收查询
	retch := make(chan *Peer)
	query := func(peers map[string]*Peer) {
		retch <- peers[pid]
	}
	// 发送查询
	go func() {
		h.opCh <- query
	}()

	// 等待结果返回
	return <-retch
}

// 只向一个Peer发送消息
func (h *Host) SendTo(pid, msg string) {
	p := h.GetPeer(pid)
	p.WriteMsg(msg)
}

/*
总结：
	版本一：使用锁共享数据，
	版本二：使用channel传输数据，没办法实现GetPeer
	版本三：使用channel传输函数,不能对Peer进行同时读，不容易做单元测试

这3种方式本身并无优劣之分，具体要用那种实现，要依赖自身的实际场景进行取舍
*/
