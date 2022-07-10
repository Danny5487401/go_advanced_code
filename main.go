package main

import (
	"log"
	"net/http"
)

func main() {

	err := http.ListenAndServe(":8080", http.FileServer(http.Dir("/Users/python/Desktop/go_advanced_code")))
	if err != nil {
		log.Fatal("...")

	}
}
