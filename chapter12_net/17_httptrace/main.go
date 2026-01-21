package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"
)

func main() {

	// Prepare request with trace attached to it.
	// httpbin这个网站能测试 HTTP 请求和响应的各种信息，比如 cookie、ip、headers 和登录验证等，且支持 GET、POST 等多种方法，对 web 开发和测试很有帮助。
	// 模拟get 请求
	req, err := http.NewRequest(http.MethodGet, "https://httpbin.org/get", nil)
	if err != nil {
		log.Fatalln("request error", err)
	}

	/*
		Go 官方设计是：

		Request 在 Transport / Middleware / Redirect 中 可能被并发读取

		为了安全：

		- 不能原地改

		- 所有修改 API 都返回新对象
	*/
	// 注意接收返回值
	req = traceHTTPRequest(req)
	// Make a request.
	res, err := client().Do(req)
	if err != nil {
		log.Fatalln("client error", err)
	}
	defer res.Body.Close()

	// Read response.
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}

func client() *http.Client {
	return &http.Client{
		Transport: transport(),
	}
}

func transport() *http.Transport {
	return &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig:   tlsConfig(),
	}
}

func tlsConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}

func traceHTTPRequest(req *http.Request) *http.Request {
	start := time.Now()
	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			fmt.Printf("HTTP Trace: DNS resolution started for host %s \n", info.Host)
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			if info.Err != nil {
				fmt.Printf("HTTP Trace: DNS resolution failed: %v \n ", info.Err)
			} else {
				fmt.Printf("HTTP Trace: DNS resolution completed in %v, addresses: %v \n", time.Since(start), info.Addrs)
			}
		},
		ConnectStart: func(network, addr string) {
			fmt.Printf("HTTP Trace: TCP connection started to %s://%s  \n ", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			if err != nil {
				fmt.Printf("HTTP Trace: TCP connection failed: %v \n", err)
			} else {
				fmt.Printf("HTTP Trace: TCP connection completed in %v \n", time.Since(start))
			}
		},
		GotConn: func(info httptrace.GotConnInfo) {
			fmt.Printf("HTTP Trace: Got connection: reused=%v, local=%v, remote=%v \n",
				info.Reused, info.Conn.LocalAddr(), info.Conn.RemoteAddr())
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			if info.Err != nil {
				fmt.Printf("HTTP Trace: Wrote request failed: %v \n", info.Err)
			} else {
				fmt.Printf("HTTP Trace: Request written in %v \n", time.Since(start))
			}
		},
		GotFirstResponseByte: func() {
			fmt.Printf("HTTP Trace: First response byte received in %v \n", time.Since(start))
		},
	}

	// 把 ClientTrace 放进 context.Context, 然后新建一个 request 对象
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	return req
}
