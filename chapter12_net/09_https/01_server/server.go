package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// 读取签发客户端证书的CA证书
	CACrt, err := ioutil.ReadFile("chapter12_net/09_https/ca_pem/rootCA.pem")
	if err != nil {
		log.Fatalln(err)
	}
	// 创建用于信任校验客户端的CA证书池
	CAPool := x509.NewCertPool()
	// 添加CA证书到池（信任签发客户端证书的CA）
	CAPool.AppendCertsFromPEM(CACrt)

	server := &http.Server{
		Addr:    ":1280",
		Handler: &myHandler{},
		TLSConfig: &tls.Config{
			// 设置服务端对客户端的CA信任池
			ClientCAs: CAPool,
			// 设置TLS/SSL鉴权认证时要求客户端提供其证书并校验
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}

	// 设置：
	// * 服务端要发送给客户端认证的服务端证书
	// * 服务端自己的私钥，用以和客户端协商对称加解密密钥时的加解密
	log.Fatal(server.ListenAndServeTLS("chapter12_net/09_https/server_pem/server.crt", "chapter12_net/09_https/server_pem/server.key"))
}

type myHandler struct {
}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HTTP/1.1 TLS/SSL 双向校验"))
}
