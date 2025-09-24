# ***高级Goland学习代码*** _github.com/Danny5487401/go_advanced_code_
![](https://changkun.de/urlstat?mode=github&repo=)
[![Go Report Card](https://goreportcard.com/badge/github.com/talkgo/night?style=flat-square)](https://goreportcard.com/report/github.com/Danny5487401/github.com/Danny5487401/go_advanced_code)
[![GitHub stars](https://img.shields.io/github/stars/talkgo/night.svg?label=Stars&style=flat-square)](https://github.com/Danny5487401/github.com/Danny5487401/go_advanced_code)
[![GitHub forks](https://img.shields.io/github/forks/talkgo/night.svg?label=Fork&style=flat-square)](https://github.com/Danny5487401/github.com/Danny5487401/go_advanced_code)
![](https://img.shields.io/github/contributors/talkgo/night.svg?style=flat-square&color=orange&label=all%20contributors)
[![Documentation](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](http://godoc.org/github.com/Danny5487401/github.com/Danny5487401/go_advanced_code)
[![GitHub issues](https://img.shields.io/github/issues/talkgo/night.svg?label=Issue&style=flat-square)](https://github.com/Danny5487401/github.com/Danny5487401/go_advanced_code/issues)
![](https://changkun.de/urlstat?mode=github&repo=Danny5487401/github.com/Danny5487401/go_advanced_code)
[![license](https://img.shields.io/github/license/talkgo/night.svg?style=flat-square)](https://github.com/Danny5487401/github.com/Danny5487401/go_advanced_code/blob/master/LICENSE)

![高级go编程](.assets/logo/golang.jpeg)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [*目录*](#%E7%9B%AE%E5%BD%95)
  - [必备知识:](#%E5%BF%85%E5%A4%87%E7%9F%A5%E8%AF%86)
  - [第一章 I/O](#%E7%AC%AC%E4%B8%80%E7%AB%A0-io)
  - [第二章 协程Goroutine](#%E7%AC%AC%E4%BA%8C%E7%AB%A0-%E5%8D%8F%E7%A8%8Bgoroutine)
  - [第三章 通道Channel](#%E7%AC%AC%E4%B8%89%E7%AB%A0-%E9%80%9A%E9%81%93channel)
  - [第四章 interface 和反射](#%E7%AC%AC%E5%9B%9B%E7%AB%A0-interface-%E5%92%8C%E5%8F%8D%E5%B0%84)
  - [第五章 切片和数组](#%E7%AC%AC%E4%BA%94%E7%AB%A0-%E5%88%87%E7%89%87%E5%92%8C%E6%95%B0%E7%BB%84)
  - [第六章 指针](#%E7%AC%AC%E5%85%AD%E7%AB%A0-%E6%8C%87%E9%92%88)
  - [第七章 系统调用](#%E7%AC%AC%E4%B8%83%E7%AB%A0-%E7%B3%BB%E7%BB%9F%E8%B0%83%E7%94%A8)
  - [第八章 defer函数及汇编语言理解](#%E7%AC%AC%E5%85%AB%E7%AB%A0-defer%E5%87%BD%E6%95%B0%E5%8F%8A%E6%B1%87%E7%BC%96%E8%AF%AD%E8%A8%80%E7%90%86%E8%A7%A3)
  - [第九章 设计模式-OOP七大准则](#%E7%AC%AC%E4%B9%9D%E7%AB%A0-%E8%AE%BE%E8%AE%A1%E6%A8%A1%E5%BC%8F-oop%E4%B8%83%E5%A4%A7%E5%87%86%E5%88%99)
  - [第十章 函数式编程](#%E7%AC%AC%E5%8D%81%E7%AB%A0-%E5%87%BD%E6%95%B0%E5%BC%8F%E7%BC%96%E7%A8%8B)
  - [第十一章 汇编理解go语言底层源码(AMD芯片运行代码)](#%E7%AC%AC%E5%8D%81%E4%B8%80%E7%AB%A0-%E6%B1%87%E7%BC%96%E7%90%86%E8%A7%A3go%E8%AF%AD%E8%A8%80%E5%BA%95%E5%B1%82%E6%BA%90%E7%A0%81amd%E8%8A%AF%E7%89%87%E8%BF%90%E8%A1%8C%E4%BB%A3%E7%A0%81)
  - [第十二章 net 网络--涉及性能指标,协议栈统计,套接字信息](#%E7%AC%AC%E5%8D%81%E4%BA%8C%E7%AB%A0-net-%E7%BD%91%E7%BB%9C--%E6%B6%89%E5%8F%8A%E6%80%A7%E8%83%BD%E6%8C%87%E6%A0%87%E5%8D%8F%E8%AE%AE%E6%A0%88%E7%BB%9F%E8%AE%A1%E5%A5%97%E6%8E%A5%E5%AD%97%E4%BF%A1%E6%81%AF)
  - [第十三章 CGO调用C语言](#%E7%AC%AC%E5%8D%81%E4%B8%89%E7%AB%A0-cgo%E8%B0%83%E7%94%A8c%E8%AF%AD%E8%A8%80)
  - [第十四章 Context上下文-源码分析涉及父类EmptyCtx](#%E7%AC%AC%E5%8D%81%E5%9B%9B%E7%AB%A0-context%E4%B8%8A%E4%B8%8B%E6%96%87-%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90%E6%B6%89%E5%8F%8A%E7%88%B6%E7%B1%BBemptyctx)
  - [第十五章 接口嵌套编程](#%E7%AC%AC%E5%8D%81%E4%BA%94%E7%AB%A0-%E6%8E%A5%E5%8F%A3%E5%B5%8C%E5%A5%97%E7%BC%96%E7%A8%8B)
  - [第十六章 并发编程](#%E7%AC%AC%E5%8D%81%E5%85%AD%E7%AB%A0-%E5%B9%B6%E5%8F%91%E7%BC%96%E7%A8%8B)
  - [第十七章 数据结构及算法](#%E7%AC%AC%E5%8D%81%E4%B8%83%E7%AB%A0-%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84%E5%8F%8A%E7%AE%97%E6%B3%95)
  - [第十八章 错误跟踪 error 和 panic](#%E7%AC%AC%E5%8D%81%E5%85%AB%E7%AB%A0-%E9%94%99%E8%AF%AF%E8%B7%9F%E8%B8%AA-error-%E5%92%8C-panic)
  - [第十九章 nil 预定义标识](#%E7%AC%AC%E5%8D%81%E4%B9%9D%E7%AB%A0-nil-%E9%A2%84%E5%AE%9A%E4%B9%89%E6%A0%87%E8%AF%86)
  - [第二十章 for-range 源码分析](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E7%AB%A0-for-range-%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [第二十一章 time标准包源码分析](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E4%B8%80%E7%AB%A0-time%E6%A0%87%E5%87%86%E5%8C%85%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [第二十二章 数据驱动模板源码分析-->kratos工具](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E4%BA%8C%E7%AB%A0-%E6%95%B0%E6%8D%AE%E9%A9%B1%E5%8A%A8%E6%A8%A1%E6%9D%BF%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90--kratos%E5%B7%A5%E5%85%B7)
  - [第二十三章 调试内部对象](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E4%B8%89%E7%AB%A0-%E8%B0%83%E8%AF%95%E5%86%85%E9%83%A8%E5%AF%B9%E8%B1%A1)
  - [第二十四章 命令行参数解析](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E5%9B%9B%E7%AB%A0-%E5%91%BD%E4%BB%A4%E8%A1%8C%E5%8F%82%E6%95%B0%E8%A7%A3%E6%9E%90)
  - [第二十四章 Flag命令行参数及源码分析](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E5%9B%9B%E7%AB%A0-flag%E5%91%BD%E4%BB%A4%E8%A1%8C%E5%8F%82%E6%95%B0%E5%8F%8A%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [第二十五章 结构体类型方法](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E4%BA%94%E7%AB%A0-%E7%BB%93%E6%9E%84%E4%BD%93%E7%B1%BB%E5%9E%8B%E6%96%B9%E6%B3%95)
  - [第二十六章 strconv 字符串和数值型转换源码分析](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E5%85%AD%E7%AB%A0-strconv-%E5%AD%97%E7%AC%A6%E4%B8%B2%E5%92%8C%E6%95%B0%E5%80%BC%E5%9E%8B%E8%BD%AC%E6%8D%A2%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [第二十七章 image 图片处理](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E4%B8%83%E7%AB%A0-image-%E5%9B%BE%E7%89%87%E5%A4%84%E7%90%86)
  - [第二十八章 如何进行测试](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E5%85%AB%E7%AB%A0-%E5%A6%82%E4%BD%95%E8%BF%9B%E8%A1%8C%E6%B5%8B%E8%AF%95)
  - [第二十九章 module包管理](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E4%B9%9D%E7%AB%A0-module%E5%8C%85%E7%AE%A1%E7%90%86)
  - [第三十章 内存管理](#%E7%AC%AC%E4%B8%89%E5%8D%81%E7%AB%A0-%E5%86%85%E5%AD%98%E7%AE%A1%E7%90%86)
  - [第三十一章 go开发套件](#%E7%AC%AC%E4%B8%89%E5%8D%81%E4%B8%80%E7%AB%A0-go%E5%BC%80%E5%8F%91%E5%A5%97%E4%BB%B6)
  - [第三十二章 Generic 泛型](#%E7%AC%AC%E4%B8%89%E5%8D%81%E4%BA%8C%E7%AB%A0-generic-%E6%B3%9B%E5%9E%8B)
  - [第三十三章 makefile 使用](#%E7%AC%AC%E4%B8%89%E5%8D%81%E4%B8%89%E7%AB%A0-makefile-%E4%BD%BF%E7%94%A8)
  - [第三十四章 regexp 正则表达式 Regular Expression](#%E7%AC%AC%E4%B8%89%E5%8D%81%E5%9B%9B%E7%AB%A0-regexp-%E6%AD%A3%E5%88%99%E8%A1%A8%E8%BE%BE%E5%BC%8F-regular-expression)
  - [第三十五章 编码 Unicode](#%E7%AC%AC%E4%B8%89%E5%8D%81%E4%BA%94%E7%AB%A0-%E7%BC%96%E7%A0%81-unicode)
  - [第三十六章 unique 包--go 1.23](#%E7%AC%AC%E4%B8%89%E5%8D%81%E5%85%AD%E7%AB%A0-unique-%E5%8C%85--go-123)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->
# *目录*
Note: 目录同级为 *代码展示*，推荐在 Goland 版本 2022.2.1+ 运行,*推荐 GoVersion: 1.21+*


## 必备知识:
- [module包管理](chapter29_module/02_discipline/module.md)
- [golangci 规范并优化代码 + pre-commit工具](golangci.md)


## [第一章 I/O](chapter01_input_output/io.md)
- 1 os 操作系统模块
    - [1.1 os中 FileInfo 文件属性](chapter01_input_output/01_OS_module/01FileInfo/fileinfo.md)
    - [1.2 os文件操作](chapter01_input_output/01_OS_module/02FileOperation/main.go)   
    - 1.3 io包底层 Reader 和 Writer 接口   
      - 1.3.1 os,bytes,strings包   
    - [1.4 io 断点续传及网络支持](chapter01_input_output/01_OS_module/04seeker/resume.md) 
      - seeker 偏移量
      - http header: 客户端Range，Content-Range
    - [1.5 FilePath包 ](chapter01_input_output/01_OS_module/05filePath/walk.go)    
    - [1.5.1 walkPath遍历目录及文件(匹配文件名)](chapter01_input_output/01_OS_module/05filePath/walk.go)
- [2 bufio 缓存读写](chapter01_input_output/02_bufio/bufio.md)
  - [2.1 reader](chapter01_input_output/02_bufio/01reader/main.go)
  - [2.2 writer](chapter01_input_output/02_bufio/02writer/main.go)
- [3 Go 1.16 io.FS: Go 理解的文件系统](chapter01_input_output/03_io_fs/io_fs.md)
  - [go 1.16 前后的文件io对比](chapter01_input_output/03_io_fs/embed.go)
- [4 io.Pipe 对比使用 json.NewDecoder 流式解析 http body中间数据-->在 minio 下载数据实现](chapter01_input_output/04_io_pipe/io_pipe.md)

---
## 第二章 协程Goroutine
- [1 线程模型分类及Goroutine切换原则(GPM模型)](chapter02_goroutine/01_GPM/GPM.md)
    - [1.1 trace查看宏观调度流程(Goroutine启动时长)](chapter02_goroutine/01_GPM/01_trace/trace.md)
    - [1.2 使用 go.uber.org/automaxprocs 在容器里设置 GOMAXPROCS 的正确姿势](chapter02_goroutine/01_GPM/02_automaxprocs/GOMAXPROCS.md)
    - [1.3 runtime.LockOSThread() 将某个goroutine锁定到某个系统线程，这个线程只调度这个goroutine，进而可以被优先调度（相对其他goroutine）-->cni terway 中实现](chapter02_goroutine/01_GPM/03_os_thread_management/lockosthread.md)
- [2 runtime运行时模块](chapter02_goroutine/02_runtime/runtime.md)
    - [2.1 runtime核心功能及系统信息调用](chapter02_goroutine/02_runtime/01basic_use/main.go)
    - [2.2 Goexit()终止线程: defer 语句仍然执行](chapter02_goroutine/02_runtime/02GoExit/main.go)
    - [2.3 data race 资源竞争一致性问题分析](chapter02_goroutine/02_runtime/03ResourceCompetition/01problem/resource_competion.md)
        - [2.3.1 -race 标志分析问题产生](chapter02_goroutine/02_runtime/03ResourceCompetition/01problem/main.go)
        - [2.3.2 mutex解决问题](chapter02_goroutine/02_runtime/03ResourceCompetition/02Fix_Resource_data_consistency/main.go)
    - [2.4 监控代码性能 pprof](chapter02_goroutine/02_runtime/04_pprof/intro.md)
        - [2.4.1 标准包 runtime/pprof 及 net/http/pprof 使用](chapter02_goroutine/02_runtime/04_pprof/01_pprof/main.go)
        - [2.4.2 第三方包 github.com/pkg/profile ](chapter02_goroutine/02_runtime/04_pprof/02_pkg_profile/pkg_profile.md)
        - [2.4.3 dlsniper/debugger: 添加 pprof 标签调试 goroutine](chapter02_goroutine/02_runtime/04_pprof/03_pprof_label/client/main.go)
- [3 多 goroutine 的缓存一致性(涉及cpu伪共享)](chapter02_goroutine/03_cache/cache.md)
- [4 线程池(池化技术)](chapter02_goroutine/04_concurrent_pool/pool.md)
    - [4.1 使用channel实现Goroutine最大数量限制(令牌桶方式)](chapter02_goroutine/04_concurrent_pool/01_goroutine_max_control/main.go)
    - [4.2 百万请求处理案例](chapter02_goroutine/04_concurrent_pool/02_millionRequests/main.go)
    - [4.3 第三方包线程池ants](chapter02_goroutine/04_concurrent_pool/03_antsPool/ants.md)
    - [4.4 标准库连接池database/sql源码分析](chapter02_goroutine/04_concurrent_pool/04_database_sql/sql.md)
      - [4.4.1 连接池Benchmark对比](chapter02_goroutine/04_concurrent_pool/04_database_sql/database_pool_test.go)
- [5 channel导致goroutine泄漏分析及处理](chapter02_goroutine/05_goroutine_leaks/goroutine_leak.md)
  - [5.1 channel未正常关闭导致goroutine泄漏-->使用 goleak 工具检查](chapter02_goroutine/05_goroutine_leaks/01_leaks_happen/goroutine_leak_test.go)
  - [5.2 channel监听避免goroutine泄漏](chapter02_goroutine/05_goroutine_leaks/02_avoid_leaks/main.go)
- [6 Go routine 编排框架：oklog/run 包: 将各个组件作为一个整体运行，并有序地结束-->vault 应用](chapter02_goroutine/06_oklog_run/oklog_run.md)
---

## 第三章 通道Channel
- [1 Channel内部结构及源码分析(含PPT分析)](chapter03_channel/01_channel_use/channel.md)
    - [1.1 channel 初始化及引用传递](chapter03_channel/01_channel_use/00introdution/main.go)
      - [1.1.1 无缓存 channel](chapter03_channel/01_channel_use/01_initialize/01_unbuffered_channel/main.go)
      - [1.1.2 有缓冲 channel](chapter03_channel/01_channel_use/01_initialize/02_bufferChan/main.go)
      - [1.1.3 chanx: 使用 RingBuffer 实现无限缓存 channel](chapter03_channel/01_channel_use/01_initialize/03_unbounded_chan/unbounder_chan.md)
    - [1.2 使用 channel 实现 goroutine 父子通信](chapter03_channel/01_channel_use/02_parent_children_communication/main.go)
    - [1.3 死锁：range 未关闭的 channel](chapter03_channel/01_channel_use/03_deadlock/main.go)
    - [1.4 通道遍历:for range 语法 ](chapter03_channel/01_channel_use/05ChannelRange/main.go)
    - [1.5 优雅关闭 channel 与粗暴关闭 channel](chapter03_channel/01_channel_use/04channelClose/ChanClose.md)
      - [1.5.1 SPMC(Single-Producer Multi-Consumer 1 个 sender，N 个 receiver): 发送者通过关闭数据通道说 「不要再发送了」](chapter03_channel/01_channel_use/05_channel_close/case1_1sender_nreceiver/main.go)
      - 1.5.2 SPSC(Single-Producer Single-Consumer 1 个 sender，1 个 receiver):发送者通过关闭数据通道说 「不要再发送了」
      - [1.5.3 MPSC(Multi-Producer Single-Consumer N 个 sender，1 个 receiver): 接收者通过关闭一个信号通道说 「请不要再发送数据了」](chapter03_channel/01_channel_use/05_channel_close/case3_nsender_1receiver/nsender_1receiver.go)
      - [1.5.4 MPMC(Multi-Producer Multi-Consumer N 个 sender，M 个 receiver): 任意一个通过通知一个主持人去关闭一个信号通道说「让我们结束这场游戏吧」 ](chapter03_channel/01_channel_use/05_channel_close/case4_nsender_nreceiver/nsender_nreceiver.go)
    - [1.6 单向与双向通道](chapter03_channel/01_channel_use/06_single-directional_and_bi-directional_chan/main.go)
    - [1.7 读取 nil channel 实现阻塞](chapter03_channel/01_channel_use/07_read_nil_channel/main.go)
    - [1.8 使用 channel 传递 channel](chapter03_channel/01_channel_use/08_chan_pass_chan/main.go)
    - [1.9 循环读取关闭的通道值是否阻塞](chapter03_channel/01_channel_use/09_read_closed_chan/readCloseChan.go)
    - [1.10 select 实现 channel 优先级-->k8s中Node 的更新操作优先于 Pod 的更新](chapter03_channel/01_channel_use/10_priority_channel/priority_chan.md)
- [2 channel应用:TimerChan模块源码分析及使用陷阱](chapter03_channel/02_TimerChan/timer.md)
    - [2.1 reset重新等待被触发](chapter03_channel/02_TimerChan/01_TimerReset/timer_reset.md)
    - [2.2 timerStop使用](chapter03_channel/02_TimerChan/02_TimerStop/timer_stop.md)
    - [2.3 TimerAfter给数据库操作增加超时](chapter03_channel/02_TimerChan/03_TimeAfter/main.go)
- [3 Select 多路复用](chapter03_channel/03_select/select.md)
  - [3.1 配合 default 实现不阻塞发送](chapter03_channel/03_select/01_default_unblock/main.go)
  - [3.2 多 case 随机选择](chapter03_channel/03_select/02_random_select/main.go)
- [4 基于消息传递并发模型：Actor模型和CSP模型-->Golang 在 CSP 模型中应用](chapter03_channel/04_CSP/CSP.md)
---

## 第四章 interface 和反射 
- [1 interface 分类：eface 和 iface, 及两者之间关系转换](chapter04_interface_n_reflect/01_interface/interface.md)
    - [1.1 汇编分析不含方法eface和带方法iface](chapter04_interface_n_reflect/01_interface/01_interface_in_asm/main.go)
    - [1.2 接口值 iface == nil 是指动态类型 iface.tab._type 和动态值 iface.data 都为 nil ](chapter04_interface_n_reflect/01_interface/02_interface_compare_with_nil/main.go)
    - [1.3 模拟打印出接口 eface 的动态类型 itab 和 data 值](chapter04_interface_n_reflect/01_interface/03_print_dynamic_value_n_type/main.go)
- [2 反射](chapter04_interface_n_reflect/02_reflect/reflect.md)
    - [2.0 常见需求: 不能预先确定参数类型，需要动态的执行不同参数类型行为](chapter04_interface_n_reflect/02_reflect/00_kind_route/kind_route_test.go)
    - [2.1 反射三大定律](chapter04_interface_n_reflect/02_reflect/01_three_laws/threeLaw.md)
    - [2.2 四种类型转换:断言、强制、显式、隐式](chapter04_interface_n_reflect/02_reflect/02TypeAssert/type_assertion.md)
        - [2.2.1 断言的类型T是一个**具体类型** 或则 **接口类型**](chapter04_interface_n_reflect/02_reflect/02TypeAssert/01_eface_n_iface_type_assert/main.go)
        - [2.2.2 类型断言性能分析](chapter04_interface_n_reflect/02_reflect/02TypeAssert/02_type_assert_performance/typeAssert_test.go)
            - 空接口类型直接类型断言具体的类型
            - 空接口类型使用TypeSwitch 只有部分类型
            - 空接口类型使用TypeSwitch 所有类型
            - 直接使用类型转换
            - 非空接口类型判断一个类型是否实现了该接口 12个方法
            - 直接调用方法
    - [2.3 动态创建类型](chapter04_interface_n_reflect/02_reflect/03_dynamic_make/main.go)
    - [2.4 通过 reflect 基本函数修改值，调用结构体方法，调用普通函数](chapter04_interface_n_reflect/02_reflect/04_reflect_method/main.go)
    - [2.5 反射性能优化演变案例](chapter04_interface_n_reflect/02_reflect/05_performance_inprove/main.go)
    - [2.6 通过reflect.DeepEqual进行深度比较引用类型](chapter04_interface_n_reflect/02_reflect/06_deepEqual/deepEqual.md)
        - *底层类型相同，相应的值也相同,两个自定义类型*是否“深度”相等
        - *一个 nil 值的map*和*非 nil 值但是空的map*是否“深度”相等
        - *带有环的数据*对比是否“深度”相等
    - [2.7 reflect.implements 判断 struct 是否实现某接口](chapter04_interface_n_reflect/02_reflect/07_implement_interface/main.go)
    - [2.8 go-cmp-->reflect.DeepEqual 的替代品](chapter04_interface_n_reflect/02_reflect/08_go-cmp/go-cmp.md)
      - [2.8.1 结构体内嵌指针：与 == 对比进行相等判断](chapter04_interface_n_reflect/02_reflect/08_go-cmp/01_struct_compare_with_pointer/main.go)
      - [2.8.2 IgnoreUnexported 忽略未导出字段,AllowUnexported 指定某些类型的未导出字段需要比较](chapter04_interface_n_reflect/02_reflect/08_go-cmp/02_ignoreUnexported/main.go)
      - [2.8.3 切片变量值为 nil 与 长度为 0 的切片相等判断，map 实现元素对比 ](chapter04_interface_n_reflect/02_reflect/08_go-cmp/03_nil_and_empty_slice_or_map/main.go)
      - [2.8.4 切片 及 map 相等判断](chapter04_interface_n_reflect/02_reflect/08_go-cmp/04_slice_and_map_equal/main.go)
      - [2.8.5 diff 打印结构体成员区别](chapter04_interface_n_reflect/02_reflect/08_go-cmp/05_diff/main.go)
      - [2.8.6 自定义Equal方法](chapter04_interface_n_reflect/02_reflect/08_go-cmp/06_custom_equal/main.go)
---

## 第五章 切片和数组
- 1 参数传递
  - [1,1 值传递-->数组拷贝，数组作为函数参数传递](chapter05_slice_n_array/01_pass_as_param/01passByValue_array/main.go)
  - [1.2 引用传递-->数组指针，切片和指针切片传递](chapter05_slice_n_array/01_pass_as_param/02passByReference/main.go)
  - [1.3 切片和数组作为参数传递性能对比及注意项](chapter05_slice_n_array/01_pass_as_param/03_array_n_slice_pass_performance/main_test.go)
- [2 切片传递的疑惑](chapter05_slice_n_array/02_slice_pass/slice_n_array_pass.md)
  - [2.1 没有足够容量时函数中切片传递的疑惑](chapter05_slice_n_array/02_slice_pass/01_slice_pass_confusition_without_enough_cap/main.go)
  - [2.2 没有足够容量切片传递疑惑揭秘：底层扩容指向的数据变化](chapter05_slice_n_array/02_slice_pass/02_slice_pass_reality_without_enough_cap/main.go)
  - [2.3 有足够容量时函数中切片传递的疑惑](chapter05_slice_n_array/02_slice_pass/03_slice_pass_confusition_fix_with_enough_cap)
  - [2.4 有足够容量时函数传递疑惑揭秘: 底层len长度没变](chapter05_slice_n_array/02_slice_pass/04_slice_pass_confusition_with_enough_cap)
- [3 带索引初始化数组和切片](chapter05_slice_n_array/03_make_slice_with_index/make_slice_with_index.go)
- 4 底层数据结构
  - [4.1 数组数据结构](chapter05_slice_n_array/04_structure_of_array_n_slice/01_array/arrayStructure.md)
  - [4.2 切片数据结构及拷贝copy源码分析](chapter05_slice_n_array/04_structure_of_array_n_slice/02_slice/sliceStructure.md)
  - [4.3 array 转 slice,slice 转 array 在版本 1.20 前后变化](chapter05_slice_n_array/04_structure_of_array_n_slice/03_slice_to_array/main.go)
- [5 nil 切片和 empty 切片](chapter05_slice_n_array/05nilSlice_n_NoneSlice/nil_n_empty_slice.md)
  - [5.1 优雅的清空切片，复用内存](chapter05_slice_n_array/05nilSlice_n_NoneSlice/clear_slice.go)
- [6 扩容策略](chapter05_slice_n_array/06GrowSlice/grow_size_policy.md)
- [7 不同类型的切片间互转](chapter05_slice_n_array/07_transfer_slice_in_different_type/main.go)
- [8 切片复制方式对比: copy和=复制](chapter05_slice_n_array/08_reslice_n_copy/slice_copy.md)
- [9 append 切片常用考题](chapter05_slice_n_array/09_append/main.go)
- [10 并发访问 slice 如何做到优雅和安全](chapter05_slice_n_array/10_concurrency_slice/slice_concurrency.md)
- [11 go1.21 切片泛型库](chapter05_slice_n_array/11_slice_in_1_21/slice_in_1_21.md)
---

## 第六章 指针
- [1 指针类型转换及修改值](chapter06_pointer/01_ptrOperation/ptr_operation.md)
- [2 指针分类及unsafe包使用](chapter06_pointer/02unsafe/unsafe.md)
    - [2.1 sizeof 获取类型其占用的字节数, Offsetof 修改结构体私有成员, Alignof 内存中的地址对齐](chapter06_pointer/02unsafe/01_basic_api/unsafe.go)
    - [2.2 指针获取切片长度和容量](chapter06_pointer/02unsafe/02_slice_operaion/slice_len_n_cap.go)
    - [2.3 指针获取Map的元素数量](chapter06_pointer/02unsafe/03_map_count/main.go)
    - [2.4 使用指针来访问数组里的所有元素](chapter06_pointer/02unsafe/04_array_filed/array_field.go)
- [3 获取并修改结构体私有变量值](chapter06_pointer/03PointerSetPrivateValue/main.go)
- [4 []byte 切片 与 string 字符串实现零拷贝互转(指针和反射方式)](chapter06_pointer/04SliceToString/sliceToString.go)
- [5 结构体的内存对齐规则](chapter06_pointer/05_struct_align/struct_align.md)
    - [5.1 结构体排序优化内存占用](chapter06_pointer/05_struct_align/01_struct_mem/align.go)
    - [5.2 空 struct{} 结构体使用](chapter06_pointer/05_struct_align/02_empty_struct/empty_struct.go)
        - 空结构体作为第一个元素
        - 空结构体作为最后一个元素
- [6 弱指针--go 1.24 ](chapter06_pointer/06_weak_pointer/week_pointer.md)
  - [6.1 基本使用](chapter06_pointer/06_weak_pointer/01_basic_use/main.go)
  - [6.2 缓存使用](chapter06_pointer/06_weak_pointer/02_cache/main.go))
- [7 keepalive 配合 finalizer 使用](chapter06_pointer/07_keepalive_n_finalizer/keepalive.md)
  - [7.1 错误使用 finalizer 优化](chapter06_pointer/07_keepalive_n_finalizer/01_wrong_use/main.go)
  - [7.2 正确使用 finalizer 优化](chapter06_pointer/07_keepalive_n_finalizer/02_correct_use/main.go)
---

## [第七章 系统调用](chapter07_system_call/Syscall.md)
- [1 基于 kqueue event loop 的 TCP server（涉及各种linux系统调用](chapter07_system_call/01_kqueue_server/kqueue.md)
- [2 使用 strace 工具追踪进程与内核的交互情况:如系统调用](chapter07_system_call/02_ptrace/ptrace.md)
  - [2.1 syscall.PtraceGetRegs 获取所有寄存器的值](chapter07_system_call/02_ptrace/01_register/main.go)
  - [2.2 使用 libseccomp-golang 查看 echo hello 的系统调用及次数](chapter07_system_call/02_ptrace/02_follow_system_call/main.go)
- [3 exec 执行命令](chapter07_system_call/03_exec/exec.md)
---

## [第八章 defer函数及汇编语言理解](chapter08_defer/defer.md)
- [1 注册延迟调用机制定义及使用](chapter08_defer/01_defer_definiton/main.go)
- [2 defer 陷阱](chapter08_defer/02_defer_common_mistakes/main.go)
- [3 分阶段解析 defer 函数](chapter08_defer/03_defer_params_n_return/main.go)
- [4 defer 循环性能问题](chapter08_defer/04_defer_loop_performance/main.go)
- [5 汇编理解 defer 函数(AMD)](chapter08_defer/05_defer_assembly/defer_amd.s)
---

## [第九章 设计模式-OOP七大准则](chapter09_design_pattern/introduction.md)
- 1 创建型模式 Creational Patterns
    - 1.1 工厂模式(Factory Design Pattern)
      - [1.1.1 简单工厂模式-->new关键字函数实现简单工厂](chapter09_design_pattern/01_construction/01_factory/01_StaticFactory/static_factory.md)
      - [1.1.2 工厂方法模式-->k8s中实现](chapter09_design_pattern/01_construction/01_factory/02_factory_mode/factory.md)
      - [1.1.3 抽象工厂模式](chapter09_design_pattern/01_construction/01_factory/03_abstract_factory/abstract_factory.md)
    - [1.2 单例模式(Singleton Design Pattern)-->标准库strings/replace实现](chapter09_design_pattern/01_construction/03_singleton/singleton.md)
    - [1.3 原型模式(Prototype Design Pattern)-->confluent-kafka中map实现](chapter09_design_pattern/01_construction/04_prototype/prototype.md)
    - [1.4 建造者模式(Builder Design Pattern)-->xorm,k8s,zap中实现](chapter09_design_pattern/01_construction/05_builder/builder_info.md)
- 2 结构型模式 Structural Patterns
    - 2.1 组合模式(Composite Design Pattern)
        - [2.1.1 修改前：使用面向对象处理](chapter09_design_pattern/02_structure/01_Composite/01_modify_before/composite.go)
        - [2.1.2 修改后：使用组合模式处理](chapter09_design_pattern/02_structure/01_Composite/02_modify_after/conposite.go)
    - [2.2 装饰模式(Decorator Design Pattern)-->grpc源码体现](chapter09_design_pattern/02_structure/02_Decorate/decorate.md)
        - [2.2.1 闭包实现--多个装饰器同时使用](chapter09_design_pattern/02_structure/02_Decorate/01_closure_decorate/main.go)
        - [2.2.2 结构体装饰](chapter09_design_pattern/02_structure/02_Decorate/02_struct_decorate_inGrpc/main.go)
        - [2.2.3 反射实现--泛型装饰器](chapter09_design_pattern/02_structure/02_Decorate/03_reflect_decorate/decorate.go)
    - [2.3 享元模式(Flyweight Design Pattern)-->线程池,缓存思想](chapter09_design_pattern/02_structure/03_FlyweightPattern/flyWeightPattern.md)
    - [2.4 适配器模式(Adapter Design Pattern)](chapter09_design_pattern/02_structure/04_adopter/adopter.md)
    - [2.5 桥接模式(Bridge Design Pattern)](chapter09_design_pattern/02_structure/05_bridgeMethod/bridge_method.md)
    - [2.6 门面模式(外观模式Facade Design Pattern)-->在gin中render应用(封装多个子服务)](chapter09_design_pattern/02_structure/06_facade_pattern/facade.md)
    - [2.7 代理模式(Proxy Design Pattern)](chapter09_design_pattern/02_structure/07_proxy/proxy.md)
- 3 行为型模式 Behavioral Patterns
    - [3.1  访问者模式(Visitor Design Pattern)-->k8s中kubectl实现](chapter09_design_pattern/03_motion/01_visitor/vistor.md)
    - [3.2  迭代器(Iterator Design Pattern)-->标准库container/ring中实现](chapter09_design_pattern/03_motion/02_Iterator/main.go)
    - [3.3  状态模式(State Design Pattern)](chapter09_design_pattern/03_motion/03_State/introduction.md)
    - [3.4  责任链模式(Chain Of Responsibility Design Pattern)-->gin 中间件中使用](chapter09_design_pattern/03_motion/04_duty_chain_method/duty_chain.md)
    - [3.5  模版模式(Template Method Design Pattern)](chapter09_design_pattern/03_motion/05_templateMethod/templateMethod.md)
    - [3.6  策略模式(Strategy Method Design Pattern)-->if-else的另类写法(内部算法封装)](chapter09_design_pattern/03_motion/06_strategyMethod/strategy.md)
    - [3.7  解释器模式(Interpreter Design Pattern)](chapter09_design_pattern/03_motion/07_InterpreterMethod/interpreter.md)
    - [3.8  命令模式(Command Design Pattern)-->go-redis中实现](chapter09_design_pattern/03_motion/08_CommandMethod/command.md)
    - [3.9  备忘录模式(Memento Design Pattern)](chapter09_design_pattern/03_motion/09_memento/introduction.md)
    - [3.10 观察者模式(Observer Design Pattern)-->官方Signal包及etcd的watch机制](chapter09_design_pattern/03_motion/10_ObserverPattern/introduction.md)
    - [3.11 中介者模式(Mediator Design Pattern)](chapter09_design_pattern/03_motion/11_mediator/introduction.md)
- [4 函数选项:成例模式-->在日志库zap中实现](chapter09_design_pattern/04_fuctional_option/option.md)
    - [4.1 未使用函数选项初始化结构体的现状](chapter09_design_pattern/04_fuctional_option/01_problem/ServerConfig.md)
    - [4.2 区分必填项和选项](chapter09_design_pattern/04_fuctional_option/02_method_splitConfig/SplitConfig.go)
    - 4.3 带参数的选项模式
      - [不返回error](chapter09_design_pattern/04_fuctional_option/03_FunctionalOption/01_simple_solution/main.go)
      - [返回error](chapter09_design_pattern/04_fuctional_option/03_FunctionalOption/02_complexed_with_error/main.go)
- [5 插件式编程-->grpc中实现](chapter09_design_pattern/05_plugin_programming/plugin.md)
- [6 同步模式(sync同步原语以及扩展原语)](chapter09_design_pattern/06_Synchronization_mode/01_sync/sync.md)
    - [6.1 waitGroup同步等待组对象](chapter09_design_pattern/06_Synchronization_mode/01_sync/01waitGroup/waitGroup.md)
    - 6.2 使用互斥锁（sync.Mutex）实现读写功能和直接使用读写锁（sync.RWMutex）性能对比
      - [6.2.1 使用互斥锁（sync.Mutex）实现读写功能](chapter09_design_pattern/06_Synchronization_mode/01_sync/02RWMutex_vs_mutex/Mutex/main.go)
      - [6.2.2 直接使用读写锁（sync.RWMutex）实现读写功能](chapter09_design_pattern/06_Synchronization_mode/01_sync/02RWMutex_vs_mutex/RWMutex/main.go)
      - [Mutex 和 RWMutex 源码分析](chapter09_design_pattern/06_Synchronization_mode/01_sync/02RWMutex_vs_mutex/mutex.md)
    - [6.3 Once源码分析](chapter09_design_pattern/06_Synchronization_mode/01_sync/03Once/once.md)
    - [6.4 并发安全的sync.Map与sync.RWMutex封装的map对比及源码分析](chapter09_design_pattern/06_Synchronization_mode/01_sync/04map/sync_map.md)
    - [6.5 Pool对象池模式( *非连接池*)-->官方包对象池fmt](chapter09_design_pattern/06_Synchronization_mode/01_sync/05Pool/pool.md)
        - [6.5.1 错误使用：未使用newFunc](chapter09_design_pattern/06_Synchronization_mode/01_sync/05Pool/01Without_newFunc/main.go)
        - [6.5.2 newFunc与GC前后Get对比](chapter09_design_pattern/06_Synchronization_mode/01_sync/05Pool/02NewFunc/newFunc.go)
        - [6.5.3 何时使用对象缓存](chapter09_design_pattern/06_Synchronization_mode/01_sync/05Pool/03When2Use_object_pool/main.go)
        - [6.5.4 第三方对象池object pool(bytebufferpool)](chapter09_design_pattern/06_Synchronization_mode/01_sync/05Pool/04_byteBufferPool/main.go)
    - [6.6 Cond 条件变量协调想要访问共享资源的goroutine-->熔断框架 hystrix-go 优秀实现](chapter09_design_pattern/06_Synchronization_mode/01_sync/06Cond/Cond.md)
    - [6.7 atomic原子操作源码分析-->zerolog源码中实现](chapter09_design_pattern/06_Synchronization_mode/01_sync/07Atomic/atomic.md)
        - [6.7.0 Value的load和store](chapter09_design_pattern/06_Synchronization_mode/01_sync/07Atomic/00_value/main.go)
        - [6.7.1 add及补码减](chapter09_design_pattern/06_Synchronization_mode/01_sync/07Atomic/01_add/main.go)
        - [6.7.2 cas算法和自旋锁](chapter09_design_pattern/06_Synchronization_mode/01_sync/07Atomic/02_CompareAndSwap/main.go)
        - [6.7.3 load和store用法](chapter09_design_pattern/06_Synchronization_mode/01_sync/07Atomic/03_load_n_store/main.go)
        - [6.7.4 swap交换](chapter09_design_pattern/06_Synchronization_mode/01_sync/07Atomic/04_swap/main.go)
    - [6.8 ErrorGroup获取协程中error](chapter09_design_pattern/06_Synchronization_mode/01_sync/08ErrorGroup/errGroup.md)
        - [6.8.1 不带context](chapter09_design_pattern/06_Synchronization_mode/01_sync/08ErrorGroup/01WithoutContext/main.go)
        - [6.8.2 带context](chapter09_design_pattern/06_Synchronization_mode/01_sync/08ErrorGroup/02WithContext/main.go)
    - [6.9 信号量 Semaphore-->mutex是二进制信号量 binary semaphore](chapter09_design_pattern/06_Synchronization_mode/01_sync/09Semaphore/semaphore.md)
    - [6.10 SingleFlight避免缓存击穿](chapter09_design_pattern/06_Synchronization_mode/01_sync/10SingleFlight/singleFlight.md)
        - [6.10.1 Do方法](chapter09_design_pattern/06_Synchronization_mode/01_sync/10SingleFlight/01_do/main.go)
        - [6.10.2 DoChan方法](chapter09_design_pattern/06_Synchronization_mode/01_sync/10SingleFlight/02_do_chan/main.go)
    - [6.11 NoCopy 机制](chapter09_design_pattern/06_Synchronization_mode/01_sync/11_nocopy/nocopy.md)
---

## [第十章 函数式编程](chapter10_function/func.md)
- 1 函数应用
    - [1.1 闭包基本使用](chapter10_function/01_func_application/01_closure/main.go)
    - [1.2 匿名函数应用:回调函数](chapter10_function/01_func_application/02_anonymousFunc/main.go)
    - [1.3 函数模版:定义行为](chapter10_function/01_func_application/03_func_template/main.go)
- [2 高级函数filter,map,reduce](chapter10_function/02_advanced_function/introduction.md)
    - [2.1 简单案例理解 filter,map,reduce](chapter10_function/02_advanced_function/01_simple_solution/main.go)
    - 2.2 interface{} + reflect 实现泛型->Go1.18之前
        - [filter](chapter10_function/02_advanced_function/01_simple_solution/main.go)
        - [map](chapter10_function/02_advanced_function/02_generic_n_parameter_check/map/main.go)
        - [reduce](chapter10_function/02_advanced_function/02_generic_n_parameter_check/reduce/main.go)
    - [2.3 go-zero框架实现 map-reduce](chapter10_function/02_advanced_function/03_go_zero_map_reduce/go_zero_map_reduce.md)
        - [2.3.1 Finish函数：一个应用依托于很多服务,在没有强依赖关系下,优雅地实现并发编排任务](chapter10_function/02_advanced_function/03_go_zero_map_reduce/01_producer/main.go)
    - [2.4 RXGo基于pipelines实现ReactiveX 编程模型](chapter10_function/02_advanced_function/04_rxgo/rxgo.md)
        - [2.4.1 map,reduce 使用](chapter10_function/02_advanced_function/04_rxgo/main.go)
- 3 一等公民案例
    - [网络管理中需求](chapter10_function/03_Firstclassfunction/problem_desc.md)
    - 网络管理中三种处理对比
        - [3.1 通过同享内存通信](chapter10_function/03_Firstclassfunction/01_communicate_by_sharing_memory/main.go)
        - [3.2 通过通信(具体数据)共享内存](chapter10_function/03_Firstclassfunction/02_sharing_memory_by_communicating/main.go)
        - [3.3 通过通信(函数)共享内存](chapter10_function/03_Firstclassfunction/03_send_func_to_channel/main.go)
---

## 第十一章 汇编理解go语言底层源码(AMD芯片运行代码)
- [1 汇编基本指令](chapter11_assembly_language/01asm/introduction.md)
- [2 plan9 手写汇编](chapter11_assembly_language/02plan9/introduction.md)
    - [2.1  变量var，常量constant](chapter11_assembly_language/02plan9/01_pkg_constant_string/main.go)
    - [2.2  array数组](chapter11_assembly_language/02plan9/02_pkg_array/main.go)
    - [2.3  bool类型](chapter11_assembly_language/02plan9/03_pkg_bool/main.go)
    - [2.4  int,int32,uint32类型](chapter11_assembly_language/02plan9/04_pkg_int/main.go)
    - [2.5  float32，float64类型](chapter11_assembly_language/02plan9/05_pkg_float/main.go)
    - [2.6  slice切片([]byte)](chapter11_assembly_language/02plan9/06_pkg_slice/main.go)
    - [2.7  引用类型map和channel](chapter11_assembly_language/02plan9/07_pkg_channel_n_map/main.go)
    - [2.8  Go 函数申明](chapter11_assembly_language/02plan9/08_pkg_func/func.md)
    - [2.9  局部变量](chapter11_assembly_language/02plan9/09_local_param/local_params.md)
    - [2.10 流程控制](chapter11_assembly_language/02plan9/10_control_process/main.go)
    - [2.11 伪寄存器 SP 、伪寄存器 FP 和硬件寄存器 SP关系](chapter11_assembly_language/02plan9/11_FalseSP_fp_SoftwareSP_relation/main.go)
    - [2.12 结构体方法](chapter11_assembly_language/02plan9/12_struct_method/main.go)
    - [2.13 递归函数](chapter11_assembly_language/02plan9/13_recursive_func/main.go)
    - [2.14 闭包函数](chapter11_assembly_language/02plan9/14_closure/main.go)
    - [2.15 两种方式获取GoroutineId](chapter11_assembly_language/02plan9/15_GoroutineId/main.go)
    - [2.16 汇编调用非汇编Go函数](chapter11_assembly_language/02plan9/16_assembly_call_NonassemblyFunc/main.go)

---
## [第十二章 net 网络--涉及性能指标,协议栈统计,套接字信息](chapter12_net/net.md)
- [socket 套接字缓冲区](chapter12_net/socket.md)
- [tcp 传输控制协议](chapter12_net/tcp.md)
- [I/O 多路复用及 epoll 在 Golang 工作模型体现](chapter12_net/io_multiplexing.md)
- [http的三个版本知识介绍](chapter12_net/http.md)
---

- 1 http 服务端高级封装演变: ServeHTTP 是 HTTP 服务器响应客户端的请求接口
  - [1.1 高级封装：使用 DefaultServeMux](chapter12_net/01_http_server/01_use_DefaultServeMux/main.go)
  - [1.2 中级封装：使用内置 serveMux 生成函数](chapter12_net/01_http_server/02_use_http_NewServeMux/main.go)
  - [1.3 原始封装：自定义实现 serveMux](chapter12_net/01_http_server/03_use_cutomized_mux/main.go)
- 2 http 客户端高级封装演变
  - [request 源码](chapter12_net/02_http_client/http_request.md)
  - [response 源码](chapter12_net/02_http_client/http_response.md)
  - [http.RoundTripper 接口实现源码: 调用方将请求作为参数获取请求对应的响应,并管理连接](chapter12_net/02_http_client/http_transport.md)
  - [http.Client 源码](chapter12_net/02_http_client/http_client.md) 
  - [2.1 官方库版(爬虫获取邮箱案例-未封装)](chapter12_net/02_http_client/01_standard_pkg/client.go)
  - [2.2 go-resty(推荐使用)](chapter12_net/02_http_client/02_go_resty/rest_client.go)
- [3 Tcp 实现 Socket 编程 (服务端 netpoll 分析)](chapter12_net/03_tcp/tcp_server.md)
  - [3.1 客户端](chapter12_net/03_tcp/client/main.go)
  - [3.2 服务端](chapter12_net/03_tcp/server/main.go)
- [4 Tcp 黏包分析及处理(大小端介绍)](chapter12_net/04_tcp_sticky_problem/big_n_small_endian.md)
  - [4.1 TCP 粘包问题](chapter12_net/04_tcp_sticky_problem/01_problem)
  - [4.2 TCP 粘包解决方式](chapter12_net/04_tcp_sticky_problem/02_solution)
  - [4.3 可变长度编码 Variable-length encoding](chapter12_net/04_tcp_sticky_problem/03_varint/main.go)
- [5 fastHttp(源码分析)](chapter12_net/05_fasthttp/fasthttp.md)
  - [5.1 服务端](chapter12_net/05_fasthttp/server/main.go)
  - [5.2 客户端](chapter12_net/05_fasthttp/client/client.go)
- [6 优雅退出原理分析-涉及linux信号介绍（go-zero实践）](chapter12_net/06_grateful_stop/grateful_stop.md)
  - [6.1 信号监听处理](chapter12_net/06_grateful_stop/signal.go)
- [7 URL的解析 Parse，query 数据的转义与反转义](chapter12_net/07_url/url.md)
- [8 使用 alice 优雅编排中间件](chapter12_net/08_middleware/middleware.md)
  - [5.1 jwt 中间件载体 Symmetric 对称加密->HSA](chapter12_net/08_middleware/01_symmetric/jwt_test.go)
  - [5.2 jwt 中间件载体 asymmetric 非对称加密(更安全)->RSA](chapter12_net/08_middleware/02_asymmetric/jwt_test.go)
- [9 HTTPS, SAN, SLS, TLS及源码分析握手过程](chapter12_net/09_https/https.md)
  - 9.1 https 单向认证
    - [9.1.1 服务端修改 tls 版本](chapter12_net/09_https/01_sign_one/01_server/server.go)
    - [9.1.2 客户端不校验证书 或则 添加到证书池](chapter12_net/09_https/01_sign_one/02_client/client.go)
  - 9.2 https 双向认证
- [10 unix domain socket 本地 IPC 进程间通信](chapter12_net/10_unix_domain_socket/uds.md)
- [11 获取本机内网和外网Ip](chapter12_net/11_internal_n_external_ip/main.go)
- [12 http2 使用](chapter12_net/12_http2/http2.md)
  - [let's encrypt 免费证书 开发server](chapter12_net/12_http2/01_server/server.go)
- [13 flusher 实现 stream 流式返回](chapter12_net/13_flusher/flusher.md)
- [14 ReverseProxy 反向代理](chapter12_net/14_reverse_proxy/reverse_proxy.md)
- [15 go1.18 netip 处理网络地址和相关操作](chapter12_net/15_netip/netip.md)
- [16 dns 解析](chapter12_net/16_dns/dns.md)

---
## [第十三章 CGO调用C语言](chapter13_Go_call_C_or_C++/introduction.md)
[cgo在confluent-kafka-go源码使用](https://github.com/Danny5487401/go_grpc_example/blob/master/03_amqp/02_kafka/02_confluent-kafka/confluent_kafka_source_code.md)

**Note: 内部c代码需要自己编译成对应本地 静态库 或则 动态库,[可参考C基本知识](https://github.com/Danny5487401/c_learning)**

- [1 Go调用自定义C函数-未模块化](chapter13_Go_call_C_or_C++/01_call_C_func/main.go)
- [2 Go调用自定义C函数-模块化](chapter13_Go_call_C_or_C++/02_call_C_module/main.go)
- 3 Go重写C定义函数
- [4 cgo错误用法：引入其他包的变量](chapter13_Go_call_C_or_C++/04_import_other_pkg/main.go)
- [5 #Cgo语句](chapter13_Go_call_C_or_C++/05_cgo/main.go)
- [6 Go获取C函数的errno](chapter13_Go_call_C_or_C++/06_return_err/main.go)
- [7 C的void返回](chapter13_Go_call_C_or_C++/07_void_return/main.go)
- 8 Go调用模块化C++库
- [9 调用静态C库](chapter13_Go_call_C_or_C++/09_static_c_lib/main.go)
- [10 调用动态C库](chapter13_Go_call_C_or_C++/10_dynamic_c_lib/main.go)
- [11 Golang使用pkg-config自动获取头文件和链接库的方法](chapter13_Go_call_C_or_C++/11_pkg_config/pkg_config.md)
---
## [第十四章 Context上下文-源码分析涉及父类EmptyCtx](chapter14_context/introduction.md)
- 1 Context使用背景
    - [1.1 问题：如何通过父进程优雅释放子goroutine资源](chapter14_context/01_Reason_To_Use_Context/01_problem/main.go)
    - [1.2 方式一：全局参数方式解决的优缺点](chapter14_context/01_Reason_To_Use_Context/02_Method1_Global_Param/main.go)
    - [1.3 方式二: 通道channel方式解决的优缺点](chapter14_context/01_Reason_To_Use_Context/03_Method2_Channel/main.go)
    - [1.4 方式三: 最优方式Context](chapter14_context/01_Reason_To_Use_Context/04_Method3_Context/main.go)
- [2 WithCancel 使用](chapter14_context/02_WithCancel/main.go)
- [3 WithDeadline 使用](chapter14_context/03_WithDeadline/main.go)
- [4 WithValue 使用](chapter14_context/04_WIthValue/main.go)
- [5 WithTimeout 对 WithDeadline 封装的使用](chapter14_context/05_WithTimeout/main.go)
- [6 Go1.21 增加取消原因以及回调函数的增添](chapter14_context/06_cancel_reason_in_go1_21/main.go)
---
## 第十五章 接口嵌套编程
- [1 常见冗余代码写法](chapter15_interfaceProgramming/01_problem/main.go)
- [2 简单优化](chapter15_interfaceProgramming/02_simple_method/main.go)
- [3 更优方式](chapter15_interfaceProgramming/03_better_solution/main.go)
---
## 第十六章 并发编程
- [1 简单流水线模型](chapter16_concurrentProgramming/01_pipeline/pipeline.md)
- [2 FAN-IN和FAN-OUT模型](chapter16_concurrentProgramming/02_fanin_fanout/fanin_fanout.md)
---
## 第十七章 数据结构及算法
- [1 queue 双端单向队列](chapter17_dataStructure_n_algorithm/01_queue/queue.md)
  - [1.1 array 数组实现非阻塞 Queue-->没有实现容量限制和容量收缩](chapter17_dataStructure_n_algorithm/01_queue/01_unblock_queue/queue_test.go)
  - [1.2 array + channel 实现阻塞队列 Queue](chapter17_dataStructure_n_algorithm/01_queue/02_block_queue/blockqueue_test.go)
- [2 加解密](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/encryption.md)
  - 2.1 对称式加密
    - [2.1.1 AES高级加密标准(Advanced Encryption Standard)](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/01_symmetric_encryption/aes/aes.md)
      - [Cipher FeedBack密码反馈模式](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/01_symmetric_encryption/aes/01_cfb/aes_cfb.go)
      - [Cipher Block Chaining密码分组链接模式](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/01_symmetric_encryption/aes/02_cbc/aes_cbc.go)
    - [2.1.2 des美国数据加密标准(不推荐)](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/01_symmetric_encryption/des/des.md)
  - [2.2 哈希算法及其在数字签名中应用(hmac,md5,sha1,sha256)](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/02_digital_signature/hash.md)
    - [MD 5信息摘要算法(Message-DigestAlgorithm 5)](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/02_digital_signature/md5/md5.md)
    - [HMAC 哈希运算消息认证码(Hash-based Message Authentication Code)](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/02_digital_signature/hmac/hmac.md)
    - [SHA安全散列算法(secure Hash Algorithm)](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/02_digital_signature/sha1_n_sha256/sha.md)
    - [号称计算速度最快的哈希 xxhash](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/02_digital_signature/xxhash/main.go)
    - [consistant hash 一致性hash](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/02_digital_signature/consistent_hash/consistent_hash.md)
  - [2.3 非对称加密算法 rsa](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/03_rsa/rsa.md)
    - [分段与不分段进行加解密](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/03_rsa/main.go)
- [3 随机算法（伪随机和真随机)](chapter17_dataStructure_n_algorithm/03_rand/rand.md)
  - [3.1 math_rand使用](chapter17_dataStructure_n_algorithm/03_rand/01_math_rand/main.go)
  - [3.2 crypto_rand使用](chapter17_dataStructure_n_algorithm/03_rand/02_crypto_rand/main.go)
  - [3.3 fastrand 优化使用](chapter17_dataStructure_n_algorithm/03_rand/03_fastrand/main.go)
- [4 排序算法分类及图解(sort包源码分析)](chapter17_dataStructure_n_algorithm/04_sort/sort.md)
  - [4.1 sort 排序](chapter17_dataStructure_n_algorithm/04_sort/01_sorted_info/main.go)
    - 不同结构体切片根据反射reflect实现自定义排序: sort.Sort 接口实现
    - 相同结构体切片排序: sort.Slice
    - map根据key实现排序: sort.Strings
    - int 类型切片排序: sort.IntSlice
  - [4.2 sort.Search 二分查找: 根据排序切片找索引](chapter17_dataStructure_n_algorithm/04_sort/02_search/main.go)
- [5 container](chapter17_dataStructure_n_algorithm/05_container/container.md)
  - [5.1 heap 堆](chapter17_dataStructure_n_algorithm/05_container/01_heap/heap.md)
    - [5.1.1 最小堆](chapter17_dataStructure_n_algorithm/05_container/01_heap/01_basic/main.go)
    - [5.1.2 优先队列](chapter17_dataStructure_n_algorithm/05_container/01_heap/02_priority)
  - [5.2 list 双向链表](chapter17_dataStructure_n_algorithm/05_container/02_list/list.md) 
  - [5.3 ring 环形链表](chapter17_dataStructure_n_algorithm/05_container/03_ring/ring.md) 
- [6 certificate 证书-->openssl 使用](chapter17_dataStructure_n_algorithm/06_certificate/certificate.md)
  - [6.1 pem(Privacy Enhanced Mail Certificate保密增强邮件协议](chapter17_dataStructure_n_algorithm/06_certificate/01_pem/pem.md)
    - [6.1.1 生成公私钥的 .pem 文件(公钥使用RSA算法)](chapter17_dataStructure_n_algorithm/06_certificate/01_pem/01_pem_generate/main.go)
    - [6.1.2 解析.pem文件获取公私钥](chapter17_dataStructure_n_algorithm/06_certificate/01_pem/02_get_pem_info/main.go) 
  - [6.2 x509 库源码](chapter17_dataStructure_n_algorithm/06_certificate/02_x509/x509.md)
    - [6.2.1 ca 创建根证书并签署终端证书](chapter17_dataStructure_n_algorithm/06_certificate/02_x509/main.go)
- [7 Base64编码解析](chapter17_dataStructure_n_algorithm/07_base64_encoding/base64.md)
- [8 trie前缀树](chapter17_dataStructure_n_algorithm/08_trie/trie.md)
- [9 Golang底层数据结构-涉及数值类型占用的bit](chapter17_dataStructure_n_algorithm/09_golang_data_structure/data.md)
    - [9.1 Map底层结构](chapter17_dataStructure_n_algorithm/09_golang_data_structure/01_map_structure/map_intro.md)
        - [9,1,1 根据预期初始大小和 loadfactor负载因子 6.5 计算一个数组的长度的对数B](chapter17_dataStructure_n_algorithm/09_golang_data_structure/01_map_structure/01_buckets_B/main.go)
        - [9,1,2 map的指针优化场景](chapter17_dataStructure_n_algorithm/09_golang_data_structure/01_map_structure/02_Improvement/map_test.go)
        - [9.1.3 map 的 Key 类型取值](chapter17_dataStructure_n_algorithm/09_golang_data_structure/01_map_structure/03_map_key/key.md)
        - [9.1.4 简单实现 hashMap](chapter17_dataStructure_n_algorithm/09_golang_data_structure/01_map_structure/04_self_HashMap/self_hash_map.md)
    - [9.2 String 底层结构,字符集和字符编码,性能分析及内存泄漏分析](chapter17_dataStructure_n_algorithm/09_golang_data_structure/02_string_structure/str.md)
    - [9.3 Struct 底层结构,内存布局,空结构体内存对齐](chapter17_dataStructure_n_algorithm/09_golang_data_structure/03_struct_structure/struct.md)
- [10 copy_on_write 写入时复制-->golang官方库btree](chapter17_dataStructure_n_algorithm/10_copy_on_write/copy_on_write.md)
  - [10.1 使用 RWMutex 缺点](chapter17_dataStructure_n_algorithm/10_copy_on_write/01_rwmutex/main.go)
  - [10.2 使用 copy_on_write 优化](chapter17_dataStructure_n_algorithm/10_copy_on_write/02_copy_on_write/main.go)
## 第十八章 错误跟踪 error 和 panic
- [0 错误(err)和异常（exception）区别及处理方式](chapter18_error_n_panic/00_diff_between_err_n_exception/main.go)
- [1 自定义错误类型打印错误栈](chapter18_error_n_panic/01_customized_error/customized_error.go)
- [2 扩展包pkg.errors](chapter18_error_n_panic/02_pkg_errors/pkg_errors.md)
- [3 Gin的错误recover分析(panic和recover源码分析)](chapter18_error_n_panic/03_recover/panic.md)
- 4 生成errCode错误码及信息
  - [4.1 传统方式：命名错误码、状态码的同时，又要map同步写码对应的翻译](chapter18_error_n_panic/04_errorCode/01_traditional/main.go)
  - [4.2 stringer + go generate 自带工具生成errCode错误码及信息->效率高于map映射错误](chapter18_error_n_panic/04_errorCode/02_generate_n_stringer/intro.md)
- [5 error如何正确比较](chapter18_error_n_panic/05_err_comparision/main.go)
- [6 收集多个errors-->go-multierror实现](chapter18_error_n_panic/06_multi_error/01_one_goroutine_n_errors/main.go)
- [7 错误链](chapter18_error_n_panic/07_chain_error/chain_err.md)
  - [7.1 errors.Unwrap 获取错误链中最后面的一个 root error](chapter18_error_n_panic/07_chain_error/01_root_error/main.go)
  - [7.2 errors.As函数 提取 error chain中特定类型的error](chapter18_error_n_panic/07_chain_error/02_error_as/main.go)  
- [8 debug.SetCrashOutput -->go 1.23 允许设置未被捕获的错误、异常的日志写入](chapter18_error_n_panic/08_SetCrashOutput/main.go)
## [第十九章 nil 预定义标识](chapter19_nil/nil.md)
- [1 pointer, channel, func, interface, map, or slice type 为nil时的地址和size大小](chapter19_nil/01_nil_size_n_addr/main.go)
- 2 不同类型与nil的比较
  - [2.1 nil==nil不可以比较](chapter19_nil/02_comparison/01_nil/main.go)
  - [2.2 两个 nil 值未必相等](chapter19_nil/02_comparison/02_nil_to_nil/main.go)
  - [2.3 interface为nil时:数据段和类型](chapter19_nil/02_comparison/03_interface/interface.go)
  - [2.4 ptr,channel,func,map为nil必须地址未分配](chapter19_nil/02_comparison/04_ptr_chan_func_map/main.go)
  - [2.5 slice 的长度和容量不决定nil](chapter19_nil/02_comparison/05_slice/slice.go)
- 3 不同类型 nil 时的特点
  - [3.1 channel为 nil 时的接收，发送，关闭及select](chapter19_nil/03_Attribute/01_channel/chan.go)
  - [3.2 map 为 nil 时可读不可写](chapter19_nil/03_Attribute/02_map/map.go)
  - [3.3 结构体指针为 nil 时是否可以调用方法](chapter19_nil/03_Attribute/03_struct_method/ptr.go)
---
## [第二十章 for-range 源码分析](chapter20_for_range/for_range.md)
- [1 遍历数组,切片,结构体数组](chapter20_for_range/01_for_range_slice_n_array/main.go)
- [2 正确遍历 Goroutine 捕获变量 (解析协程启动时间) 及在 GO 1.21 使用 EXPERIMENT=loopvar重新定义循环](chapter20_for_range/02_for_range_goroutine/main.go)
- [3 遍历 Map(增加或删除map元素时)](chapter20_for_range/03_for_range_map/main.go)
---
## [第二十一章 time标准包源码分析](chapter21_time/time.md)
- [1 比time.Now()更优雅获取时间戳（go:link技术）](chapter21_time/01_time_sec.go)
- [2 time.Format()优化写法](chapter21_time/02_append_format.go)
---
## [第二十二章 数据驱动模板源码分析-->kratos工具](chapter22_template/template.md)
- [1 加载多个模版并指定模版生成结构体方法](chapter22_template/01_multi_template/main.go)
- [2 自定义扩展模版函数 FuncMap ](chapter22_template/02_template_func/main.go)
- [3 html模版](chapter22_template/03_html_template/main.go)
- [4 generate 根据模版代码生成](chapter22_template/04_gen_template/gen_main.go)
- [5 推荐第三方 sprig 模版函数](chapter22_template/05_sprig_func/sprig.md)
---
## 第二十三章 调试内部对象
- [1 fmt 打印结构体中含有指针对象, 数组或者map中是指针对象, 循环结构时的困难](chapter23_debug_program/01_fmt_problem/main.go)
- [2 go-spew 优化调试](chapter23_debug_program/02_go_spew/go_spew.md)

---
## [第二十四章 命令行参数解析](chapter24_flag/flag.md)
- [1 flag 基本使用及自定义帮助信息](chapter24_flag/01_flag/nginx.go) 
- [2 pflag 完全兼容flag](chapter24_flag/02_pflag/pflag.md) 

---
## [第二十四章 Flag命令行参数及源码分析](chapter24_flag/flag.md)
- [1 标准包flag基本使用及自定义帮助信息](chapter24_flag/01_flag/nginx.go)
- [2 第三方包pflag：兼容标准包flag](chapter24_flag/02_pflag/pflag.md)

## 第二十五章 结构体类型方法
- [1 方法调用语法糖](chapter25_struct_method/01_struct_method/main.go)
---

## [第二十六章 strconv 字符串和数值型转换源码分析](chapter26_strconv/strconv.md)
---
## [第二十七章 image 图片处理](chapter27_image/image.md)
---

## [第二十八章 如何进行测试](chapter28_test/test.md)
- [1 testing](chapter28_test/01_testing/testing.md) 
  - 1.1 sub 测试并发
  - [1.2 testing.M 将测试交给TestMain调度](chapter28_test/01_testing/02_m/m_test.go)
  - [1.3 testing.F 模糊测试](chapter28_test/01_testing/03_f/sum_test.go)
- [2 go-mock接口测试](chapter28_test/02_gomock/gomock.md)
- 3 web 测试
    - [3.1 使用标准包 httptest 进行 server handler 测试](chapter28_test/03_httptest/01_httptest/httptest.md)
    - [3.2 gock 模拟HTTP流量](chapter28_test/03_httptest/02_gock/gock.md) 
- 4 数据库测试
    - [4.1 sqlmock](chapter28_test/04_database/01_go-sqlmock/go-sqlmock.md)
    - [4.2 miniredis](chapter28_test/04_database/02_miniredis/miniredis.md)
- [5 测试框架 ginkgo-->k8s 使用](chapter28_test/05_ginkgo/ginkgo.md)
- [6 gomonkey 打桩测试(暂不支持arm)](chapter28_test/01_gomonkey/gomonkey.md)
- [7 测试框架 goconvey](chapter28_test/07_goconvey/goconvey.md)
- [8 测试框架testify-->gin 使用)](chapter28_test/08_testify/testify.md)
  - [8.1 assert断言](chapter28_test/08_testify/01_assert/calculate_test.go)
  - [8.2 mock测试替身](chapter28_test/08_testify/02_mock/main_test.go)
  - [8.3 suite测试套件](chapter28_test/08_testify/03_suite/suite_test.go)
- [9 testcontainers 使用容器依赖进行集成测试](chapter28_test/09_testcontainers/testcontainers.md)
  - [9.1 使用 mysql 容器](chapter28_test/09_testcontainers/01_mysql/main_test.go)

## 第二十九章 module包管理
- [1 go-module 实践篇](chapter29_module/01_use/module_operation.md)
  - 模块缓存
  - GOPROXY
- [2 go-module 原理篇](chapter29_module/02_discipline/module.md)
  - Minimal Version Selection 最小版本选择算法
- [3 go1.16 retract 撤回版本](chapter29_module/03_go1_16_module/module.md)
- [4 go1.17 module 依赖图修剪及延迟 module 加载](chapter29_module/04_go1_17_module/module.md)
- [5 go1.18 workspace 工作区模式-->k8s 使用](chapter29_module/05_go1_18_workspace/workspace.md)

---

## 第三十章 内存管理
- 1 Linux内存及Go内存结构管理
  - [1.1 Linux 内存管理](chapter30_memory_management/01_memory/linux_mem.md)
  - [1.2 Go 内存结构管理](chapter30_memory_management/01_memory/go_mem.md)
- [2 GC垃圾回收机制](chapter30_memory_management/02_GC/gc.md)
  - [2.1 下次GC的时机](chapter30_memory_management/02_GC/01_next_gc_stage/main.go)
  - [2.2 删除Map元素时通过 runtime.MemStats 查看GC回收流程](chapter30_memory_management/02_GC/02_map_GC/main.go)
  - 2.3 内存对象中有指针与无指针的GC对比,检测内存对象中的指针
    - [2.3.1 gc运行时间: 切片中存储10亿个指针](chapter30_memory_management/02_GC/03_gc_between_pointer_and_not/01_with_pointer/pointer.go)
    - [2.3.2 gc运行时间: 切片中存储10亿个非指针](chapter30_memory_management/02_GC/03_gc_between_pointer_and_not/02_without_pointer/not_pointer.go)
- [3 逃逸分析](chapter30_memory_management/03_escape_to_heap/escape_to_heap.md)
    - [3.1 argument content escapes(fmt参数内容逃逸)](chapter30_memory_management/03_escape_to_heap/01_fmt_interface.go)
    - [3.2 局部变量指针返回时被外部引用](chapter30_memory_management/03_escape_to_heap/02_params_ptr_return.go)
    - [3.3 接口类型](chapter30_memory_management/03_escape_to_heap/03_interface_method.go)
- [4 内存模型:happened before](chapter30_memory_management/05_happened_before/happened_before.md)
---   
## [第三十一章 go开发套件](chapter31_tool/go_toolsets.md)
- [1 build == compile编译 + link链接，附Go包导入路径讲解](chapter31_tool/01_build/build.md)
  - [1.1 Go build 构建约束（build constraint），也叫做构建标记（build tag）](chapter31_tool/01_build/01_tags/tags.md)
  - [1.2 Go build 选项给 go 链接器传入参数 -ldflags="-X key=value来重写一个符号定义"-->符号表应用](chapter31_tool/01_build/02_ldflags/build.go)
  - [1.3 Go build 选项 -n 查看构建过程用到的命令](chapter31_tool/01_build/03_n/hello.go)
- [2 Go tool 自带工具](chapter31_tool/02_tool/tool.md)
  - [2.1 tool compile 编译](chapter31_tool/02_tool/01_compile/compile.md)
  - [2.2 tool link 链接](chapter31_tool/02_tool/02_link/link.md)
  - [2.3 generate 批量执行任何命令](chapter31_tool/02_tool/03_generate/genarate.md)
- [3 Golang程序调试工具: delve-->调试器分类及实现](chapter31_tool/03_delve/delve.md)
- [4 ast 抽象语法树](chapter31_tool/04_ast/ast.md)
  - [4.1 scanner 词法分析，生成 Token](chapter31_tool/04_ast/01_scan/main.go)
  - 4.2 ast.Expr 
  - 4.3 ImportSpec 打印
  - [4.4 yacc+手动定义词法解析器 lexer-->prometheus 应用](chapter31_tool/04_ast/04_yacc/yacc.md)
    - 实现美国手机号码格式解析
- [5 go1.21 toolchain 规则](chapter31_tool/05_toolchain_in_go1.21/toolchain.md)
---
## [第三十二章 Generic 泛型](chapter32_generic/generic.md)
- [1 泛型在算法上的基本使用](chapter32_generic/01_basic_algorithm_application/main.go)
- [2 ~int 底层类型及 Go1.18 接口分为两种类型: 基本接口(Basic interface) 和 一般接口(General interface)](chapter32_generic/02_typeParam_n_typeArgument/main.go)
- [3 泛型性能测试](chapter32_generic/03_performance/generic_test.go)
- [4 comparable 要求结构体中每个成员变量都是可比较](chapter32_generic/04_comparable/main.go)
- [5 泛型工具库-->samber/lo](chapter32_generic/05_samber_lo/samber_lo.md)
---
## [第三十三章 makefile 使用](chapter33_makefile/Makefile_info.md)
- [1 Makefile常用函数列表](chapter33_makefile/makefile_func.md)
- [2 golang makefile 最佳实践](chapter33_makefile/go_makefile.md)
---
## [第三十四章 regexp 正则表达式 Regular Expression](chapter34_regexp/regexp.md)
- [1 基本正则表达式使用](chapter34_regexp/01_basic_grammar/main.go)
---

## [第三十五章 编码 Unicode](chapter35_unicode/unicode.md)
---
## 第三十六章 unique 包--go 1.23
- [1 字符串驻留（string interning)](chapter36_unique/01_interning/interning.md)

