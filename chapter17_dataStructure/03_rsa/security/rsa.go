package security

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"log"
)

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
