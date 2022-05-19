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

## 证书种类
证书分为根证书、服务器证书、客户端证书。根证书文件（ca.crt）和根证书对应的私钥文件（ca.key）由 CA（证书授权中心，国际认可）生成和保管。那么服务器如何获得证书呢？向 CA 申请！步骤如下：

1. 服务器生成自己的公钥（server.pub）和私钥（server.key）。后续通信过程中，客户端使用该公钥加密通信数据，服务端使用对应的私钥解密接收到的客户端的数据；

2. 服务器使用公钥生成请求文件（server.req），请求文件中包含服务器的相关信息，比如域名、公钥、组织机构等

3. 服务器将 server.req 发送给 CA。CA 验证服务器合法后，使用 ca.key 和 server.req 生成证书文件（server.crt）——使用私钥生成证书的签名数据；

4. CA 将证书文件（server.crt）发送给服务器。

由于ca.key 和 ca.crt 是一对，ca.crt 文件中包含公钥，因此 ca.crt 可以验证 server.crt是否合法——使用公钥验证证书的签名。


### 自建CA

```shell
[test1280@test1280 https-server]$ openssl genrsa -out rootCA.key 2048
Generating RSA private key, 2048 bit long modulus
....................+++
........................................+++
e is 65537 (0x10001)
[test1280@test1280 https-server]$ openssl req -x509 -new -nodes -key rootCA.key -days 365 -out rootCA.pem
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [XX]:
State or Province Name (full name) []:
Locality Name (eg, city) [Default City]:
Organization Name (eg, company) [Default Company Ltd]:
Organizational Unit Name (eg, section) []:
Common Name (eg, your name or your server's hostname) []:test1280'CA
Email Address []:
```
其中，rootCA.key是我们CA的私钥，rootCA.pem是我们CA的证书


### 签发证书
```shell
[test1280@test1280 https-server]$ openssl genrsa -out server.key 2048
Generating RSA private key, 2048 bit long modulus
.+++
.+++
e is 65537 (0x10001)
[test1280@test1280 https-server]$ openssl req -new -key server.key -out server.csr
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [XX]:
State or Province Name (full name) []:
Locality Name (eg, city) [Default City]:
Organization Name (eg, company) [Default Company Ltd]:
Organizational Unit Name (eg, section) []:
Common Name (eg, your name or your server's hostname) []:test1280
Email Address []:

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:
[test1280@test1280 https-server]$ openssl x509 -req -in server.csr -CA rootCA.pem -CAkey rootCA.key -CAcreateserial -out server.crt -days 365
Signature ok
subject=/C=XX/L=Default City/O=Default Company Ltd/CN=test1280
Getting CA Private Key

```

其中，server.key是服务端私钥，server.crt是由自建CA签名发布的证书。

Note:在生成server.csr(Certificate Signing Request)时，主机名填写的是test1280（后续客户端访问将用到）



## 验证方式

验证方式分为单向验证和双向验证。


### 单向验证

单向验证是指通信双方中一方验证另一方是否合法。通常是指客户端验证服务器。

客户端需要：ca.crt

服务器需要：server.crt，server.key


PS：我们平时使用 PC 上网时使用的就是单向验证的方式。即，我们验证我们要访问的网站的合法性。
PC 中的浏览器（火狐、IE、chrome等）已经包含了很多 CA 的根证书（ca.crt）。当我们访问某个网站（比如：https://www.baidu.com）时，网站会将其证书（server.crt）发送给浏览器，浏览器会使用 ca.crt 验证 server.crt 是否合法。如果发现访问的是不合法网站，浏览器会给出提示

现实中，有的公司会使用自签发证书，即公司自己生成根证书（ca.crt）。如果我们信任此网站，那么需要手动将其证书添加到系统中。


### 双向验证

双向验证是指通信双方需要互相验证对方是否合法。服务器验证客户端，客户端验证服务器。

服务器需要：ca.crt，server.crt，server.key

客户端需要：ca.crt，client.crt，client.key


双向验证通常用于支付系统中，比如支付宝。我们在使用支付宝时必须下载数字证书，该证书就是支付宝颁发给针对我们这台机器的证书，我们只能使用这台机器访问支付宝。如果换了机器，那么需要重新申请证书


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