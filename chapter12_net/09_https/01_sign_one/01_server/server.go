package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	// 记录请求协议
	log.Printf("Got connection: %s, tls version :%v , tls suite: %v", r.Proto, r.TLS.Version, r.TLS.CipherSuite)
	// 向客户发送一条消息
	w.Write([]byte("Hello"))
}

func main() { /*修改TLS版本*/
	http.HandleFunc("/", handle)

	server := &http.Server{
		Addr:    ":443",
		Handler: nil,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS11,
			MaxVersion: tls.VersionTLS13,
		},
	}

	err := server.ListenAndServeTLS("chapter17_dataStructure_n_algorithm/06_certificate/02_x509/cert.pem", "chapter17_dataStructure_n_algorithm/06_certificate/02_x509/key.pem")
	if err != nil {
		fmt.Println("error:", err)
	}

}
