package main

import (
	"log"
	"net/http"
)

func main() {
	// 文件服务器
	err := http.ListenAndServe(":8080", http.FileServer(http.Dir("")))
	if err != nil {
		log.Fatal("...")
	}
}
