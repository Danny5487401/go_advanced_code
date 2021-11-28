package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

func main() {
	getRes()
	postRes()

}

// 简单的get请求
func getRes() {
	url := `http://httpbin.org/get`

	status, resp, err := fasthttp.Get(nil, url)
	if err != nil {
		fmt.Println("请求失败:", err.Error())
		return
	}

	if status != fasthttp.StatusOK {
		fmt.Println("请求没有成功:", status)
		return
	}

	fmt.Println(string(resp))
}

// 简单的post请求
func postRes() {
	url := `http://httpbin.org/post?key=123`

	// 填充表单，类似于net/url
	args := &fasthttp.Args{}
	args.Add("name", "test")
	args.Add("age", "18")

	status, resp, err := fasthttp.Post(nil, url, args)
	if err != nil {
		fmt.Println("请求失败:", err.Error())
		return
	}

	if status != fasthttp.StatusOK {
		fmt.Println("请求没有成功:", status)
		return
	}

	fmt.Println(string(resp))
}

// 在有性能要求的代码中，通过 AcquireRequest 和 AcquireResponse 来获取 req 和 resp
func customizedReq() {
	url := `http://httpbin.org/post?key=123`

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 用完需要释放资源
	req.SetRequestURI(url)

	requestBody := []byte(`{"request":"test"}`)
	req.SetBody(requestBody)

	// 默认是application/x-www-form-urlencoded
	req.Header.SetContentType("application/json")
	req.Header.SetMethod("POST")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源

	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		fmt.Println("请求失败:", err.Error())
		return
	}

	b := resp.Body()

	fmt.Println("result:\r\n", string(b))
}
