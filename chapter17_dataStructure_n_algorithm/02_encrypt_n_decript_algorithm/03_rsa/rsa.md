# RSA（Rivest–Shamir–Adleman）加密
RSA以麻省理工学院的科学家（Rivest、Shamir 和 Adleman）的名字命名。

RSA算法的安全性基于RSA问题的困难性，也就是基于大整数因子分解的困难性上。但是RSA问题不会比因子分解问题更加困难，
也就是说，在没有解决因子分解问题的情况下可能解决RSA问题，因此RSA算法并不是完全基于大整数因子分解的困难性上的。

它是一种非对称加密算法，也叫”单向加密“。用这种方式，任何人都可以很容易地对数据进行加密，而只有用正确的”秘钥“才能解密
RSA加密是最常用的非对称加密方式，原理是对一极大整数做因数分解的困难性来保证安全性。通常个人保存私钥，公钥是公开的（可能同时多人持有）

## RSA加密有常见的三种情况：

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

实现：  
    go标准库中仅有"公钥加密，私钥解密"，而没有“私钥加密、公钥解密”。经过考虑，我认为GO的开发者是故意这样设计的

原因如下：
1. 非对称加密相比对称加密的好处就是：私钥自己保留，公钥公布出去，公钥加密后只有私钥能解开，私钥加密后只有公钥能解开。  
2. 如果仅有一对密钥，与对称加密区别就不大了。

    假如你是服务提供方，使用私钥进行加密后，接入方使用你提供的公钥进行解密，一旦这个公钥泄漏，带来的后果和对称加密密钥泄漏是一样的。
    只有双方互换公钥（均使用对方公钥加密，己方私钥解密），才能充分发挥非对称加密的优势。

## 应用
RSA 算法需要的计算量比 AES 高，但速度要慢得多。它比较适合用于加密少量数据。

## 相关概念

### x.509

X.509标准是密码学里公钥证书的格式标准。X.509 证书己应用在包括TLS/SSL（WWW万维网安全浏览的基石）在内的众多 Internet协议里，同时它也有很多非在线的应用场景，比如电子签名服务。

X.509证书含有公钥和标识（主机名、组织或个人），并由证书颁发机构（CA）签名（或自签名）。对于一份经由可信的证书签发机构签名（或者可以通过其它方式验证）的证书，证书的拥有者就可以用证书及相应的私钥来创建安全的通信，以及对文档进行数字签名。
X.509证书的结构是用ASN.1(Abstract Syntax Notation One：抽象语法标记)来描述其数据结构，并使用ASN1语法进行编码。

X.509 v3数字证书的结构如下：
```
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            10:e6:fc:62:b7:41:8a:d5:00:5e:45:b6
    Signature Algorithm: sha256WithRSAEncryption
        Issuer: C=BE, O=GlobalSign nv-sa, CN=GlobalSign Organization Validation CA - SHA256 - G2
        Validity
            Not Before: Nov 21 08:00:00 2016 GMT
            Not After : Nov 22 07:59:59 2017 GMT
        Subject: C=US, ST=California, L=San Francisco, O=Wikimedia Foundation, Inc., CN=*.wikipedia.org
        Subject Public Key Info:
            Public Key Algorithm: id-ecPublicKey
                Public-Key: (256 bit)
                pub: 
                    04:c9:22:69:31:8a:d6:6c:ea:da:c3:7f:2c:ac:a5:
                    af:c0:02:ea:81:cb:65:b9:fd:0c:6d:46:5b:c9:1e:
                    ed:b2:ac:2a:1b:4a:ec:80:7b:e7:1a:51:e0:df:f7:
                    c7:4a:20:7b:91:4b:20:07:21:ce:cf:68:65:8c:c6:
                    9d:3b:ef:d5:c1
                ASN1 OID: prime256v1
                NIST CURVE: P-256
        X509v3 extensions:
            X509v3 Key Usage: critical
                Digital Signature, Key Agreement
            Authority Information Access: 
                CA Issuers - URI:http://secure.globalsign.com/cacert/gsorganizationvalsha2g2r1.crt
                OCSP - URI:http://ocsp2.globalsign.com/gsorganizationvalsha2g2
 
            X509v3 Certificate Policies: 
                Policy: 1.3.6.1.4.1.4146.1.20
                  CPS: https://www.globalsign.com/repository/
                Policy: 2.23.140.1.2.2
 
            X509v3 Basic Constraints: 
                CA:FALSE
            X509v3 CRL Distribution Points: 
 
                Full Name:
                  URI:http://crl.globalsign.com/gs/gsorganizationvalsha2g2.crl
 
            X509v3 Subject Alternative Name: 
                DNS:*.wikipedia.org, DNS:*.m.mediawiki.org, DNS:*.m.wikibooks.org, DNS:*.m.wikidata.org, DNS:*.m.wikimedia.org, DNS:*.m.wikimediafoundation.org, DNS:*.m.wikinews.org, DNS:*.m.wikipedia.org, DNS:*.m.wikiquote.org, DNS:*.m.wikisource.org, DNS:*.m.wikiversity.org, DNS:*.m.wikivoyage.org, DNS:*.m.wiktionary.org, DNS:*.mediawiki.org, DNS:*.planet.wikimedia.org, DNS:*.wikibooks.org, DNS:*.wikidata.org, DNS:*.wikimedia.org, DNS:*.wikimediafoundation.org, DNS:*.wikinews.org, DNS:*.wikiquote.org, DNS:*.wikisource.org, DNS:*.wikiversity.org, DNS:*.wikivoyage.org, DNS:*.wiktionary.org, DNS:*.wmfusercontent.org, DNS:*.zero.wikipedia.org, DNS:mediawiki.org, DNS:w.wiki, DNS:wikibooks.org, DNS:wikidata.org, DNS:wikimedia.org, DNS:wikimediafoundation.org, DNS:wikinews.org, DNS:wikiquote.org, DNS:wikisource.org, DNS:wikiversity.org, DNS:wikivoyage.org, DNS:wiktionary.org, DNS:wmfusercontent.org, DNS:wikipedia.org
            X509v3 Extended Key Usage: 
                TLS Web Server Authentication, TLS Web Client Authentication
            X509v3 Subject Key Identifier: 
                28:2A:26:2A:57:8B:3B:CE:B4:D6:AB:54:EF:D7:38:21:2C:49:5C:36
            X509v3 Authority Key Identifier: 
                keyid:96:DE:61:F1:BD:1C:16:29:53:1C:C0:CC:7D:3B:83:00:40:E6:1A:7C
 
    Signature Algorithm: sha256WithRSAEncryption
         8b:c3:ed:d1:9d:39:6f:af:40:72:bd:1e:18:5e:30:54:23:35:
         ...
```
- Certificate 证书

    - Version Number版本号
    
    - Serial Number序列号

    - ID Signature Algorithm ID签名算法

        - Issuer Name颁发者名称
        
        - Validity period 有效期 
        
            - Not before起始日期 
            - Not after截至日期
    
        - Subject Name主题名称
        
        - Subject pbulic Key Info 主题公钥信息 
    
            - Public Key Algorithm公钥算法
    
                - Subject Public Key主题公钥

    - Issuer Unique Identifier (optional)颁发者唯一标识符（可选）
    
    - Subject Unique Identifier (optional)主题唯一标识符（可选）

    - Extensions (optional) 证书的扩展项（可选）

    - Certificate Sigature Algorithm证书签名算法
    
    - Certificate Signature证书的签名

Note：对于所有的版本，同一个CA颁发的证书序列号都必须是唯一的。

### PKCS（The Public-Key Cryptography Standards）公钥密码学标准
PKCS代表“公钥密码学标准”。这是一组由RSA Security Inc.设计和发布的公钥密码标准，始于20世纪90年代初，该公司发布这些标准是为了推广使用他们拥有专利的密码技术，如RSA算法、Schnorr签名算法和其他一些算法。

尽管不是行业标准（因为该公司保留了对它们的控制权），但近年来某些标准已经开始进入IETF和PKIX工作组等相关标准化组织的“标准跟踪”过程

● PKCS＃8 1.2私钥信息语法标准，请参见RFC5958。用于携带私钥证书密钥对（加密或未加密）。

● PKCS＃9 2.0选定的属性类型[，请参见RFC2985。定义选定的属性类型，以便在PKCS＃6扩展证书、PKCS＃7数字签名消息、PKCS＃8私钥信息和PKCS＃10证书签名请求中使用。

● PKCS＃10 1.7认证请求标准，请参阅RFC2986。发送给认证机构以请求公钥证书的消息格式，请参阅证书签名请求。

● PKCS＃11 2.40密码令牌接口，也称为“ Cryptoki”。定义密码令牌通用接口的API（另请参阅硬件安全模块）。常用于单点登录，公共密钥加密和磁盘加密[10]系统。 RSA Security已将PKCS＃11标准的进一步开发移交给了OASIS PKCS 11技术委员会。

● PKCS＃12 1.1个人信息交换语法标准，请参阅RFC7292。定义一种文件格式，个人信息交换语法标准[11]见RFC 7292。定义一种文件格式，通常用于存储私钥和附带的公钥证书，并使用基于Password的对称密钥进行保护。PFX是PKCS#12的前身。

## 密钥
常见的几种秘钥存储格式有：字符串、证书文件、n/e参数等

### 1. 字符串格式
这是最常见的一种形式，通常RSA的秘钥都是以hex、base64编码后的字符串提供，如证书内的秘钥格式即是base64编码的字符串，然后添加前后的具体标识实现的。可以通过解码字符串，构建公钥/私钥。

Note:base64存在几种细节不同的编码格式，StdEncoding、URLEncoding、RawStdEncoding、RawURLEncoding，使用时还需要进一步确认秘钥具体编码格式，避免解码出错。
以下未特殊说明的例子中均默认使用StdEncoding。

（1）公钥
直接hex、base64解码后调用x509.ParsePKIXPublicKey即可
```go
key, _ := hex.DecodeString(publicKeyStr)
publicKey, _ := x509.ParsePKIXPublicKey(key)

```

（2）私钥
由于RSA私钥存在PKCS1和PKCS8两种格式，因此解码后需要根据格式类型调用对应的方法即可。一般java使用pkcs8格式的私钥，其他语言使用pkcs1格式的私钥。使用时，记得确认下格式。
```go
//解析pkcs1格式私钥
key, _ := base64.StdEncoding.DecodeString(pkcs1keyStr)
privateKey, _ := x509.ParsePKCS1PrivateKey(key)

//解析pkcs8格式私钥
key, _ := hex.DecodeString(pkcs8keyStr)
privateKey, err := x509.ParsePKCS8PrivateKey(key)

```

### 2. 证书文件扩展名

Note: 其中一些扩展名也有其它用途，就是说具有这个扩展名的文件可能并不是证书，比如说可能只是保存了私钥。

（1）.pem、.cert、.cer、.crt

.pem、.cert、.cer、.crt等都是pem格式的文件，只是文件后缀不一。

- PEM是Privacy Enhance Mail(隐私增强型电子邮件)的缩写，DER编码的证书再进行Base64编码(即对字符串格式私钥的文件化处理)，再加上开始和结束行,即数据存放于“--- BEGIN CERTIFICATE ---”和“ --- END CERTIFICATE ---”之间

- .cer，.crt，.der：通常采用二进制DER形式，但Base64编码也很常见

解析方式：读取文件，调用pem.Decode，然后按照base64解码，再解析成公钥/私钥。
```go
key,_ := ioutil.ReadFile("pem_file_path")
block, _ := pem.Decode(key)
//解析成pkcs8格式私钥
privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
//解析成pkcs1格式私钥
privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
//解析成公钥
publicKey, _ := x509.ParsePKIXPublicKey(key)


```

（2）.pkcs12、.pfx、.p12

.pkcs12、.pfx、.p12这些文件格式存储的是已加密后的内容，可以通过openssl转换成pem文件后进行处理。

-  .p12-PKCS＃12：可以包含证书（公钥），也可同时包含受密码保护的私钥
-  .pfx ：PKCS＃12的前身（通常用PKCS＃12格式，例如IIS产生的PFX文件）

提取密钥对：
```shell
openssl pkcs12 -in in.p12 -out out.pem -nodes
```

### 3. N,E参数

例如：login with apple keys的公钥就是这种格式的，需要根据n,e构造出公钥。
```shell
{
      "kty": "RSA",
      "kid": "eXaunmL",
      "use": "sig",
      "alg": "RS256",
      "n": "4dGQ7bQK8LgILOdLsYzfZjkEAoQeVC_aqyc8GC6RX7dq_KvRAQAWPvkam8VQv4GK5T4ogklEKEvj5ISBamdDNq1n52TpxQwI2EqxSk7I9fKPKhRt4F8-2yETlYvye-2s6NeWJim0KBtOVrk0gWvEDgd6WOqJl_yt5WBISvILNyVg1qAAM8JeX6dRPosahRVDjA52G2X-Tip84wqwyRpUlq2ybzcLh3zyhCitBOebiRWDQfG26EH9lTlJhll-p_Dg8vAXxJLIJ4SNLcqgFeZe4OfHLgdzMvxXZJnPp_VgmkcpUdRotazKZumj6dBPcXI_XID4Z4Z3OM1KrZPJNdUhxw",
      "e": "AQAB"
    }

```

使用时就需要，将N，E解析成big.Int格式，注意N、E的base64的具体编码格式：

```go
pubN, _ := parse2bigInt(n)
pubE, _ := parse2bigInt(e)
pub = &rsa.PublicKey{
    N: pubN,
    E: int(pubE.Int64()),
}

// parse string to big.Int
func parse2bigInt(s string) (bi *big.Int, err error) {
    bi = &big.Int{}
    b, err := base64.RawURLEncoding.DecodeString(s)//此处使用的是RawURLEncoding
    if err != nil {
        return
    }
    bi.SetBytes(b)
    return
}

```