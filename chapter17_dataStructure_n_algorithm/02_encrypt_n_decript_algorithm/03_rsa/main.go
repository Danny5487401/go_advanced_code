package main

import (
	"encoding/base64"
	"github.com/Danny5487401/go_advanced_code/chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/03_rsa/security"
	"log"
)

func main() {
	var mingwen = "Danny最帅"

	// 1. 生成RSA密钥对der
	privateKey, publicKey, _ := security.GenRSAKey(1024)
	//RSA的内容使用base64打印
	log.Println("rsa私钥:\t", base64.StdEncoding.EncodeToString(privateKey))
	log.Println("rsa公钥:\t", base64.StdEncoding.EncodeToString(publicKey))

	// 方式一：不分段
	// 2. 不分段用公钥加密数据
	miwen, err := security.RsaEncrypt([]byte(mingwen), publicKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("加密后：\t", base64.StdEncoding.EncodeToString(miwen))

	// 3. 不分段用私钥解密数据
	jiemi, err := security.RsaDecrypt(miwen, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("解密后：\t", string(jiemi))

	// 方式二：分段
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

//Note：加密的内容不能太大，最大不能超过秘钥的长度
