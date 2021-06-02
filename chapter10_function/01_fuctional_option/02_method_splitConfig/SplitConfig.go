package _2_method_splitConfig

import (
	"crypto/tls"
	"fmt"
	"time"
)

// 简单做法：分离变化点和不变点。这里，我们可以将必填项认为是不变点，而非必填则是变化点。

type Config struct {
	Protocol string
	Timeout  time.Duration
	MaxConns int
	TLS      *tls.Config
}

type Server struct {
	Addr string
	Port int
	Conf *Config
}

func NewServer(addr string, port int, conf *Config) (*Server, error) {
	return &Server{
		Addr: addr,
		Port: port,
		Conf: conf,
	}, nil
}

func main() {
	srv1, _ := NewServer("localhost", 9000, nil)

	conf := Config{Protocol: "tcp", Timeout: 60 * time.Second}
	srv2, _ := NewServer("localhost", 9000, &conf)

	fmt.Println(srv1, srv2)
}

// 满足大部分的开发需求
