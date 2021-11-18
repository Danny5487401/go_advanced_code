package main

import (
	"fmt"
	"net/http"
)

/*
net/http代码块分为客户端和服务端两个部分。

服务端和服务端通信的过程：

1. 服务端创建socket，绑定/监听指定的ip地址和端口，即Listen Socket
2. 客户端与Listen Socket连接，Listen Socket接受客户端的请求，得到Client Socket，接下来通过Client Socket与客户端通信
3. 创建go线程服务的一个连接，处理客户端的请求。首先从Client Socket读取HTTP请求的协议头，
	如果是POST方法，还可能要读取客户端提交的数据。然后交给相应的handler处理请求，handler处理完毕准备好客户端需要的数据, 通过Client Socket写给客户端
*/

// 简单写法--高级封装-第一个版本

func main() {
	// 1. 注册函数
	http.HandleFunc("/hello", sayHello)
	/*
		func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
		  DefaultServeMux.HandleFunc(pattern, handler)
		}
			我们发现它直接调用了一个名为DefaultServeMux对象的HandleFunc()方法。
		type HandlerFunc func(ResponseWriter, *Request)

		func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
		  f(w, r)
		}
			HandlerFunc实际上是以函数类型func(ResponseWriter, *Request)为底层类型，为HandlerFunc类型定义了方法ServeHTTP。是的，Go 语言允许为（基于）函数的类型定义方法
	*/

	// 2. 开启服务
	//ListenAndServe监听TCP地址addr，并且会使用handler参数调用Serve函数处理接收到的连接。
	//handler参数一般会设为nil，此时会使用DefaultServeMux。
	http.ListenAndServe(":8080", nil)

}

/*
	*http.Request表示 HTTP 请求对象，该对象包含请求的所有信息，如 URL、首部、表单内容、请求的其他内容等
	http.ResponseWriter是一个接口类型:
		用于向客户端发送响应，实现了ResponseWriter接口的类型显然也实现了io.Writer接口。所以在处理函数index中，可以调用fmt.Fprintln()向ResponseWriter写入响应信息
*/
func sayHello(res http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(res, "hello,danny,\nreq= %+v\n", req)
}
