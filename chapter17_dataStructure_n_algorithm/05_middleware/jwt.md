#Jwt（JSON Web Token）
JSON Web Token（JWT）是一个开放标准（RFC 7519），它定义了一种紧凑且自包含的方式，用于在各方之间以JSON方式安全地传输信息。
由于此信息是经过数字签名的，因此可以被验证和信任。可以使用秘密（使用HMAC算法）或使用RSA或ECDSA的公钥/私钥对对JWT进行签名

直白的讲jwt就是一种用户认证（区别于session、cookie）的解决方案。

## jwt构成：
- Header：TOKEN 的类型，就是JWT，签名的算法，如 HMAC SHA256、HS384
- Payload：载荷又称为Claim，携带的信息，比如用户名、过期时间等，一般叫做 Claim
- Signature：签名，是由header、payload 和你自己维护的一个 secret 经过加密得来的


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

