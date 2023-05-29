<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [SHA系列](#sha%E7%B3%BB%E5%88%97)
  - [SHA-1](#sha-1)
  - [SHA-2/SHA-256](#sha-2sha-256)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# SHA系列
SHA系列包括多个散列算法标准,其中, SHA-1是数字签名标准中要求使用的算法。
- SHA-0：正式地称作SHA，这个版本在发行后不久被指出存在弱点。
- SHA-1：NIST于1994年发布的,它与MD4和MD5散列算法非常相似,被认为是MD4和MD5的后继者。(160位)
- SHA-2：实际上分为SHA-224、 SHA-256、SHA-384和SHA512算法。

## SHA-1
SHA-1（英语：Secure Hash Algorithm 1，中文名：安全散列算法1）是一种密码散列函数，美国国家安全局设计，并由美国国家标准技术研究所（NIST）发布为联邦数据处理标准（FIPS）。
SHA-1可以生成一个被称为消息摘要的160位（20字节）散列值，散列值通常的呈现形式为40个十六进制数


SHA-1已经不再视为可抵御有充足资金、充足计算资源的攻击者。2005年，密码分析人员发现了对SHA-1的有效攻击方法，这表明该算法可能不够安全，不能继续使用，自2010年以来，许多组织建议用SHA-2或SHA-3来替换SHA-1。Microsoft、Google以及Mozilla都宣布，它们旗下的浏览器将在2017年前停止接受使用SHA-1算法签名的SSL证书。
2017年2月23日，CWI Amsterdam与Google宣布了一个成功的SHA-1碰撞攻击，发布了两份内容不同但SHA-1散列值相同的PDF文件作为概念证明。


```go

func Sha1(data string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(data))
	return hex.EncodeToString(sha1.Sum([]byte("")))
}

```

## SHA-2/SHA-256

SHA-2有多种不同的位数，导致这个名词有一些混乱。但是无论是“SHA-2”，“SHA-256”或“SHA-256位”，其实都是指同一种加密算法。
但是SHA-224”，“SHA-384”或“SHA-512”，表示SHA-2的二进制长度。还要另一种就是会把算法和二进制长度都写上，如“SHA-2 384”。
SSL行业选择SHA作为数字签名的散列算法，从2011到2015，一直以SHA-1位主导算法。但随着互联网技术的提升，SHA-1的缺点越来越突显。从去年起，SHA-2成为了新的标准，所以现在签发的SSL证书，必须使用该算法签名。
也许有人偶尔会看到SHA-2 384位的证书，很少会看到224位，因为224位不允许用于公共信任的证书，512位，不被软件支持。
初步预计，SHA-2的使用年限为五年，但也许会被提前淘汰。这需要时间来验证
