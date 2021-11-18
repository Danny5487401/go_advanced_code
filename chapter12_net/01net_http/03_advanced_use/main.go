package main

// 低层封装：第三个版本

import (
	"io"
	"log"
	"net/http"
	"time"
)

//通过map保存注册的handler
//然后通过底层的serveHTTP进行转发，这是效率最高的，因为没有进任何封装
//参考beego的那些路由

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	server := http.Server{
		Addr:              ":8080",
		Handler:           &myHandler{}, //注册
		ReadHeaderTimeout: 5 * time.Second,
	}

	// 根据路由前缀注册handler
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/hello"] = sayHello
	mux["/bye"] = sayBye

	err := server.ListenAndServe() //使用自定义的map来实现路由时，使用ListenAndServe方法，上面用的是ListenAndServe函数。
	if err != nil {
		log.Fatal(err)
	}
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r) //根据URL从map中取出函数名，然后调用。
		return
	}

	io.WriteString(w, "Version 3 "+"URL: "+r.URL.String())
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Version 3 "+"Hello ")
}

func sayBye(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Version 3 "+"Bye")
}
