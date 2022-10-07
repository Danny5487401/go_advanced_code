package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			/* 不校验服务端证书，直接信任 */
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: tr,
	}

	/* 协议：https，非http */
	resp, err := client.Get("https://127.0.0.1:443/hello")
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
