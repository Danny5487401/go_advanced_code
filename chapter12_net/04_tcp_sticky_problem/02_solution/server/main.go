package main

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/Danny5487401/go_advanced_code/chapter12_net/04_tcp_sticky_problem/02_solution/proto"
)

func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		msg, err := proto.Decode(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("decode msg failed, err:", err)
			return
		}
		fmt.Println("收到client发来的数据：", msg)
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
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn)
	}
}

/*

收到client发来的数据： Hello, Hello. How are you?
收到client发来的数据： Hello, Hello. How are you?
收到client发来的数据： Hello, Hello. How are you?

*/
