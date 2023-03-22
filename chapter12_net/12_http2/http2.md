# http2


HTTP/2协议定义了两个版本：h2 和 h2c：

- h2版本的协议是建立在TLS层之上的HTTP/2协议，这个标志被用在TLS应用层协议协商（TLS-ALPN）域和任何其它的TLS之上的HTTP/2协议。
- h2c版本是建立在明文的TCP之上的HTTP/2协议，这个标志被用在HTTP/1.1的升级协议头域和其它任何直接在TCP层之上的HTTP/2协议





## 参考资料
1. [http2 官方文档](https://httpwg.org/specs/rfc7540.html)
2. [go官方x/http2 包](https://pkg.go.dev/golang.org/x/net/http2)