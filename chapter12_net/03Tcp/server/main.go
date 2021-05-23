package main

// TCP server端

/*
TCP/IP(Transmission Control Protocol/Internet Protocol) 即传输控制协议/网间协议，是一种面向连接（连接导向）的、可靠的、
	基于字节流的传输层（Transport layer）通信协议，因为是面向连接的协议，数据像水流一样传输，会存在黏包问题
TCP服务端程序的处理流程：
    1.监听端口
    2.接收客户端请求建立链接
    3.创建goroutine处理链接。
 */
import (
	"bufio"
	"fmt"
	"net"
)

func process(conn net.Conn)  {
	defer conn.Close()
	for  {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("收到client端发来的数据：", recvStr)
		conn.Write([]byte("回复"+recvStr)) // 发送数据
	}

}

func main()  {
	listen,err := net.Listen("tcp","127.0.0.1:20000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for  {
		conn,err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn) // 启动一个goroutine处理连接
	}
}
