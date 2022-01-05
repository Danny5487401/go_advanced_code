package socket

import (
	"fmt"
	"net"
	"syscall"
)

// 监听一个 Socket
func Listen(ip string, port int) (*Socket, error) {
	socket := &Socket{}

	// 地址类型：我们用的是 AF_INET (IPv4),
	//socket 类型：我们用 SOCKET_STREAM，代表基于字节流连续、可靠的双向连接,
	//协议类型：0 在 SOCKET_STREAM 类型下代表的是 TCP。
	socketFileDescriptor, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create socket (%v)", err)
	}

	socket.FileDescriptor = socketFileDescriptor

	socketAddress := &syscall.SockaddrInet4{Port: port}
	copy(socketAddress.Addr[:], net.ParseIP(ip))

	if err = syscall.Bind(socket.FileDescriptor, socketAddress); err != nil {
		return nil, fmt.Errorf("failed to bind socket (%v)", err)
	}

	// 监听的第二个参数是连接请求队列的最大长度。我们使用了内核参数 SOMAXCONN,可以使用 sysctl kern.ipc.somaxconn获取
	if err = syscall.Listen(socket.FileDescriptor, syscall.SOMAXCONN); err != nil {
		return nil, fmt.Errorf("failed to listen on socket (%v)", err)
	}

	return socket, nil
}
