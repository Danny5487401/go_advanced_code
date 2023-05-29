<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [HTTP中间件](#http%E4%B8%AD%E9%97%B4%E4%BB%B6)
  - [优雅连接中间件：alice](#%E4%BC%98%E9%9B%85%E8%BF%9E%E6%8E%A5%E4%B8%AD%E9%97%B4%E4%BB%B6alice)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# HTTP中间件

当你构建一个web应用程序时,可能有一些共享的功能,你想参加许多(甚至是全部)HTTP请求。 您可能想要记录每个请求,gzip每个响应,或做一些繁重的处理之前检查缓存。


组织这个共享功能的一种方法是设置它 中间件 ——独立的代码独立作用于正常的应用程序请求之前或之后处理程序。 在一个共同的地方去使用ServeMux之间的中间件和应用程序处理程序,以便控制流为一个HTTP请求的样子:

```go
ServeMux => Middleware Handler => Application Handler
```

## 优雅连接中间件：alice

原来
```go
http.Handle("/", myLoggingHandler(authHandler(enforceXMLHandler(finalHandler))))
```

改造
```go
http.Handle("/", alice.New(myLoggingHandler, authHandler, enforceXMLHandler).Then(finalHandler))
```


Alice 真正的好处是,它允许您指定一个处理程序链，并重用它为多个路线。
```go

stdChain := alice.New(myLoggingHandler, authHandler, enforceXMLHandler)

http.Handle("/foo", stdChain.Then(fooHandler))
http.Handle("/bar", stdChain.Then(barHandler))
```