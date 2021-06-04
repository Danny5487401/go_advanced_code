package _1_problem

import (
	"crypto/tls"
	"time"
)

// 需求：常见的HTTP服务器的配置，它区分了2个必填参数与4个非必填参数

type ServerCfg struct {
	Addr     string        // 必填
	Port     int           // 必填
	Protocol string        // 非必填
	Timeout  time.Duration // 非必填
	MaxConns int           // 非必填
	TLS      *tls.Config   // 非必填
}
type Server struct {
}

/*  不好的做法
func NewServer(addr string, port int) (*Server, error)                                   {}
func NewTLSServer(addr string, port int, tls *tls.Config) (*Server, error)               {}
func NewServerWithTimeout(addr string, port int, timeout time.Duration) (*Server, error) {}
func NewTLSServerWithMaxConnAndTimeout(addr string, port int, maxconns int, timeout time.Duration, tls *tls.Config) (*Server, error) {
}
*/

// 问题：我们要实现非常多种方法，来支持各种非必填的情况，示例如下
