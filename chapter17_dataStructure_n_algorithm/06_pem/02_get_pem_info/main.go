package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/pkg/errors"

	"io/ioutil"
	"log"
)

// 全局变量
var privateKey, publicKey []byte

func init() {
	var err error
	publicKey, err = ioutil.ReadFile("chapter17_dataStructure_n_algorithm/06_pem/pem_file/public.pem")
	if err != nil {
		log.Fatalln(err)
	}
	privateKey, err = ioutil.ReadFile("chapter17_dataStructure_n_algorithm/06_pem/pem_file/private.pem")
	if err != nil {
		log.Fatalln(err)
	}
}

/**
 * 功能：获取RSA公钥长度
 * 参数：public
 * 返回：成功则返回 RSA 公钥长度，失败返回 error 错误信息
 */
func GetPubKeyLen(pubKey []byte) (int, error) {
	if pubKey == nil {
		return 0, errors.New("input arguments error")
	}
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return 0, errors.New("public rsaKey error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return 0, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return pub.N.BitLen(), nil
}

/*
   获取RSA私钥长度
   PriKey
   成功返回 RSA 私钥长度，失败返回error
*/
func GetPriKeyLen(priKey []byte) (int, error) {
	if priKey == nil {
		return 0, errors.New("input arguments error")
	}
	block, _ := pem.Decode(priKey)
	if block == nil {
		return 0, errors.New("private rsaKey error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return 0, err
	}
	return priv.N.BitLen(), nil
}

func main() {
	// 获取rsa 公钥长度
	pubKeyLen, _ := GetPubKeyLen(publicKey)
	fmt.Println(pubKeyLen)

	// 获取rsa 私钥长度
	privateLen, _ := GetPriKeyLen(privateKey)
	fmt.Println(privateLen)
}
