# Client

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