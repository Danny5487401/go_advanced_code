package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	// 记录请求协议
	log.Printf("Got connection: %s", r.Proto)
	// 向客户发送一条消息
	w.Write([]byte("Hello"))
}

//func main() {/*默认情况下*/
//	http.HandleFunc("/", HelloWorld)
//
//	err := http.ListenAndServeTLS(":1280", "./cert/server.crt", "./cert/server.key", nil)
//	if err != nil {
//		fmt.Println(err)
//	}
//}

func main() { /*修改TLS版本*/
	http.HandleFunc("/", handle)

	server := &http.Server{
		Addr:    ":443",
		Handler: nil,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS10,
			MaxVersion: tls.VersionTLS13,
		},
	}

	err := server.ListenAndServeTLS("chapter17_dataStructure_n_algorithm/06_certificate/02_x509/cert.pem", "chapter17_dataStructure_n_algorithm/06_certificate/02_x509/key.pem")
	if err != nil {
		fmt.Println("error:", err)
	}

}
