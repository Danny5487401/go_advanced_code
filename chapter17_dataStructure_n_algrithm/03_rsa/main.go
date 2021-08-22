package main

import (
	"encoding/base64"
	"go_advenced_code/chapter17_dataStructure_n_algrithm/03_rsa/security"
	"log"
)

/*
go标准库中仅有"公钥加密，私钥解密"，而没有“私钥加密、公钥解密”。经过考虑，我认为GO的开发者是故意这样设计的，原因如下：

	1.非对称加密相比对称加密的好处就是：私钥自己保留，公钥公布出去，公钥加密后只有私钥能解开，私钥加密后只有公钥能解开。
	2.如果仅有一对密钥，与对称加密区别就不大了。
	假如你是服务提供方，使用私钥进行加密后，接入方使用你提供的公钥进行解密，一旦这个公钥泄漏，带来的后果和对称加密密钥泄漏是一样的。只有双方互换公钥（均使用对方公钥加密，己方私钥解密），才能充分发挥非对称加密的优势。

*/

func main() {
	var mingwen = "Danny最帅"
	//md5 := security.MD5([]byte(mingwen))
	////MD5打印为16进制字符串
	//log.Println(hex.EncodeToString(md5))

	//RSA的内容使用base64打印
	privateKey, publicKey, _ := security.GenRSAKey(1024)
	log.Println("rsa私钥:\t", base64.StdEncoding.EncodeToString(privateKey))
	log.Println("rsa公钥:\t", base64.StdEncoding.EncodeToString(publicKey))

	miwen, err := security.RsaEncryptBlock([]byte(mingwen), publicKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("加密后：\t", base64.StdEncoding.EncodeToString(miwen))

	jiemi, err := security.RsaDecryptBlock(miwen, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("解密后：\t", string(jiemi))

}
