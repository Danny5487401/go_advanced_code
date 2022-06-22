package main

import (
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example using https in golang.\n"))
}

func main() {
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServeTLS(":443", "chapter12_net/09_https/ca_pem/cert.pem", "chapter12_net/09_https/ca_pem/key.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
