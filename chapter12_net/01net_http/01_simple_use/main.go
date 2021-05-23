package main

import (
	"fmt"
	"net/http"
)

/*
net/http代码块分为客户端和服务端两个部分。以下是net/http标准库的各个组成部分在客户端/服务中的划分：见structure.png

执行流程：process_diagram.png
服务端和服务端通信的过程：

1.服务端创建socket，绑定/监听指定的ip地址和端口，即Listen Socket
2. 客户端与Listen Socket连接，Listen Socket接受客户端的请求，得到Client Socket，接下来通过Client Socket与客户端通信
3. 创建go线程服务的一个连接，处理客户端的请求。首先从Client Socket读取HTTP请求的协议头，
	如果是POST方法，还可能要读取客户端提交的数据。然后交给相应的handler处理请求，handler处理完毕准备好客户端需要的数据, 通过Client Socket写给客户端
 */

// 简单写法--高级封装-第一个版本

func main()  {
	// 1. 注册函数
	http.HandleFunc("/hello",sayHello)

	// 2. 开启服务
	//ListenAndServe监听TCP地址addr，并且会使用handler参数调用Serve函数处理接收到的连接。
	//handler参数一般会设为nil，此时会使用DefaultServeMux。
	http.ListenAndServe(":8080",nil)

}

func sayHello(res http.ResponseWriter,req *http.Request)  {
	fmt.Fprintf(res, "hello,danny,\nreq= %+v\n",req)
}

