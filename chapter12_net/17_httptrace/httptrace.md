<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [httptrace](#httptrace)
  - [结构体](#%E7%BB%93%E6%9E%84%E4%BD%93)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# httptrace


Go1.7 引入了HTTP trace，可以在HTTP客户端请求过程中收集一些更细粒度的信息，httptrace包提供了HTTP trace的支持，收集的信息可用于调试延迟问题，服务监控，编写自适应系统等。

httptrace包提供了许多钩子，在HTTP往返期间收集各种事件的信息，包括连接的创建、复用、DNS解析查询、写入请求和读取响应


## 结构体
```go
// go1.24.3/src/net/http/httptrace/trace.go
type ClientTrace struct {
	// 在创建连接之前调用
	GetConn func(hostPort string)

	// 连接成功后调用
	GotConn func(GotConnInfo)

	// 当连接用完时需要放回池子调用
	PutIdleConn func(err error)

	// 当读到响应的第一个字节时
	GotFirstResponseByte func()

	// Got100Continue is called if the server replies with a "100
	// Continue" response.
	Got100Continue func()

	// Got1xxResponse is called for each 1xx informational response header
	// returned before the final non-1xx response. Got1xxResponse is called
	// for "100 Continue" responses, even if Got100Continue is also defined.
	// If it returns an error, the client request is aborted with that error value.
	Got1xxResponse func(code int, header textproto.MIMEHeader) error

	// 当DNS查询开始时调用
	DNSStart func(DNSStartInfo)

	// 当DNS查询结束时调用
	DNSDone func(DNSDoneInfo)

	// 当一个新连接 Dial开始时
	ConnectStart func(network, addr string)

	// 当一个新的连接的Dial完成时
	ConnectDone func(network, addr string, err error)

	// TLSHandshakeStart is called when the TLS handshake is started. When
	// connecting to an HTTPS site via an HTTP proxy, the handshake happens
	// after the CONNECT request is processed by the proxy.
	TLSHandshakeStart func()

	// TLSHandshakeDone is called after the TLS handshake with either the
	// successful handshake's connection state, or a non-nil error on handshake
	// failure.
	TLSHandshakeDone func(tls.ConnectionState, error)

	// WroteHeaderField is called after the Transport has written
	// each request header. At the time of this call the values
	// might be buffered and not yet written to the network.
	WroteHeaderField func(key string, value []string)

	// 在传输时write完所有请求头后调用的
	WroteHeaders func()

	// Wait100Continue is called if the Request specified
	// "Expect: 100-continue" and the Transport has written the
	// request headers but is waiting for "100 Continue" from the
	// server before writing the request body.
	Wait100Continue func()

	// WroteRequest is called with the result of writing the
	// request and any body. It may be called multiple times
	// in the case of retried requests.
	WroteRequest func(WroteRequestInfo)
}
```


使用

```go
func WithClientTrace(ctx context.Context, trace *ClientTrace) context.Context {
	if trace == nil {
		panic("nil trace")
	}
	old := ContextClientTrace(ctx)
	trace.compose(old)

	ctx = context.WithValue(ctx, clientEventContextKey{}, trace)
	if trace.hasNetHooks() { // 有 trace 相关 hook
		nt := &nettrace.Trace{
			ConnectStart: trace.ConnectStart,
			ConnectDone:  trace.ConnectDone,
		}
		if trace.DNSStart != nil {
			nt.DNSStart = func(name string) {
				trace.DNSStart(DNSStartInfo{Host: name})
			}
		}
		if trace.DNSDone != nil {
			nt.DNSDone = func(netIPs []any, coalesced bool, err error) {
				addrs := make([]net.IPAddr, len(netIPs))
				for i, ip := range netIPs {
					addrs[i] = ip.(net.IPAddr)
				}
				trace.DNSDone(DNSDoneInfo{
					Addrs:     addrs,
					Coalesced: coalesced,
					Err:       err,
				})
			}
		}
		ctx = context.WithValue(ctx, nettrace.TraceKey{}, nt)
	}
	return ctx
}

```

## 参考
