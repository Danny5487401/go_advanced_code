package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("收到%s请求\n", r.RemoteAddr)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "hello world")
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
	http.HandleFunc("/hello", HelloWorld)

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
		fmt.Println(err)
	}

}
