package security

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"log"
)

/*
PEM
	将X.509基础证书用base64重新编码存为ASCII文件，用于和邮件一起传输保证邮件安全性。
X.509
	这是由国际电信联盟（ITU-T）制定的ASN.1规范下数字证书标准，它规定了证书应包含哪些信息和使用什么样的编码格式（默认DER二进制编码）。
PKCS，
	The Public-Key Cryptography Standards，公钥密码学标准，由美帝的RSA公司制定的一系列标准，这里我们只讨论#7/#8/#12，
	#7/#12是对X.509证书进行扩展、加密用于交换。#8是一种私钥格式标准。openssl生成的私钥，可以转换成pkcs8格式
*/

//生成RSA密钥对
func GenRSAKey(size int) (privateKeyBytes, publicKeyBytes []byte, err error) {
	//生成密钥
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return
	}
	privateKeyBytes = x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyBytes = x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	return
}

//公钥加密
func RsaEncrypt(src, publicKeyByte []byte) (bytes []byte, err error) {
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyByte)
	if err != nil {
		return
	}
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, src)
}

// 公钥加密-分段
func RsaEncryptBlock(src, publicKeyByte []byte) (bytesEncrypt []byte, err error) {
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyByte)
	if err != nil {
		return
	}
	keySize, srcSize := publicKey.Size(), len(src)
	log.Println("密钥长度：", keySize, "\t明文长度：\t", srcSize)
	//单次加密的长度需要减掉padding的长度，PKCS1为11
	offSet, once := 0, keySize-11
	buffer := bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + once
		if endIndex > srcSize {
			endIndex = srcSize
		}
		// 加密一部分
		bytesOnce, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, src[offSet:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	bytesEncrypt = buffer.Bytes()
	return
}

// 私钥解密
func RsaDecrypt(src, privateKeyBytes []byte) (bytesDecrypt []byte, err error) {
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return
	}
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, src)
}

// 私钥解密-分段
func RsaDecryptBlock(src, privateKeyBytes []byte) (bytesDecrypt []byte, err error) {
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return
	}
	keySize := privateKey.Size()
	srcSize := len(src)
	log.Println("密钥长度：", keySize, "\t密文长度：\t", srcSize)
	var offSet = 0
	var buffer = bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + keySize
		if endIndex > srcSize {
			endIndex = srcSize
		}
		bytesOnce, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, src[offSet:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	bytesDecrypt = buffer.Bytes()
	return
}
