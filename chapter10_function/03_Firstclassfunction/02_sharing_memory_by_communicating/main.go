package _2_sharing_memory_by_communicating

import "fmt"

/*
在Go中，推荐使用CSP实现并发，而不是习惯性的使用Lock，使用channel传递数据，达到多goroutine间共享数据的目的，也就是share memory by communicating。
见sharing_memory_by_communicating.png：
	把数据流动画出来，并把要流动的数据标上，然后那些数据流动的线条，就是channel，线条上的数据就是channel要传递的数据，图中也把这些线条和数据标上了

注意：
	数据是双向还是单向
*/

// Peer代表一个节点
type Peer struct {
	ID string
	// peer其他信息，比如网络连接，地址，协议等等
}

func (p *Peer) WriteMsg(msg string) {
	fmt.Printf("send to:%v,msg:%v\n", p.ID, msg)

}

type Host struct {
	add       chan *Peer
	broadcast chan string
	remove    chan string
	stop      chan struct{}

	info chan Info
}

func NewHost() *Host {
	h := &Host{
		add:       make(chan *Peer),
		broadcast: make(chan string),
		remove:    make(chan string),
		stop:      make(chan struct{}),
	}
	return h

}

// 启停Host
func (h *Host) start() {
	go h.loop()
}
func (h *Host) top() {
	close(h.stop)
}

/*
Host增加了2个方法：

	Start()用于启动1个goroutine运行loop()，loop保存所有的peers。
	Stop()用于关闭Host，让loop退出。
*/

// 通信
func (h *Host) loop() {
	peers := make(map[string]*Peer)

	for {
		select {
		case p := <-h.add:
			peers[p.ID] = p
		case pid := <-h.remove:
			delete(peers, pid)
		case msg := <-h.broadcast:
			for _, p := range peers {
				p.WriteMsg(msg)
			}
		case <-h.stop:
			return
		}
	}
}

// 具体方法
func (h *Host) Add(p *Peer) {
	h.add <- p
}
func (h *Host) RemovePeer(pid string) {
	h.remove <- pid
}
func (h *Host) Broadcast(msg string) {
	h.broadcast <- msg
}

type Info struct {
	pid string
	Peer
}

/*
问题：GetPeer是双向的
	有很多goroutine调用GetPeer，我们需要向每一个goroutine发送结果，这就需要每一个goroutine都需要对应的1个接收结果的channel。

做法：
	可以增加1个query channel，channel里传递Peer.ID和接收结果的channel


*/
