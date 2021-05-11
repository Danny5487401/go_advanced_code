// 中层封装：第二个版本

// 看源码http

package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

//对比第一版本：传入的handler是nil，接下来我们要把nil换掉，换成我们自己实现的handler
//其实并不是一个handler，真实是ServerMux（路由类型），首先需要实现一个mux

func main() {

	mux := http.NewServeMux() //实例化一个mux,返回一个mux
	//设置handler操作，不可以用HandlerFunc这个函数，这个函数用的是默认的路由器DefaultServeMux
	//这里用的是默认的的handler进行注册，要自己实现handler，注册到mux中

	// 注册方式一
	mux.Handle("/", &myHandler{})

	// 注册方式二
	mux.HandleFunc("/hello", sayHello)

	//实现文件服务器-简易静态文件实现
	wd, err := os.Getwd() //os.Getwd()返回一个对应当前工作目录的根路径
	if err != nil {
		log.Fatal(err)
	}
	// 注册方式三
	mux.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(wd))))

	err = http.ListenAndServe(":8080", mux) //使用自定义的路由器mux时，用http包的ListenAndServe函数，此时要传入mux
	if err != nil {
		log.Fatal(err)
	}

}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "URL:"+r.Host+r.URL.String())
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello Danny, this is version 2")
}
