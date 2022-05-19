package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// 读取签发服务端证书的CA证书
	CACrt, err := ioutil.ReadFile("chapter12_net/09_https/ca_pem/rootCA.pem")
	if err != nil {
		log.Fatalln(err)
	}
	// 创建用于信任校验服务端的CA证书池
	CAPool := x509.NewCertPool()
	// 添加CA证书到池（信任签发服务端证书的CA）
	CAPool.AppendCertsFromPEM(CACrt)

	// 读取客户端的证书，以及客户端私钥
	clientCrt, err := tls.LoadX509KeyPair("chapter12_net/09_https/client_pem/client.crt", "chapter12_net/09_https/client_pem/client.key")
	if err != nil {
		log.Fatalln(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			/* 设置服务端CA证书池，池里的CA证书，均为可信的 */
			RootCAs: CAPool,
			/* 设置要发送给服务端校验的客户端证书；和客户端自用私钥（不发给服务端） */
			Certificates: []tls.Certificate{
				clientCrt,
			},
		},
	}

	client := &http.Client{
		Transport: tr,
	}

	/* 协议：https，非http */
	resp, err := client.Get("https://danny-host:1280/")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%s\n", body)
}
