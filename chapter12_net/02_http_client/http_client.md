<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Client](#client)
  - [源码](#%E6%BA%90%E7%A0%81)
  - [发送请求流程](#%E5%8F%91%E9%80%81%E8%AF%B7%E6%B1%82%E6%B5%81%E7%A8%8B)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Client

http.Client 表示一个http client端，用来处理HTTP相关的工作，例如cookies, redirect, timeout等工作，

## 源码
Client 结构体

```go
type Client struct { 
    Transport RoundTripper  // 表示 HTTP 事务，用于处理客户端的请求连接并等待服务端的响应；
    CheckRedirect func(req *Request, via []*Request) error  // 用于指定处理重定向的策略
    Jar CookieJar  // 用于管理和存储请求中的 cookie
    Timeout time.Duration // 指定客户端请求的最大超时时间，该超时时间包括连接、任何的重定向以及读取相应的时间
}
```

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

