package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HTTP/1.1 with TLS/SSL"))
	})
	log.Fatal(http.ListenAndServeTLS(":1280", "chapter12_net/09_https/openssl_conf/zchd.crt", "chapter12_net/09_https/openssl_conf/ca.key", nil))
}
