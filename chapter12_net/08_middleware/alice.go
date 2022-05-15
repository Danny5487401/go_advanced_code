package main

import (
	"fmt"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"time"
)

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started,请求方式 %s 路径%s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing index handler")
	fmt.Fprintf(w, "welcome!")
}

func middlewareFirst(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("MiddlewareFirst - before Handler")
		next.ServeHTTP(w, r)
		log.Println("MiddlewareFirst - after Handler")
	})
}

func middlewareSecond(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("MiddlewareSecond - before Handler")
		if r.URL.Path == "/message" {
			if r.URL.Query().Get("password") == "123" {
				log.Println("Authorized to system...")
				next.ServeHTTP(w, r)
			} else {
				log.Println("Failed to authorize to the system")
				return
			}
		} else {
			next.ServeHTTP(w, r)
		}
		log.Println("MiddlewareSecond - after Handler")
	})
}
func message(w http.ResponseWriter, r *http.Request) {
	log.Println("Executed message Handler")
	fmt.Fprintf(w, "message...")
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Executed foo Handler")
	fmt.Fprintf(w, "foo...")
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Executed bar Handler")
	fmt.Fprintf(w, "bar...")
}

func main() {
	// 1. 单个中间件
	indexHandler := http.HandlerFunc(index)
	http.Handle("/", loggingHandler(indexHandler))

	// 2. 多个中间件
	messageHandler := http.HandlerFunc(message)
	http.Handle("/message", middlewareFirst(middlewareSecond(messageHandler)))

	// 3. 优化写法
	stdChain := alice.New(middlewareFirst, middlewareSecond)
	http.Handle("/foo", stdChain.Then(http.HandlerFunc(fooHandler)))
	http.Handle("/bar", stdChain.Then(http.HandlerFunc(barHandler)))

	server := &http.Server{
		Addr: ":8080",
	}

	log.Println("Listening...")
	server.ListenAndServe()
}
