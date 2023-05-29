<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [httptest](#httptest)
  - [httptest 方法介绍](#httptest-%E6%96%B9%E6%B3%95%E4%BB%8B%E7%BB%8D)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# httptest


## httptest 方法介绍

```go
// NewRequest 方法用来创建一个 http 的请求体。
func NewRequest(method, target string, body io.Reader) *http.Request


// NewRecorder(响应体)
func NewRecorder() *ResponseRecorder
```


