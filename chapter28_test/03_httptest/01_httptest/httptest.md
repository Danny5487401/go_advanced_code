# httptest


## httptest 方法介绍

```go
// NewRequest 方法用来创建一个 http 的请求体。
func NewRequest(method, target string, body io.Reader) *http.Request


// NewRecorder(响应体)
func NewRecorder() *ResponseRecorder
```


