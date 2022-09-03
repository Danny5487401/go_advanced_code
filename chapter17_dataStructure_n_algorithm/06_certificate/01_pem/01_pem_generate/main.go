package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var bits int
	flag.IntVar(&bits, "b", 1024, "密钥长度，默认是1024")
	flag.Parse()
	//生成rsa密钥文件
	if GenRsaKey(bits) != nil {
		log.Fatalln("密钥文件生成失败！")
	}
	log.Println("密钥文件生成成功！")
}

func GenRsaKey(bits int) error {
	// 1. 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 2. 使用x509.MarshalPKCS1PrivateKey序列化私钥为derText
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)

	// 3. 使用pem.Block转为Block
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY", // 头部的type，-----BEGIN Type-----
		Bytes: derStream,         //内容
	}
	file, err := os.Create("chapter17_dataStructure_n_algorithm/06_certificate/pem_file/private.pem")
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 4. 使用pem.Encode写入文件
	err = pem.Encode(file, block)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 5. 从私钥中获取公钥
	publicKey := &privateKey.PublicKey
	// 6. 序列化公钥为derStream
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 7. 使用pem.Block转为Block
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create("chapter17_dataStructure_n_algorithm/06_certificate/pem_file/public.pem")
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 8. 使用pem.Encode写入文件
	err = pem.Encode(file, block)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
