package kqueue

import (
	"fmt"
	"github.com/Danny5487401/go_advanced_code/chapter07_system_call/08_kqueue_server/socket"

	"syscall"
)

type EventLoop struct {
	KqueueFileDescriptor int
	SocketFileDescriptor int
}
type Handler func(s *socket.Socket)

func NewEventLoop(s *socket.Socket) (*EventLoop, error) {
	// 订阅新事件和轮询队列
	kQueue, err := syscall.Kqueue()
	if err != nil {
		return nil,
			fmt.Errorf("failed to create kqueue file descriptor (%v)", err)
	}

	changeEvent := syscall.Kevent_t{
		Ident:  uint64(s.FileDescriptor),           // Ident 的文件描述符：值是我们 socket 的文件描述
		Filter: syscall.EVFILT_READ,                //处理事件的 Filter：设置为 EVFILT_READ，当和监听 socket 一起用时，它代表我们只关心传入连接的事件
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE, // 我们想要添加（EV_ADD）事件到 kqueue，比如说订阅事件，同时要启用（EV_ENABLE）它.
		Fflags: 0,
		Data:   0,
		Udata:  nil,
	}

	changeEventRegistered, err := syscall.Kevent(
		kQueue,
		[]syscall.Kevent_t{changeEvent},
		nil,
		nil,
	)
	if err != nil || changeEventRegistered == -1 {
		return nil,
			fmt.Errorf("failed to register change event (%v)", err)
	}

	return &EventLoop{
		KqueueFileDescriptor: kQueue,
		SocketFileDescriptor: s.FileDescriptor,
	}, nil
}

// Handle 处理函数
func (eventLoop *EventLoop) Handle(handler Handler) {
	for {
		newEvents := make([]syscall.Kevent_t, 10)
		numNewEvents, err := syscall.Kevent(
			eventLoop.KqueueFileDescriptor,
			nil,
			newEvents,
			nil,
		)
		if err != nil {
			continue
		}

		for i := 0; i < numNewEvents; i++ {
			currentEvent := newEvents[i]
			eventFileDescriptor := int(currentEvent.Ident)

			if currentEvent.Flags&syscall.EV_EOF != 0 {
				// 我们要处理 EV_EOF 事件
				// client closing connection
				syscall.Close(eventFileDescriptor)
			} else if eventFileDescriptor == eventLoop.SocketFileDescriptor {
				// 监听 socket 有连接请求
				// new incoming connection
				socketConnection, _, err := syscall.Accept(eventFileDescriptor)
				// 从 TCP 连接请求队列中获取连接请求，它会为监听 socket 创建一个新的客户端 socket 和新的文件描述符。
				if err != nil {
					continue
				}

				socketEvent := syscall.Kevent_t{
					Ident:  uint64(socketConnection),
					Filter: syscall.EVFILT_READ, // 这个新创建的 socket 订阅一个新的 EVFILT_READ 事件
					Flags:  syscall.EV_ADD,
					Fflags: 0,
					Data:   0,
					Udata:  nil,
				}
				socketEventRegistered, err := syscall.Kevent(
					eventLoop.KqueueFileDescriptor,
					[]syscall.Kevent_t{socketEvent},
					nil,
					nil,
				)
				if err != nil || socketEventRegistered == -1 {
					continue
				}
			} else if currentEvent.Filter&syscall.EVFILT_READ != 0 {
				// data available -> forward to handler
				handler(&socket.Socket{
					FileDescriptor: int(eventFileDescriptor),
				})
			}

			// ignore all other events
		}
	}
}
