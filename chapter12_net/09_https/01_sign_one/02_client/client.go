package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	// 方式一： 不校验服务端证书，直接信任
	//tr := &http.Transport{
	//	TLSClientConfig: &tls.Config{
	//		/* 不校验服务端证书，直接信任 */
	//		InsecureSkipVerify: true,
	//	},
	//}
	// 方式二：添加CA证书到池
	// 读取CA证书
	CACrt, err := os.ReadFile("chapter17_dataStructure_n_algorithm/06_certificate/02_x509/cert.pem")
	if err != nil {
		log.Fatalln(err)
	}
	// 创建CA证书池
	CAPool := x509.NewCertPool()
	// 添加CA证书到池
	CAPool.AppendCertsFromPEM(CACrt)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			/* 设置CA证书池，池里的CA证书，均为可信的 */
			RootCAs:    CAPool,
			MaxVersion: tls.VersionTLS13,
		},
	}

	client := &http.Client{
		Transport: tr,
	}

	/* 协议：https，非http */
	resp, err := client.Get("https://localhost:443/")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("body:%s\n", body)
}
