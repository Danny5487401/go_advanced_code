<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [RSA（Rivest–Shamir–Adleman）加密](#rsarivestshamiradleman%E5%8A%A0%E5%AF%86)
  - [前置知识](#%E5%89%8D%E7%BD%AE%E7%9F%A5%E8%AF%86)
    - [1. 欧拉函数](#1-%E6%AC%A7%E6%8B%89%E5%87%BD%E6%95%B0)
    - [2. 欧拉定理](#2-%E6%AC%A7%E6%8B%89%E5%AE%9A%E7%90%86)
    - [3. 模反元素](#3-%E6%A8%A1%E5%8F%8D%E5%85%83%E7%B4%A0)
  - [RSA密钥对生成算法](#rsa%E5%AF%86%E9%92%A5%E5%AF%B9%E7%94%9F%E6%88%90%E7%AE%97%E6%B3%95)
  - [RSA加密有常见的三种情况](#rsa%E5%8A%A0%E5%AF%86%E6%9C%89%E5%B8%B8%E8%A7%81%E7%9A%84%E4%B8%89%E7%A7%8D%E6%83%85%E5%86%B5)
    - [实现](#%E5%AE%9E%E7%8E%B0)
  - [应用](#%E5%BA%94%E7%94%A8)
  - [Go源码RSA](#go%E6%BA%90%E7%A0%81rsa)
    - [crypto/x509包](#cryptox509%E5%8C%85)
      - [序列化](#%E5%BA%8F%E5%88%97%E5%8C%96)
      - [解析](#%E8%A7%A3%E6%9E%90)
    - [crypto/rsa包](#cryptorsa%E5%8C%85)
      - [生成RSA密钥对](#%E7%94%9F%E6%88%90rsa%E5%AF%86%E9%92%A5%E5%AF%B9)
  - [参考链接](#%E5%8F%82%E8%80%83%E9%93%BE%E6%8E%A5)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# RSA（Rivest–Shamir–Adleman）加密
1977年，三位数学家Rivest、Shamir 和 Adleman 设计了一种算法，可以实现非对称加密。
这种算法用他们三个人的名字命名，叫做RSA算法。从那时直到现在，RSA算法一直是最广为使用的"非对称加密算法"。
毫不夸张地说，只要有计算机网络的地方，就有RSA算法。

这种算法非常可靠，密钥越长，它就越难破解。根据已经披露的文献，目前被破解的最长RSA密钥是768个二进制位。也就是说，长度超过768位的密钥，还无法破解（至少没人公开宣布）。
因此可以认为，1024位的RSA密钥基本安全，2048位的密钥极其安全。

RSA算法的安全性基于RSA问题的困难性，也就是基于大整数因子分解的困难性上。但是RSA问题不会比因子分解问题更加困难，
也就是说，在没有解决因子分解问题的情况下可能解决RSA问题，因此RSA算法并不是完全基于大整数因子分解的困难性上的。

它是一种非对称加密算法，也叫”单向加密“。用这种方式，任何人都可以很容易地对数据进行加密，而只有用正确的”秘钥“才能解密
RSA加密是最常用的非对称加密方式，原理是对一极大整数做因数分解的困难性来保证安全性。通常个人保存私钥，公钥是公开的（可能同时多人持有）

## 前置知识

### 1. 欧拉函数
>任意给定正整数n，请问在小于等于n的正整数之中，有多少个与n构成互质关系？（比如，在1到8之中，有多少个数与8构成互质关系？）

计算这个值的方法就叫做欧拉函数，以φ(n)表示。在1到8之中，与8形成互质关系的是1、3、5、7，所以 φ(n) = 4。

### 2. 欧拉定理
拉函数的用处，在于[欧拉定理]。"欧拉定理"指的是:如果两个正整数a和n互质，则n的欧拉函数 φ(n) 可以让下面的等式成立：
![](.rsa_images/oula_discipline.png)

也就是说，a的φ(n)次方被n除的余数为1。或者说，a的φ(n)次方减去1，可以被n整除。
比如，3和7互质，而7的欧拉函数φ(7)等于6，所以3的6次方（729）减去1，可以被7整除（728/7=104）

### 3. 模反元素
如果两个正整数a和n互质，那么一定可以找到整数b，使得 ab-1 被n整除，或者说ab被n除的余数是1。
![](.rsa_images/mod_elem.png)  
这时，b就叫做a的"模反元素"。

比如，3和11互质，那么3的模反元素就是4，因为 (3 × 4)-1 可以被11整除。显然，模反元素不止一个， 4加减11的整数倍都是3的模反元素 {...,-18,-7,4,15,26,...}，即如果b是a的模反元素，则 b+kn 都是a的模反元素。


## RSA密钥对生成算法
![](.rsa_images/rsa_process.png)

1. 随机选择两个不相等的质数p和q。
> 爱丽丝选择了61和53。（实际应用中，这两个质数越大，就越难破解。）
   
2. 计算p和q的乘积n。
> 爱丽丝就把61和53相乘。n = 61×53 = 3233;   

n的长度就是密钥长度.  3233写成二进制是110010100001，一共有12位，所以这个密钥就是12位。实际应用中，RSA密钥一般是1024位，重要场合则为2048位。

3. 计算n的欧拉函数φ(n)。
> φ(n) = (p-1)(q-1)    

爱丽丝算出φ(3233)等于60×52，即3120。

4. 随机选择一个整数e
>爱丽丝就在1到3120之间，随机选择了17。（实际应用中，常常选择65537）
 
5. 计算e对于φ(n)的模反元素d: 所谓"模反元素"就是指有一个整数d，可以使得ed被φ(n)除的余数为1
> ed ≡ 1 (mod φ(n)) 这个式子等价于 ed - 1 = kφ(n)

于是，找到模反元素d，实质上就是对下面这个二元一次方程求解。

> ex + φ(n)y = 1

这个方程可以用"扩展欧几里得算法"求解，此处省略具体过程。总之，爱丽丝算出一组整数解为 (x,y)=(2753,-15)，即 d=2753。

至此所有计算完成。

6. 将n和e封装成公钥，n和d封装成私钥。
> 在爱丽丝的例子中，n=3233，e=17，d=2753，所以公钥就是 (3233,17)，私钥就是（3233, 2753）。

## RSA加密有常见的三种情况

1. 公钥加密，私钥解密   
    最常用的一种情况，对接过支付宝就应该碰到过。
    接收方存一个私钥，发送方保存对应的公钥用来发送消息加密，能够确认消息不被泄露。
    所以支付宝会把他生成的一份公钥给你， 你需要把你生成的公钥给支付宝，你们的通信是建立在两对公私钥之上的。

2. 私钥加密，公钥验签   
    发送者把原文和密文同时发布，客户端使用公钥确认是由真正的发送者发出的。
    这种情况一般用于确认消息发布的真实性，可用于推送、广播和公开消息验证场景。

3. 私钥加密，公钥解密   
    这个本身不在推荐规范内的场景，现在却是很常见。
    例如离线软件授权，发布出去的软件里面保存一个公钥，软件厂商使用私钥加密包含到期时间的原文，得到密文，也就是授权码，
    软件验证的时候使用公钥解密授权码，对比当前时间，确认是否过期。

### 实现  

go标准库中仅有"公钥加密，私钥解密"，而没有“私钥加密、公钥解密”。经过考虑，我认为GO的开发者是故意这样设计的

原因如下：
1. 非对称加密相比对称加密的好处就是：私钥自己保留，公钥公布出去，公钥加密后只有私钥能解开，私钥加密后只有公钥能解开。  
2. 如果仅有一对密钥，与对称加密区别就不大了。

假如你是服务提供方，使用私钥进行加密后，接入方使用你提供的公钥进行解密，一旦这个公钥泄漏，带来的后果和对称加密密钥泄漏是一样的。
只有双方互换公钥（均使用对方公钥加密，己方私钥解密），才能充分发挥非对称加密的优势。

## 应用
RSA 算法需要的计算量比 AES 高，所以速度要慢得多。它比较适合用于加密少量数据。






## Go源码RSA

### crypto/x509包
#### 序列化
```go
// /usr/local/go/src/crypto/x509/pkcs1.go
func MarshalPKCS1PrivateKey(key *rsa.PrivateKey) []byte
```
MarshalPKCS1PrivateKey将rsa私钥序列化为ASN.1 PKCS#1 DER编码。

```go
func MarshalPKIXPublicKey(pub interface{}) ([]byte, error)
```
MarshalPKIXPublicKey将公钥序列化为PKIX格式DER编码。


#### 解析
```go
func ParsePKIXPublicKey(derBytes []byte) (pub interface{}, err error)
```
ParsePKIXPublicKey解析一个DER编码的公钥。这些公钥一般在以"BEGIN PUBLIC KEY"出现的PEM块中。

```go
func ParsePKCS1PrivateKey(der []byte) (key *rsa.PrivateKey, err error)
```
ParsePKCS1PrivateKey解析ASN.1 PKCS#1 DER编码的rsa私钥。



### crypto/rsa包

```go
func EncryptPKCS1v15(rand io.Reader, pub *PublicKey, msg []byte) (out []byte, err error)
```
EncryptPKCS1v15使用PKCS#1 v1.5规定的填充方案和RSA算法加密msg。信息不能超过((公共模数的长度)-11)字节。
注意：使用本函数加密明文（而不是会话密钥）是危险的，请尽量在新协议中使用RSA OAEP。

```go
func DecryptPKCS1v15(rand io.Reader, priv *PrivateKey, ciphertext []byte) (out []byte, err error)
```
DecryptPKCS1v15使用PKCS#1 v1.5规定的填充方案和RSA算法解密密文。如果random不是nil，函数会注意规避时间侧信道攻击。

#### 生成RSA密钥对
1. 使用rsa.GenerateKey生成私钥
2. 使用x509.MarshalPKCS1PrivateKey序列化私钥为derText
3. 使用pem.Block转为Block
4. 使用pem.Encode写入文件
5. 从私钥中获取公钥
6. 使用x509.MarshalPKIXPublicKey序列化公钥为derStream
7. 使用pem.Block转为Block
8. 使用pem.Encode写入文件



## 参考链接
1. rsa算法原理：https://www.kancloud.cn/kancloud/rsa_algorithm/48484