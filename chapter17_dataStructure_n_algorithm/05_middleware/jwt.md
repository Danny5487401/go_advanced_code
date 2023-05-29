<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Jwt（JSON Web Token）](#jwtjson-web-token)
  - [jwt构成](#jwt%E6%9E%84%E6%88%90)
    - [头部 Header](#%E5%A4%B4%E9%83%A8-header)
    - [Payload：载荷](#payload%E8%BD%BD%E8%8D%B7)
    - [Signature：签名](#signature%E7%AD%BE%E5%90%8D)
  - [安全](#%E5%AE%89%E5%85%A8)
  - [缺点](#%E7%BC%BA%E7%82%B9)
  - [go-jwt源码分析](#go-jwt%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Jwt（JSON Web Token）
JSON Web Token（JWT）是一个开放标准（RFC 7519），它定义了一种紧凑且自包含的方式，用于在各方之间以JSON方式安全地传输信息。
由于此信息是经过数字签名的，因此可以被验证和信任。可以使用秘密（使用HMAC算法）或使用RSA或ECDSA的公钥/私钥对对JWT进行签名

直白的讲jwt就是一种用户认证（区别于session、cookie）的解决方案。

## jwt构成

JWT 由三部分组成：头部、数据体、签名 / 加密。这三部分以 . (英文句号) 连接，注意这三部分顺序是固定的，即 header.payload.signature 如下示例：

```shell
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

###  头部 Header

这部分用来描述 JWT 的元数据，比如该 JWT 所使用的签名 / 加密算法、媒体类型等。

这部分原始数据是一个 JSON 对象，经过 Base64Url 编码方式进行编码后得到最终的字符串。其中只有一个属性是必要的：alg—— 加密 / 签名算法，默认值为 HS256。

最简单的头部可以表示成这样
```json
{
    "alg": "none"
}
```
其他可选属性：
- typ，描述 JWT 的媒体类型，该属性的值只能是 JWT，它的作用是与其他 JOSE Header 混合时表明自己身份的一个参数（很少用到）。
- cty，描述 JWT 的内容类型。只有当需要一个 Nested JWT 时，才需要该属性，且值必须是 JWT。
- kid，KeyID，用于提示是哪个密钥参与加密。




### Payload：载荷
又称为Claim，携带的信息，比如用户名、过期时间等，一般叫做 Claim.
原始数据仍是一个 JSON 对象，经过 Base64url 编码方式进行编码后得到最终的 Payload。这里的数据默认是不加密的，所以不应存放重要数据（当然你可以考虑使用嵌套型 JWT）。官方内置了七个属性，大小写敏感，且都是可选属性，如下：


- iss (Issuer) 签发人，即签发该 Token 的主体
- sub (Subject) 主题，即描述该 Token 的用途，一般就最为用户的唯一标识
- aud (Audience) 作用域，即描述这个 Token 是给谁用的，多个的情况下该属性值为一个字符串数组，单个则为一个字符串
- exp (Expiration Time) 过期时间，即描述该 Token 在何时失效
- nbf (Not Before) 生效时间，即描述该 Token 在何时生效
- iat (Issued At) 签发时间，即描述该 Token 在何时被签发的
- jti (JWT ID) 唯一标识


除以上默认字段外，我们还可以自定义私有字段，如下例：
```json
{
"sub": "1234567890",
"name": "chongchong",
"admin": true
}

```

这里对 aud 做一个说明，有如下 Payload：
```json
{
    "iss": "server1",
    "aud": ["http://www.a.com","http://www.b.com"]
}
```
那么如果我拿这个 JWT 去 http://www.c.com 获取有访问权限的资源，就会被拒绝掉，因为 aud 属性明确了这个 Token 是无权访问 www.c.com 的，



### Signature：签名
由header、payload 和你自己维护的一个 secret 经过加密得来的
```shell
# 1.签名哈希部分是对上面两部分数据签名，通过指定的算法生成哈希，以确保数据不会被篡改。
# 2.首先，需要指定一个密码（secret），该密码仅仅为保存在服务器中，并且不能向用户公开。
# 3.然后，使用标头中7指定的签名算法（默认情况下为HMAC SHA256）根据以下公式生成签名。
# 4.HMACSHA256(base64UrlEncode(header) + "." + base64UrlEncode(payload),secret)
# 5.在计算出签名哈希后，JWT头，有效载荷和签名哈希的三个部分组合成一个字符串，每个部分用"."分隔，就构成整个JWT对象。
```

## 安全

- 因为 JWT 的前两个部分仅是做了 Base64 编码处理并非加密，所以在存放数据上不能存放敏感数据。
- 用来签名 / 加密的密钥需要妥善保存。
- 尽可能采用 HTTPS，确保不被窃听。
- 如果存放在 Cookie 中则强烈建议开启 Http Only，其实官方推荐是放在 LocalStorage 里，然后通过 Header 头进行传递。

> Cookie 的 HTTP Only 这个 Flag 和 HTTPS 并不冲突，你会发现其实还有一个 Secure 的 Flag，这个就是指 HTTPS 了，这两个 Flag 互不影响的，开启 HTTP Only 会导致前端 JavaScript 无法读取该 Cookie，更多的是为了防止 类 XSS 攻击。


## 缺点
1. 数据臃肿

2. 无法废弃和续签#

3. Token 丢失#

## go-jwt源码分析
标准载荷
```go
type StandardClaims struct {
	Audience  string `json:"aud,omitempty"`  //接收jwt的一方
	ExpiresAt int64  `json:"exp,omitempty"`  //jwt的过期时间，这个过期时间必须要大于签发时间
	Id        string `json:"jti,omitempty"`  //jwt的唯一身份标识，主要用来作为一次性token,从而回避重放攻击。
	IssuedAt  int64  `json:"iat,omitempty"` // jwt的签发时间
	Issuer    string `json:"iss,omitempty"` // jwt签发者
    NotBefore int64  `json:"nbf,omitempty"` // 定义在什么时间之前，该jwt都是不可用的.就是这条token信息生效时间.这个值可以不设置,但是设定后,一定要大于当前Unix UTC,否则token将会延迟生效.
	Subject   string `json:"sub,omitempty"` // jwt所面向的用户
}
```

```go
// token结构
type Token struct {
    Raw       string                 // 保存原始token解析的时候保存
    Method    SigningMethod          // 保存签名方法 目前库里有HMAC  RSA  ECDSA
    Header    map[string]interface{} // jwt中的头部
    Claims    Claims                 // jwt中第二部分荷载，Claims是一个接口
    Signature string                 // jwt中的第三部分 签名
    Valid     bool                   // 记录token是否正确
}

type Claims interface {
    Valid() error
}
// 签名方法 所有的签名方法都会实现这个接口
// 具体可以参考https://github.com/dgrijalva/jwt-go/blob/master/hmac.go
type SigningMethod interface {
    // 验证token的签名，如果有限返回nil
    Verify(signingString, signature string, key interface{}) error

    // 签名方法 接受头部和荷载编码过后的字符串和签名秘钥
    // 在hmac中key必须是Key must be []byte
    // 在rsa中key 必须是*rsa.PrivateKey 对象
    Sign(signingString string, key interface{}) (string, error)

    // 返回加密方法的名字 比如'HS256'
    Alg() string
}
//parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(cert.PrivateKey)


```

valid（）方法
```go

// Validate Claims
if !p.SkipClaimsValidation {
    //调用
    if err := token.Claims.Valid(); err != nil {

        // If the Claims Valid returned an error, check if it is a validation error,
        // If it was another error type, create a ValidationError with a generic ClaimsInvalid flag set
        if e, ok := err.(*ValidationError); !ok {
            vErr = &ValidationError{Inner: err, Errors: ValidationErrorClaimsInvalid}
        } else {
            vErr = e
        }
    }
}

```

Note：服务端生成的jwt返回客户端可以存到cookie也可以存到localStorage中（相比cookie容量大），存在cookie中需加上 HttpOnly 的标记，可以防止 XSS 攻击。

