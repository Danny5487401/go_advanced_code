package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", emptyHandler)
	mux.HandleFunc("/login", emptyHandler)
	mux.HandleFunc("/logout", emptyHandler)
	mux.HandleFunc("/products", emptyHandler)
	mux.HandleFunc("/product/:productID", emptyHandler)
	mux.HandleFunc("/basket", emptyHandler)
	mux.HandleFunc("/about", emptyHandler)
	http.ListenAndServe(":8080", mux)
}

var emptyHandler = func(writer http.ResponseWriter, request *http.Request) {

}
