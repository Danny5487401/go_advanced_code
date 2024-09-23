package main

// TCP 粘包问题

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func process(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	var buf [1024]byte
	for {
		n, err := reader.Read(buf[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("收到client发来的数据：", recvStr)

	}

}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer listen.Close()

	for {
		Conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(Conn)
	}

}

/*
收到client发来的数据： Hello, Hello. How are you?Hello, Hello. How are you?Hello, Hello. How are you?
*/
