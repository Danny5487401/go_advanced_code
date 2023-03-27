package main

import (
	"crypto/tls"
	"log"
	"net/http"

	// golang.org/x/crypto/acme/autocert，它是由Go的核心开发人员开发的
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("tencent.danny.games"),
		Cache:      autocert.DirCache("certs"), // 证书暂存在certs文件夹。autocert会定期自动刷新，避免证书过期
	}

	http.HandleFunc("/", handle)
	server := &http.Server{
		Addr: ":443",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))

	log.Fatal(server.ListenAndServeTLS("", ""))

}

func handle(w http.ResponseWriter, r *http.Request) {
	// 记录请求协议
	log.Printf("Got connection: %s", r.Proto)
	// 向客户发送一条消息
	w.Write([]byte("Hello"))
}
