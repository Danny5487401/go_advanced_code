package main

import (
	"encoding/base64"
	"go_advanced_code/chapter17_dataStructure_n_algrithm/03_rsa/security"
	"log"
)

/*
Rsa:
	RSA加密是最常用的非对称加密方式，原理是对一极大整数做因数分解的困难性来保证安全性。通常个人保存私钥，公钥是公开的（可能同时多人持有）
RSA加密有常见的三种情况：

	1。公钥加密，私钥解密
		最常用的一种情况，对接过支付宝就应该碰到过。
		接收方存一个私钥，发送方保存对应的公钥用来发送消息加密，能够确认消息不被泄露。
		所以支付宝会把他生成的一份公钥给你， 你需要把你生成的公钥给支付宝，你们的通信是建立在两对公私钥之上的。

	2.私钥加密，公钥验签
		发送者把原文和密文同时发布，客户端使用公钥确认是由真正的发送者发出的。
		这种情况一般用于确认消息发布的真实性，可用于推送、广播和公开消息验证场景。

	3.私钥加密，公钥解密
		这个本身不在推荐规范内的场景，现在却是很常见。
		例如离线软件授权，发布出去的软件里面保存一个公钥，软件厂商使用私钥加密包含到期时间的原文，得到密文，也就是授权码，软件验证的时候使用公钥解密授权码，对比当前时间，确认是否过期。
实现：
	go标准库中仅有"公钥加密，私钥解密"，而没有“私钥加密、公钥解密”。经过考虑，我认为GO的开发者是故意这样设计的，原因如下：
		1.非对称加密相比对称加密的好处就是：私钥自己保留，公钥公布出去，公钥加密后只有私钥能解开，私钥加密后只有公钥能解开。
		2.如果仅有一对密钥，与对称加密区别就不大了。
		假如你是服务提供方，使用私钥进行加密后，接入方使用你提供的公钥进行解密，一旦这个公钥泄漏，带来的后果和对称加密密钥泄漏是一样的。
		只有双方互换公钥（均使用对方公钥加密，己方私钥解密），才能充分发挥非对称加密的优势。

*/

func main() {
	var mingwen = "Danny最帅"

	//1.生成RSA密钥对
	privateKey, publicKey, _ := security.GenRSAKey(1024)
	//RSA的内容使用base64打印
	log.Println("rsa私钥:\t", base64.StdEncoding.EncodeToString(privateKey))
	log.Println("rsa公钥:\t", base64.StdEncoding.EncodeToString(publicKey))

	// 2.不分段
	miwen, err := security.RsaEncrypt([]byte(mingwen), publicKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("加密后：\t", base64.StdEncoding.EncodeToString(miwen))

	jiemi, err := security.RsaDecrypt(miwen, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("解密后：\t", string(jiemi))

	// 2.分段
	//// 分段加密
	//miwen, err := security.RsaEncryptBlock([]byte(mingwen), publicKey)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("加密后：\t", base64.StdEncoding.EncodeToString(miwen))
	//
	//// 分段解密
	//jiemi, err := security.RsaDecryptBlock(miwen, privateKey)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("解密后：\t", string(jiemi))

}
