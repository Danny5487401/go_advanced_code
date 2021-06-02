package main

import (
	"crypto/tls"
	"fmt"
	"time"
)

type Server struct {
	Addr     string
	Port     int
	Protocol string
	Timeout  time.Duration
	MaxConns int
	TLS      *tls.Config
}

// 定义一个Option类型的函数，它操作了Server这个对象
type Option func(*Server)

// 下面是对四个可选参数的配置函数
func Protocol(p string) Option {
	return func(s *Server) {
		s.Protocol = p
	}
}

func Timeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.Timeout = timeout
	}
}

func MaxConns(maxconns int) Option {
	return func(s *Server) {
		s.MaxConns = maxconns
	}
}

func TLS(tls *tls.Config) Option {
	return func(s *Server) {
		s.TLS = tls
	}
}

// 用到了不定参数的特性，将任意个option应用到Server上
func NewServer(addr string, port int, options ...Option) (*Server, error) {
	// 先填写默认值
	srv := Server{
		Addr:     addr,
		Port:     port,
		Protocol: "tcp",
		Timeout:  30 * time.Second,
		MaxConns: 1000,
		TLS:      nil,
	}
	// 应用任意个option
	for _, option := range options {
		option(&srv)
	}
	return &srv, nil

}
func main() {
	s1, _ := NewServer("localhost", 1024)
	s2, _ := NewServer("localhost", 2048, Protocol("udp"))
	s3, _ := NewServer("0.0.0.0", 8080, Timeout(300*time.Second), MaxConns(1000))

	fmt.Println(s1, s2, s3)

}

/*优点：
1. 可读性强，将配置都转化成了对应的函数项option
2. 扩展性好，新增参数只需要增加一个对应的方法

*/
