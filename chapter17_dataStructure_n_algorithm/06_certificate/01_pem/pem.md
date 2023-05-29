<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [encoding/pem包(Privacy Enhance Mail)](#encodingpem%E5%8C%85privacy-enhance-mail)
  - [PEM 编码格式如下](#pem-%E7%BC%96%E7%A0%81%E6%A0%BC%E5%BC%8F%E5%A6%82%E4%B8%8B)
    - [1. 证书文件](#1-%E8%AF%81%E4%B9%A6%E6%96%87%E4%BB%B6)
    - [2. 私钥文件：](#2-%E7%A7%81%E9%92%A5%E6%96%87%E4%BB%B6)
    - [3. 请求文件：](#3-%E8%AF%B7%E6%B1%82%E6%96%87%E4%BB%B6)
  - [Go 源码](#go-%E6%BA%90%E7%A0%81)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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

