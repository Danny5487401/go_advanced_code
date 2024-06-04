<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [httptest](#httptest)
  - [httptest 方法介绍](#httptest-%E6%96%B9%E6%B3%95%E4%BB%8B%E7%BB%8D)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# httptest
Web应用程序中往往需要与其他系统进行交互，比如通过http访问其他系统，此时就需要有一种方法用于打桩来模拟Web服务器和客户端，httptest包即Go语言针对Web应用提供的解决方案。

httptest可以方便的模拟各种Web服务器和客户端，以达到测试目的

## httptest 方法介绍

```go
// NewRequest 方法用来创建一个 http 的请求体。
func NewRequest(method, target string, body io.Reader) *http.Request


// NewRecorder(响应体)
func NewRecorder() *ResponseRecorder
```


