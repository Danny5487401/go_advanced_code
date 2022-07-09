package main

import (
	"bufio"
	"github.com/Danny5487401/go_advanced_code/chapter07_system_call/01_kqueue_server/kqueue"
	"github.com/Danny5487401/go_advanced_code/chapter07_system_call/01_kqueue_server/socket"

	"log"
	"os"
	"strings"
)

func main() {
	s, err := socket.Listen("127.0.0.1", 8080)
	if err != nil {
		log.Println("Failed to create Socket:", err)
		os.Exit(1)
	}

	eventLoop, err := kqueue.NewEventLoop(s)
	if err != nil {
		log.Println("Failed to create event loop:", err)
		os.Exit(1)
	}

	log.Println("Server started. Waiting for incoming connections. ^C to exit.")

	eventLoop.Handle(func(s *socket.Socket) {
		reader := bufio.NewReader(s)
		for {
			line, err := reader.ReadString('\n')
			if err != nil || strings.TrimSpace(line) == "" {
				break
			}
			s.Write([]byte(line))
		}
		s.Close()
	})
}

// 我们使用单进程和阻塞 socket 运行，另外，也没有去处理错误。其实大多数情况下，使用已经存在的库而不是自己调用操作系统内核函数会更好
/// curl localhost:8080
