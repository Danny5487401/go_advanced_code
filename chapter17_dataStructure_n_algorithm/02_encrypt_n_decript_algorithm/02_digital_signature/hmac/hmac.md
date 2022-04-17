# 消息认证
网络系统安全一般要考虑两个方面:
1. 加密保护传送的信息，使其可以抵抗被动攻击;
2. 要能防止对手对系统进行主动攻击，如伪造、篡改信息等。认证是对抗主动攻击的主要手段，它对于开放的网络中的各种信息系统的安全性有重要作用。认证分为实体认证和消息认证

在一个开放通信网络的环境中，信息面临的攻击包括窃听、伪造、修改、插入、删除、否认等。因此，需要提供用来验证消息完整性的一种机制或服务–消息认证。这种服务的主要功能包括：

- 确保收到的消息确实和发送的一样；
- 确保消息的来源真实有效

Note: 对称密码体制和公钥密码体制都可提供这种服务，但用于消息认证的最常见的密码技术是基于哈希函数的消息认证码。

## 消息认证码(MAC，Messages Authentication Codes)
与密钥相关的的单向散列函数，也称为消息鉴别码或是消息校验和。此时需要通信双方A和B共享一密钥K。


## HMAC 哈希运算消息认证码(Hash-based Message Authentication Code)
![](.hmac_images/hmac.png)
由H.Krawezyk，M.Bellare，R.Canetti于1996年提出的一种基于Hash函数和密钥进行消息认证的方法，并于1997年作为RFC2104被公布，并在IPSec和其他网络协议（如SSL）中得以广泛应用，现在已经成为事实上的Internet安全标准。它可以与任何迭代散列函数捆绑使用。
HMAC运算利用hash算法，以一个消息M和一个密钥K作为输入，生成一个定长的消息摘要作为输出。HMAC算法利用已有的Hash函数，关键问题是如何使用密钥。


和我们自定义的加salt算法不同，Hmac算法针对所有哈希算法都通用，无论是MD5还是SHA-1。
采用Hmac替代我们自己的salt算法，可以使程序算法更标准化，也更安全。


HMAC的密钥长度可以是任意大小，如果小于n（hash输出值的大小），那么将会消弱算法安全的强度。建议使用长度大于n的密钥，但是采用长度大的密钥并不意味着增强了函数的安全性。密钥应该是随机选取的，可以采用一种强伪随机发生器，并且密钥需要周期性更新，这样可以减少散列函数弱密钥的危险性以及已经暴露密钥所带来的破坏
使用SHA-1、SHA-224、SHA-256、SHA-384、SHA-512所构造的HMAC，分别称为HMAC-SHA1、HMAC-SHA-224、HMAC-SHA-384、HMAC-SHA-512。


## HMAC的Go实现

### crypto/hmac包
```go
func New(h func() hash.Hash, key []byte) hash.Hash
```
New函数返回一个采用hash.Hash作为底层hash接口、key作为密钥的HMAC算法的hash接口。

```go
func Equal(mac1, mac2 []byte) bool
```
比较两个MAC是否相同，而不会泄露对比时间信息。（以规避时间侧信道攻击：指通过计算比较时花费的时间的长短来获取密码的信息，用于密码破解）


### 使用
1. 生成MAC
```go
//key随意设置 data 要加密数据
func Hmac(key, data string) string {
	hash:= hmac.New(md5.New, []byte(key)) // 创建对应的md5哈希加密算法
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum([]byte("")))
}
func HmacSha256(key, data string) string {
	hash:= hmac.New(sha256.New, []byte(key)) //创建对应的sha256哈希加密算法
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum([]byte("")))
}

```