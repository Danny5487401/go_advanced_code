<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Client](#client)
  - [源码](#%E6%BA%90%E7%A0%81)
  - [发送请求流程](#%E5%8F%91%E9%80%81%E8%AF%B7%E6%B1%82%E6%B5%81%E7%A8%8B)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Client
Go HTTP 客户端采用分层设计，主要分为高层 API 和底层实现两部分
- 高层 API：由 Client 结构体提供，负责处理重定向、Cookie、认证等高级功能，为开发者提供简单易用的接口
- 底层实现：由 Transport 结构体提供，负责与 TCP 层交互，管理连接池，发送 HTTP 请求和接收响应

http.Client 表示一个http client端，用来处理HTTP相关的工作，例如cookies, redirect, timeout等工作，

## 源码
Client 结构体

```go
// go1.24.3/src/net/http/client.go
type Client struct { 
    Transport RoundTripper  // Transport指定了如何发送HTTP请求的机制 .如果为nil，则使用DefaultTransport
    CheckRedirect func(req *Request, via []*Request) error  // 指定处理重定向的策略.如果为nil，则使用默认策略（最多跟随10个重定向）
    Jar CookieJar  // Jar指定cookie管理器 .如果为nil，则只有在请求中明确设置的cookie才会被发送
    Timeout time.Duration // 指定客户端请求的最大超时时间，该超时时间包括连接、任何的重定向以及读取相应的时间. 值为0表示没有超时限制
}
```

关键特性:

- 连接复用：通过 Transport 复用 TCP 连接，提高性能
- 自动重定向：处理 3xx 重定向响应，可配置重定向策略
- Cookie 管理：可选的 Cookie 管理功能
- 超时控制：支持请求级别的超时控制

## 发送请求流程

1. 调用 net/http.NewRequest 根据方法名、URL 和请求体构建请求
2. 调用 net/http.Transport.RoundTrip 开启 HTTP 事务、获取连接并发送请求；
3. 在 HTTP 持久连接的 net/http.persistConn.readLoop 方法中等待响应；


![](.http_client_images/client_send_process.png)
```go
func (c *Client) Get(url string) (resp *Response, err error) {
	req, err := NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Do(req *Request) (*Response, error) {
    return c.do(req)
}
```

```go
func (c *Client) do(req *Request) (retres *Response, reterr error) {
    // ...
	for {
		// For all but the first request, create the next
		// request hop and replace req.
		if len(reqs) > 0 {
            // ...
			ireq := reqs[0]
			req = &Request{
				Method:   redirectMethod,
				Response: resp,
				URL:      u,
				Header:   make(Header),
				Host:     host,
				Cancel:   ireq.Cancel,
				ctx:      ireq.ctx,
			}
            // ...
		}

		reqs = append(reqs, req)
		var err error
		var didTimeout func() bool
		// 发送
		if resp, didTimeout, err = c.send(req, deadline); err != nil {
			// c.send() always closes req.Body
			reqBodyClosed = true
			if !deadline.IsZero() && didTimeout() {
				err = &httpError{
					// TODO: early in cycle: s/Client.Timeout exceeded/timeout or context cancellation/
					err:     err.Error() + " (Client.Timeout exceeded while awaiting headers)",
					timeout: true,
				}
			}
			return nil, uerr(err)
		}
		// 。。。 
		req.closeBody()
	}
}
```

```go
// didTimeout is non-nil only if err != nil.
func (c *Client) send(req *Request, deadline time.Time) (resp *Response, didTimeout func() bool, err error) {
	if c.Jar != nil {
		for _, cookie := range c.Jar.Cookies(req.URL) {
			req.AddCookie(cookie)
		}
	}
	// 调用transport的RoundTrip
	resp, didTimeout, err = send(req, c.transport(), deadline)
	if err != nil {
		return nil, didTimeout, err
	}
	if c.Jar != nil {
		if rc := resp.Cookies(); len(rc) > 0 {
			c.Jar.SetCookies(req.URL, rc)
		}
	}
	return resp, nil, nil
}
```



## 参考

- [golang http客户端源码解析](https://zhuanlan.zhihu.com/p/1923369392357498885)