package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"time"
)

// 配置的过程中有参数限制

type Server struct {
	Addr     string
	Port     int
	Protocol string
	Timeout  time.Duration
	MaxConns int
	TLS      *tls.Config
}

// 改造一下函数Option
// 返回错误
type OptionWithError func(*Server) error

func Protocol(p string) OptionWithError {
	return func(s *Server) error {
		if p == "" {
			return errors.New("empty protocol")
		}
		s.Protocol = p
		return nil
	}
}

func Timeout(timeout time.Duration) OptionWithError {
	return func(s *Server) error {
		if timeout.Seconds() < 1 {
			return errors.New("time out should not less than 1s")
		}
		s.Timeout = timeout
		return nil
	}
}
func NewServer(addr string, port int, options ...OptionWithError) (*Server, error) {
	srv := Server{
		Addr:     addr,
		Port:     port,
		Protocol: "tcp",
		Timeout:  30 * time.Second,
		MaxConns: 1000,
		TLS:      nil,
	}
	// 增加了一个参数验证的步骤
	for _, option := range options {
		if err := option(&srv); err != nil {
			return nil, err
		}
	}
	return &srv, nil
}

func main() {
	s1, _ := NewServer("localhost", 1024)
	s2, _ := NewServer("localhost", 2048, Protocol("udp"))
	s3, _ := NewServer("0.0.0.0", 8080, Timeout(300*time.Second), MaxConns(1000))

	fmt.Println(s1, s2, s3)

}
