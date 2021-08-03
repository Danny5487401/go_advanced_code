// 中层封装：第二个版本

// 看源码http

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"
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

	// 方式一：添加中间件--缺点复杂
	//mux.Handle("/greeting", PanicRecover(WithLogger(Metric(greeting("welcome, dj")))))
	// 方式二:
	middlewares := []Middleware{
		PanicRecover,
		WithLogger,
		Metric,
	}
	mux.Handle("/greeting", applyMiddlewares(greeting("welcome, dj"), middlewares...))

	err = http.ListenAndServe(":8080", mux) //使用自定义的路由器mux时，用http包的ListenAndServe函数，此时要传入mux
	if err != nil {
		log.Fatal(err)
	}

}

// 我们当然可以直接定义一个实现Handler接口的类型，然后注册该类型的实例：
//type Handler interface {
//	ServeHTTP(ResponseWriter, *Request)
//}
type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "URL:"+r.Host+r.URL.String())
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello Danny, this is version 2")
}

// PanicRecover 添加中间件
func PanicRecover(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(string(debug.Stack()))
			}
		}()

		handler.ServeHTTP(w, r)
	})
}

func Metric(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			fmt.Printf("path:%s elapsed:%fs\n", r.URL.Path, time.Since(start).Seconds())
		}()
		time.Sleep(1 * time.Second)
		handler.ServeHTTP(w, r)
	}
}

func WithLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("path:%s process start...\n", r.URL.Path)
		defer func() {
			fmt.Printf("path:%s process end...\n", r.URL.Path)
		}()
		handler.ServeHTTP(w, r)
	})
}

type greeting string

func (g greeting) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, g)
}

func applyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}

type Middleware func(handler http.Handler) http.Handler
