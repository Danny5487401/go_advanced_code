<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [责任链模式](#%E8%B4%A3%E4%BB%BB%E9%93%BE%E6%A8%A1%E5%BC%8F)
  - [意图](#%E6%84%8F%E5%9B%BE)
  - [主要解决](#%E4%B8%BB%E8%A6%81%E8%A7%A3%E5%86%B3)
  - [何时使用](#%E4%BD%95%E6%97%B6%E4%BD%BF%E7%94%A8)
  - [如何解决](#%E5%A6%82%E4%BD%95%E8%A7%A3%E5%86%B3)
  - [优点](#%E4%BC%98%E7%82%B9)
  - [缺点](#%E7%BC%BA%E7%82%B9)
  - [使用场景](#%E4%BD%BF%E7%94%A8%E5%9C%BA%E6%99%AF)
  - [应用实例](#%E5%BA%94%E7%94%A8%E5%AE%9E%E4%BE%8B)
    - [gin Context](#gin-context)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 责任链模式
![](example.png)

![](process.png)
![](process2.png)

一种处理请求的模式，它让多个处理器都有机会处理该请求，直到其中某个处理成功为止。责任链模式把多个处理器串成链，然后让请求在链上传递

## 意图
避免请求发送者与接受者耦合在一起，让多个对象都有可能接收请求，将这些对象连接成一条链，并且沿着这条链传递请求，直到有对象处理它为止。

## 主要解决
职责链上的处理者负责处理请求，客户只需要将请求发送到职责链上即可，无须关心请求的处理细节和请求的传递，所以职责链将请求的发送者和请求的处理者解偶了。

## 何时使用
在处理消息的时候以过滤很多道。

## 如何解决
拦截的类都实现统一接口。

## 优点

1. 降低耦合度。它将请求的发送者和接收者解耦。
2. 简化了对象。使得对象不需要知道链的结构。
3. 增强给对象指派职责的灵活性。通过改变链内的成员或者调动它们的次序，允许动态地新增或者删除责任。
4. 增加新的请求处理类很方便。


## 缺点

1. 不能保证请求一定被接收。
2. 系统性能将受到一定影响，而且在进行代码调试时不太方便，可能会造成循环调用。
3. 可能不容易观察运行时的特征，有碍于除错。


## 使用场景

1. 有多个对象可以处理同一个请求，具体哪个对象处理该请求由运行时刻自动确定。
2. 在不明确指定接受者的情况下，向多个对象中的一个提交请求。
3. 可动态指定一组对象处理请求。

## 应用实例
1. 红楼梦中的"击鼓传花"。 
2. JS 中的事件冒泡。 
3. JAVA WEB 中 Apache Tomcat 对 Encoding 的处理，Struts2 的拦截器，jsp servlet 的 Filter。


### gin Context 

Next() 会按照顺序将一个个中间件执行完毕
```go
// github.com/gin-gonic/gin@v1.9.1/context.go
type Context struct {
	// ...
	handlers HandlersChain
	// ...
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}
```

```go
// github.com/gin-gonic/gin@v1.9.1/gin.go
// HandlersChain defines a HandlerFunc slice.
type HandlersChain []HandlerFunc
```



