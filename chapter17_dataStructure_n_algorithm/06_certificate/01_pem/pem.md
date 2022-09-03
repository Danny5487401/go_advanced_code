# encoding/pem包(Privacy Enhance Mail)

pem包实现了PEM数据编码（保密增强邮件协议）。目前PEM编码主要用于TLS密钥和证书.


## PEM 编码格式如下
```text
-----BEGIN Type-----
Headers
base64-encoded Bytes
-----END Type-----
```

### 1. 证书文件
```text
-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----
```

### 2. 私钥文件：
```text
-----BEGIN RSA PRIVATE KEY-----
-----END RSA PRIVATE KEY-----
```

### 3. 请求文件：

```text
-----BEGIN CERTIFICATE REQUEST-----
-----END CERTIFICATE REQUEST----- 
```


##  Go 源码

```go
type Block struct {
    Type    string            // 得自前言的类型（如"RSA PRIVATE KEY"）
    Headers map[string]string // 可选的头项
    Bytes   []byte            // 内容解码后的数据，一般是DER编码的ASN.1结构
}
```

```go
var pemStart = []byte("\n-----BEGIN ")
var pemEnd = []byte("\n-----END ")
var pemEndOfLine = []byte("-----")


// /usr/local/go/src/encoding/pem/pem.go
func Decode(data []byte) (p *Block, rest []byte)
```
Decode函数会从输入里查找到下一个PEM格式的块（证书、私钥等）。
它返回解码得到的Block和剩余未解码的数据。如果未发现PEM数据，返回(nil, data)

