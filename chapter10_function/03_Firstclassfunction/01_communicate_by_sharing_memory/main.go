package _1_communicate_by_sharing_memory

import (
	"fmt"
	"sync"
)

// 版本一

// Peer代表一个节点
type Peer struct {
	ID string
	// peer其他信息，比如网络连接，地址，协议等等
}

func (p *Peer) WriteMsg(msg string) {
	fmt.Printf("send to:%v,msg:%v\n", p.ID, msg)
}

/*
Peer定义:
	Peer中保存了ID，我们可以通过ID来表示全网中所有的节点，Peer中还有其他字段，比如网络连接、地址、协议版本等信息，此处已经省略掉。

	Peer有一个WriteMsg的方法，实现向该Peer发送消息的功能，例子中使用打印替代。
*/

// Host 代表当前节点的连接管理
type Host struct {
	peers map[string]*Peer // 连接上所有的Peer,根据Peer.ID访问
	lock  sync.RWMutex     //保护peers互斥访问
	// 其他字段省略
}

// NewHost()用来创建一个Host对象，用来代表当前节点。
func NewHost() *Host {
	h := &Host{
		peers: make(map[string]*Peer),
	}
	return h
}

/*
Host有4个方法，分别是：

	1. AddPeer: 增加1个Peer。
	2. RemovePeer: 删除1个Peer。
	3. GetPeer: 通过Peer.ID查询1个Peer。
	4. BroadcastMsg: 向所有Peer发送消息。

每一个方法都需要获取lock，然后访问peers，如果只读取peers则使用读锁
*/

// 添加
func (h *Host) AddPeer(p *Peer) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.peers[p.ID] = p
}

// 获取
func (h *Host) GePeer(pid string) *Peer {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.peers[pid]
}

// 删除
func (h *Host) Remove(pid string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	delete(h.peers, pid)
}

// 广播消息
func (h *Host) BroadcastMsg(msg string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	for _, p := range h.peers {
		p.WriteMsg(msg)
	}

}

/*
communicate by sharing memory

每个goroutine都是1个实体，它们同时运行，调用Host的不同方法来访问peers，只有拿到当前lock的goroutine才能访问peers，仿佛当前goroutine在同其他goroutine讲：
	我现在有访问权，你们等一下。本质上就是，通过共享Host.lock这块内存，各goroutine进行交流（表明自己拥有访问权）。
*/
