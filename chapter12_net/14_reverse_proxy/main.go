package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

/*
们需要一个 Web 服务器(http://localhost:1330)能够提供以下功能：

1. 获取到请求
2. 读取请求体，特别是 proxy_condition 字段（也叫做代理域）
3. 如果 proxy_condition 字段的值为 A，则转发到 http://localhost:1331
4. 如果 proxy_condition 字段的值为 B，则转发到 http://localhost:1332
5. 否则，则转发到默认的 URL (http://localhost:1333)
*/

const PORT = "1330"
const A_CONDITION_URL = "http://localhost:1331"
const B_CONDITION_URL = "http://localhost:1332"
const DEFAULT_CONDITION_URL = "http://localhost:1333"

type requestPayloadStruct struct {
	ProxyCondition string `json:"proxy_condition"`
}

// Get the port to listen on
func getListenAddress() string {
	return ":" + PORT
}

// Log the env variables required for a reverse proxy
func logSetup() {

	log.Printf("Server will run on: %s\n", getListenAddress())
	log.Printf("Redirecting to A url: %s\n", A_CONDITION_URL)
	log.Printf("Redirecting to B url: %s\n", B_CONDITION_URL)
	log.Printf("Redirecting to Default url: %s\n", DEFAULT_CONDITION_URL)
}

// Get a json decoder for a given requests body
func requestBodyDecoder(request *http.Request) *json.Decoder {
	// Read body to buffer
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		panic(err)
	}

	// Because golang is a pain in the ass if you read the body then any susequent calls
	// are unable to read the body again....
	request.Body = io.NopCloser(bytes.NewBuffer(body))

	return json.NewDecoder(io.NopCloser(bytes.NewBuffer(body)))
}

// Parse the requests body
func parseRequestBody(request *http.Request) requestPayloadStruct {
	decoder := requestBodyDecoder(request)

	var requestPayload requestPayloadStruct
	err := decoder.Decode(&requestPayload)

	if err != nil {
		panic(err)
	}

	return requestPayload
}

// Log the typeform payload and redirect url
func logRequestPayload(requestPayload requestPayloadStruct, proxyUrl string) {
	log.Printf("proxy_condition: %s, proxy_url: %s\n", requestPayload.ProxyCondition, proxyUrl)
}

// Get the url for a given proxy condition
func getProxyUrl(proxyConditionRaw string) string {
	proxyCondition := strings.ToUpper(proxyConditionRaw)

	if proxyCondition == "A" {
		return A_CONDITION_URL
	}

	if proxyCondition == "B" {
		return B_CONDITION_URL
	}

	return DEFAULT_CONDITION_URL
}

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	urlInfo, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(urlInfo)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req)
	}

	// Update the headers to allow for SSL redirection
	// req.URL.Host = url.Host
	// req.URL.Scheme = url.Scheme
	// req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	// req.Host = url.Host
	proxy.ModifyResponse = modifyResponse

	// Note that ServeHttp is nonblocking and uses a goroutine under the hood
	proxy.ServeHTTP(res, req)
}

func modifyRequest(req *http.Request) {
	req.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
}

func modifyResponse(resp *http.Response) error {
	resp.Header.Set("X-Proxy", "Magical")
	return nil
}

// Given a request send it to the appropriate url
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	requestPayload := parseRequestBody(req)
	urlInfo := getProxyUrl(requestPayload.ProxyCondition)

	logRequestPayload(requestPayload, urlInfo)

	serveReverseProxy(urlInfo, res, req)
}

func main() {
	// Log setup values
	logSetup()

	// start server
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(getListenAddress(), nil); err != nil {
		panic(err)
	}
}

/*
✗ pnpm install http-server
✗ http-server -p 1331
✗ http-server -p 1332
✗ http-server -p 1333
✗ go run main.go

✗ curl --request GET \
  --url http://localhost:1330/ \
  --header 'content-type: application/json' \
  --data '{
    "proxy_condition": "a"
  }'
*/
