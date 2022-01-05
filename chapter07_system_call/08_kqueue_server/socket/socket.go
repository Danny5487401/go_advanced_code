package socket

import (
	"strconv"
	"syscall"
)

type Socket struct {
	FileDescriptor int //文件描述符
}

// 1.实现 io.Reader 这个接口
func (socket Socket) Read(bytes []byte) (int, error) {
	if len(bytes) == 0 {
		return 0, nil
	}
	numBytesRead, err := syscall.Read(socket.FileDescriptor, bytes)
	if err != nil {
		numBytesRead = 0
	}

	return numBytesRead, err
}

// 2.实现 io.Writer 接口
func (socket Socket) Write(bytes []byte) (int, error) {
	numBytesWritten, err := syscall.Write(socket.FileDescriptor, bytes)
	if err != nil {
		numBytesWritten = 0
	}
	return numBytesWritten, err
}

// 3.关闭 socket 可以调用 close()[9
func (socket *Socket) Close() error {
	return syscall.Close(socket.FileDescriptor)
}

// 为了稍后能打印一些有用的错误和日志，我们也需要实现 fmt.Stringer 接口。我们通过不同的文件描述符来区分不同的 socket
func (socket *Socket) String() string {
	return strconv.Itoa(socket.FileDescriptor)
}

// 定义事件循环
