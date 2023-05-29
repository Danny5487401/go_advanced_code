<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Goroutine泄漏](#goroutine%E6%B3%84%E6%BC%8F)
  - [场景](#%E5%9C%BA%E6%99%AF)
  - [泄露情况分类](#%E6%B3%84%E9%9C%B2%E6%83%85%E5%86%B5%E5%88%86%E7%B1%BB)
    - [1. 发送不接收](#1-%E5%8F%91%E9%80%81%E4%B8%8D%E6%8E%A5%E6%94%B6)
  - [检测工具 goleak](#%E6%A3%80%E6%B5%8B%E5%B7%A5%E5%85%B7-goleak)
    - [goleak 的实现原理](#goleak-%E7%9A%84%E5%AE%9E%E7%8E%B0%E5%8E%9F%E7%90%86)
  - [参考资料](#%E5%8F%82%E8%80%83%E8%B5%84%E6%96%99)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Goroutine泄漏
谈到内存管理，go语言为我们处理好了大量的细节。Go语言编译器使用逃逸分析（escape analysis）来决变量的存储。
运行时通过使用垃圾回收器跟踪和管理堆分配。虽然在应用程序中不产生内存泄漏并不可能，但几率大大降低。

一种常见的内存泄漏就是goroutines泄漏。如果你启动了一个goroutine，你期望它最终会终止，但是它并没有，这就是goroutine泄漏了。
它在应用程序的运行期内存在, 分配给该goroutine的内存无法释放。这是“永远不要在不知道它将如何停止的情况下启动一个goroutine”的建议背后的部分理由

## 场景
如果这段代码在常驻服务中执行，比如 http server，每接收到一个请求，便会启动一次 sayHello，时间流逝，每次启动的 goroutine 都得不到释放，你的服务将会离奔溃越来越近

## 泄露情况分类
按照并发的数据同步方式对泄露的各种情况进行分析。简单可归于两类，即：
1. channel 导致的泄露
2. 传统同步机制导致的泄露

传统同步机制主要指面向共享内存的同步机制，比如排它锁、共享锁等。这两种情况导致的泄露还是比较常见的。
但是go 由于 defer 的存在，第二类情况，一般情况下还是比较容易避免的。

channel 引起的泄露分析:
### 1. 发送不接收
理想情况下，我们希望接收者总能接收完所有发送的数据，这样就不会有任何问题。但现实是，一旦接收者发生异常退出，停止继续接收上游数据，发送者就会被阻塞

原因：当接收者停止工作，发送者并不知道，还在傻傻地向下游发送数据。故而，我们需要一种机制去通知发送者。

做法：Go 可以通过 channel 的关闭向所有的接收者发送广播信息。


## 检测工具 goleak

goleak主要关注两个方法即可：VerifyNone、VerifyTestMain
- VerifyNone用于单一测试用例中测试
- VerifyTestMain可以在TestMain中添加，可以减少对测试代码的入侵

### goleak 的实现原理
使用 runtime.Stack() 方法获取当前运行的所有goroutine的栈信息，默认定义不需要检测的过滤项，默认定义检测次数+检测间隔，不断周期进行检测，最终在多次检查后仍没有找到剩下的goroutine则判断没有发生goroutine泄漏


## 参考资料
1. [goleak官网](https://github.com/uber-go/goleak)
2. [uber-go/goleak检查goroutine泄漏原理](https://mp.weixin.qq.com/s/PGcutKTQZy3v9ln31dFvRg)

