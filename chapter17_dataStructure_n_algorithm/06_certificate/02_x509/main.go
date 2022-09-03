package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"io/ioutil"
	"log"
	"math/big"
	"time"
)

func main() {
	// ca 证书
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Country:            []string{"China"},
			Organization:       []string{"Devops"},
			OrganizationalUnit: []string{"unit1"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 5},
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	// 私钥及公钥 rsa 格式
	caSelfSignedPrivateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	caSelfSignedPublicKey := &caSelfSignedPrivateKey.PublicKey
	// 自签证书 []byte
	caSelfSigned, err := x509.CreateCertificate(rand.Reader, ca, ca, caSelfSignedPublicKey, caSelfSignedPrivateKey)
	if err != nil {
		log.Println("create ca failed", err)
		return
	}
	caSelfSignedFile := "ca.pem"
	log.Println("write to", caSelfSignedFile)
	ioutil.WriteFile(caSelfSignedFile, caSelfSigned, 0777) // 将自签证书写入文件

	caSelfSignedPrivateKeyFile := "ca.key"
	caSelfSignedPrivateKeyDER := x509.MarshalPKCS1PrivateKey(caSelfSignedPrivateKey) // 将私钥转换为 DER 格式
	log.Println("write to", caSelfSignedPrivateKeyFile)
	ioutil.WriteFile(caSelfSignedPrivateKeyFile, caSelfSignedPrivateKeyDER, 0777) // 将 DER 编码私钥写入文件

	// 待签署证书及其私钥公钥
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Country:            []string{"China"},
			Organization:       []string{""},
			OrganizationalUnit: []string{""},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	certPrivateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	certPublicKey := &certPrivateKey.PublicKey

	// 使用自签CA 对 证书签署
	certSigned, err2 := x509.CreateCertificate(rand.Reader, cert, ca, certPublicKey, caSelfSignedPrivateKey)
	if err != nil {
		log.Println("create cert2 failed", err2)
		return
	}

	certFile := "cert.pem"
	log.Println("write to", certFile)
	ioutil.WriteFile(certFile, certSigned, 0777) // cert 写入文件

	certPrivateKeyFile := "cert.key"
	certPrivateKeyDER := x509.MarshalPKCS1PrivateKey(certPrivateKey) // 将私钥转换为 DER 编码格式
	log.Println("write to", certPrivateKeyFile)
	ioutil.WriteFile(certPrivateKeyFile, certPrivateKeyDER, 0777) // 私钥写入文件

	caTr, _ := x509.ParseCertificate(caSelfSigned)
	certTr, _ := x509.ParseCertificate(certSigned)
	err = certTr.CheckSignatureFrom(caTr) // 验证 c 上的签名是否是来自 parent 的有效签名
	log.Println("check signature", err)
}
