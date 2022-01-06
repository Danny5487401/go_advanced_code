# ***高级Goland学习代码*** _go_advanced_code_
![](https://changkun.de/urlstat?mode=github&repo=)
[![Go Report Card](https://goreportcard.com/badge/github.com/talkgo/night?style=flat-square)](https://goreportcard.com/report/github.com/Danny5487401/go_advanced_code)
[![GitHub stars](https://img.shields.io/github/stars/talkgo/night.svg?label=Stars&style=flat-square)](https://github.com/Danny5487401/go_advanced_code)
[![GitHub forks](https://img.shields.io/github/forks/talkgo/night.svg?label=Fork&style=flat-square)](https://github.com/Danny5487401/go_advanced_code)
![](https://img.shields.io/github/contributors/talkgo/night.svg?style=flat-square&color=orange&label=all%20contributors)
[![Documentation](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](http://godoc.org/github.com/Danny5487401/go_advanced_code)
[![GitHub issues](https://img.shields.io/github/issues/talkgo/night.svg?label=Issue&style=flat-square)](https://github.com/Danny5487401/go_advanced_code/issues)
![](https://changkun.de/urlstat?mode=github&repo=Danny5487401/go_advanced_code)
[![license](https://img.shields.io/github/license/talkgo/night.svg?style=flat-square)](https://github.com/Danny5487401/go_advanced_code/blob/master/LICENSE)

![高级go编程](./img/golang.jpeg)

# *目录*

## *goVersion==1.16*

## [第一章 I/O操作](chapter01_input_output/io.md)
- 1 os操作系统模块  
    - 1.1 os中FileInfo底层的文件描述符和相关信息   
    - 1.2 os文件操作   
    - 1.3 io包底层Reader和Writer接口   
        - 1.3.1 os,bytes,strings包   
    - 1.4 io断点续传  
    - 1.5 FilePath包     
        - 1.5.1 walkPath遍历目录及文件  
        - 1.5.2 匹配文件名  
- 2 [bufio缓存读写](chapter01_input_output/02_bufio/bufio.md)

---
## 第二章 协程Goroutine
- 1 [线程模型分类及Goroutine切换原则(GPM模型)](chapter02_goroutine/01_GPM/GPM.md)
    - 1.1 [trace查看宏观调度流程(Goroutine启动时长)](chapter02_goroutine/01_GPM/trace/trace.md)
- 2 runtime模块和GC
    - 2.1 runtime核心功能及系统信息调用
    - 2.2 Goexit()终止线程
    - 2.3 资源竞争一致性问题分析
    - [2.4 垃圾回收机制(trace查看map垃圾回收）](chapter02_goroutine/02_runtime/04GC/gc.md)
    - [2.5 监控代码性能pprof](chapter02_goroutine/02_runtime/05pprof/intro.md)
        - 2.5.1 标准包runtime/pprof及net/http/pprof
        - 2.5.2 第三方包pkg/profile
    - [2.6 Go内存结构](chapter02_goroutine/02_runtime/06memory/mem.md)
- 3 [多goroutine的缓存一致性(涉及cpu伪共享)](chapter02_goroutine/03_cache/cache.md)
- 4 [线程池(池化技术)](chapter02_goroutine/04_concurrent_pool/pool.md)
    - 4.1 Goroutine最大数量限制(令牌桶方式)
    - 4.2 百万请求处理
    - [4.3 第三方包线程池ants](chapter02_goroutine/04_concurrent_pool/03_antsPool/ants.md)
    - [4.4 标准库连接池sql实现](chapter02_goroutine/04_concurrent_pool/04_database_sql/sql.md)
- [5 goroutine泄漏分析及处理](chapter02_goroutine/05_goroutine_leaks/goroutine_leak.md)
---

## 第三章 通道Channel
- 1 [Channel内部结构及源码分析(含PPT分析)](chapter03_channel/01_channel_use/channel.md)
    - 1.0 channel初始化
    - 1.1 无缓存通道
    - 1.2 父子通信
    - 1.3 死锁
    - [1.4 优雅关闭channel](chapter03_channel/01_channel_use/04channelClose/ChanClose.md)
    - 1.5 通道遍历range
    - 1.6 有缓冲channel增强并发
    - 1.7 双向通道
    - 1.8 单向通道
    - [1.9 使用channel传递channel](chapter03_channel/01_channel_use/09ChanPassChan/main.go)
    - [1.10 happened before](chapter03_channel/01_channel_use/10_happened_before/happened_before.md)
    - 1.11 读取关闭的通道值
    - [1.12 select中实现channel优先级-->k8s中实现](chapter03_channel/01_channel_use/12_priority_channel/priority_chan.md)
- 2 [channel应用:TimerChan模块](chapter03_channel/02_TimerChan/timer.md)
    - [2.1 reset陷阱](chapter03_channel/02_TimerChan/01_TimerReset/timer_reset.md)
    - [2.2 timerStop使用](chapter03_channel/02_TimerChan/02_TimerStop/timer_stop.md)
    - 2.3 TimerAfter陷阱
- 3 [Select多路复用](chapter03_channel/03_select/03Select_DataStructure/select.md)
- 4 [CSP理论中的Process/Channel](chapter03_channel/04_CSP/CSP.md)
---

## 第四章 interface和反射 
- 1 interface
    - [1.1 interface源码分析](chapter04_interface_n_reflect/01_original_code/interface.md)
- 2 [反射](chapter04_interface_n_reflect/02_reflect/reflect.md)
    - [2.1 反射三大定律](chapter04_interface_n_reflect/02_reflect/01three_laws/threeLaw.md)
    - [2.2 类型断言](chapter04_interface_n_reflect/02_reflect/02TypeAssert/type_assertion.md)
    - [2.3 获取结构体字段及获取方法](chapter04_interface_n_reflect/02_reflect/03StructField_n_method/main.go)
    - [2.4 reflect.Value修改值，调用结构体带方法，调用普通方法](chapter04_interface_n_reflect/02_reflect/04reflectValue/main.go)
    - [2.5 反射性能优化演变案例](chapter04_interface_n_reflect/02_reflect/05PerformanceInprove/main.go)
---

## 第五章 切片和数组
- 1 值传递-数组
- 2 引用传递-指针切片和指针数组
- 3 切片和数组参数传递性能对比
- 4 切片底层结构
- [5 nil切片和空切片](chapter05_slice_n_array/05nilSlice_n_NoneSlice/nil_n_empty_slice.md)
- [6 扩容策略](chapter05_slice_n_array/06GrowSlice/grow_size_policy.md)
- [7 不同类型的切片间互转](chapter05_slice_n_array/07Transfer_slice_in_different_type/main.go)
- [8 带索引初始化数组和切片](chapter05_slice_n_array/08_make_slice_with_index/make_slice_with_index.go)
- [9 使用copy替代原切片进行切片](chapter05_slice_n_array/09_reslice/reslice.go)

---

## 第六章 指针
- 1 指针类型转换及修改值
- [2 指针分类及unsafe包使用](chapter06_pointer/02unsafe/unsafe.md)
- [3 获取并修改结构体私有变量值](chapter06_pointer/03PointerSetPrivateValue/main.go)
- [4 切片与字符串零拷贝互转(指针和反射方式)](chapter06_pointer/04SliceToString/sliceToString.go)

---

## [第七章 系统调用](chapter07_system_call/Syscall.md)
- 1 自定义kqueue服务器（涉及各种linux系统调用）

## [第八章 defer函数及汇编语言理解](chapter08_defer/defer.md)
- 1 注册延迟调用机制定义及使用
- 2 defer陷阱
- 3 分解defer函数
- 4 defer循环性能问题
- 5 [汇编理解defer函数](chapter08_defer/05_defer_assembly/defer_asm.md)

## [第九章 设计模式-OOP七大准则](chapter09_design_pattern/introduction.md)
- 1 创建型模式
    - [1.1 静态工厂模式-->new关键字函数实现简单工厂](chapter09_design_pattern/01_construction/01_StaticFactoryMethod/static_factory.md)
    - 1.2 工厂方法模式-->k8s中实现
    - [1.3 单例模式-->标准库strings/replace实现](chapter09_design_pattern/01_construction/03_singleton/singleton.md)
    - [1.4 原型模式-->confluent-kafka中map实现](chapter09_design_pattern/01_construction/04_prototype/prototype.md)
    - [1.5 建造者模式-->xorm,k8s,zap中实现](chapter09_design_pattern/01_construction/05_builder/builder_info.md)
- 2 结构型模式
    - 2.1 组合模式
        - 2.1.1 修改前：使用面向对象处理
        - 2.1.2 修改后：使用组合模式处理
    - 2.2 [装饰模式-->grpc源码体现](chapter09_design_pattern/02_structure/02_Decorate/decorate.md)
        - 2.2.1 闭包实现--多个装饰器同时使用
        - 2.2.2 结构体装饰
        - 2.2.3 反射实现--泛型装饰器
    - 2.3 [享元模式-->线程池,缓存思想](chapter09_design_pattern/02_structure/03_FlyweightPattern/flyWeightPattern.md)
    - 2.4 适配器模式
    - [2.5 桥接模式](chapter09_design_pattern/02_structure/05_bridgeMethod/bridge_method.md)
    - 2.6 [门面模式(外观模式)-->在gin中render应用(封装多个子服务)](chapter09_design_pattern/02_structure/06_facade_pattern/facade.md)
    - 2.7 [代理模式](chapter09_design_pattern/02_structure/07_proxy/proxy.md)
- 3 行为型模式
    - 3.1  [访问者模式-->k8s中kubectl实现](chapter09_design_pattern/03_motion/01_visitor/vistor.md)
    - 3.2  迭代器-->标准库container/ring中实现
    - 3.3  [状态模式](chapter09_design_pattern/03_motion/03_State/introduction.md)
    - 3.4  责任链模式
    - 3.5  模版模式
    - 3.6  策略模式-->if-else的另类写法(内部算法封装)
    - 3.7  解释器模式
    - 3.8  命令模式-->go-redis中实现
    - 3.9  备忘录模式
    - 3.10 观察者模式-->etcd的watch机制
    - 3.11 中介者模式
- [4 函数选项:成例模式-->在日志库zap中实现](chapter09_design_pattern/04_fuctional_option/option.md)
    - 4.1 未使用的现状
    - 4.2 区分必填项和选项
    - 4.3 带参数的选项模式
- 5 [插件式编程-->grpc中实现](chapter09_design_pattern/05_plugin_programming/plugin.md)
- 6 [同步模式(sync同步原语以及扩展原语)](chapter09_design_pattern/06_Synchronization_mode/01_sync/sync.md)
    - 6.1 waitGroup同步等待组对象
    - 6.2 [互斥锁（sync.Mutex）和读写锁（sync.RWMutex）性能对比](chapter09_design_pattern/06_Synchronization_mode/01_sync/02RWMutex_vs_mutex/mutex.md)
    - [6.3 Once源码分析](chapter09_design_pattern/06_Synchronization_mode/01_sync/03Once/once.md)
    - 6.4 [并发安全Map(读多写少)](chapter09_design_pattern/06_Synchronization_mode/01_sync/04map/sync_map.md)
    - 6.5 [Pool对象池模式( *非连接池* !）-->官方包对象池fmt](chapter09_design_pattern/06_Synchronization_mode/01_sync/05Pool/pool.md)
        - 6.5.1 未使用newFunc
        - 6.5.2 newFunc与GC（附源码分析）
        - 6.5.3 何时使用对象缓存
        - 6.5.4 第三方对象池object pool(bytebufferpool)
    - 6.6 [Cond条件变量及NoCopy机制](chapter09_design_pattern/06_Synchronization_mode/01_sync/06Cond/Cond.md)
    - 6.7 [atomic原子操作源码分析-->zerolog源码中实现](chapter09_design_pattern/06_Synchronization_mode/01_sync/07Atomic/atomic.md)
        - 6.7.0 Value的load和store
        - 6.7.1 add
        - 6.7.2 cas算法和自旋锁
        - 6.7.3 load和store用法
        - 6.7.4 swap交换
    - 6.8 [ErrorGroup获取协程中error](chapter09_design_pattern/06_Synchronization_mode/01_sync/08ErrorGroup/errGroup.md)
    - 6.9 [信号量Semaphore](chapter09_design_pattern/06_Synchronization_mode/01_sync/09Semaphore/semaphore.md)
    - 6.10 SingleFlight避免缓存击穿

## [第十章 函数式编程](chapter10_function/func.md)
- 1 函数
    - 1.1 闭包基本使用
    - 1.2 匿名函数:回调函数
    - 1.3 函数模版:定义行为
- 2 [高级函数](chapter10_function/02_advanced_function/introduction.md)
    - 2.1 简单实现filter,map,reduce
    - 2.2 简单案例
    - 2.3 复杂实现：泛型及类型检查
- 3 一等公民案例
    - [网络管理中问题需求](chapter10_function/03_Firstclassfunction/problem_desc.md)
    - 网络管理中三种处理对比
        - 3.1 通过同享内存通信
        - 3.2 通过通信(具体数据)共享内存
        - 3.3 通过通信(函数)共享内存

## 第十一章 汇编理解go语言底层源码
- 1 [汇编基本指令](chapter11_assembly_language/01asm/introduction.md)
- 2 [ plan9汇编](chapter11_assembly_language/02plan9/introduction.md)
    - 2.1  常量constant
    - 2.2  array类型
    - 2.3  bool类型
    - 2.4  int,int32,uint32类型
    - 2.5  float32，float64类型
    - 2.6  slice类型([]byte)
    - 2.7  引用类型map和channel
    - 2.8  函数类型
    - 2.9  局部变量
    - 2.10 流程控制
    - 2.11 伪SP,FP及硬件SP关系
    - 2.12 结构体方法
    - 2.13 递归函数
    - 2.14 闭包函数

- 3 [ Golang底层数据结构](chapter11_assembly_language/03Golang_data_structure/data.md)
    - [3.1 Map底层结构](chapter11_assembly_language/03Golang_data_structure/map_structure/map_intro.md)
    - [3.2 String底层结构,字符集和字符编码,性能分析及内存泄漏分析](chapter11_assembly_language/03Golang_data_structure/string_structure/str.md)
    - [3.3 Struct底层结构,内存布局,空结构体](chapter11_assembly_language/03Golang_data_structure/struct_structure/struct.md)

## 第十二章 网络编程net
- [socket介绍](chapter12_net/socket.md)
- [tcp介绍](chapter12_net/tcp.md)
- [多路复用](chapter12_net/io_multiplexing.md)

- 1 http服务端高级封装演变
  - [1.1 使用DefaultServeMux](chapter12_net/01_http_server/01_use_DefaultServeMux/main.go)
  - [1.2 使用内置serveMux生成函数](chapter12_net/01_http_server/02_use_http_NewServeMux/main.go)
  - [1.3 自定义实现serveMux](chapter12_net/01_http_server/03_use_cutomized_mux/main.go)
- 2 爬虫获取邮箱案例(http客户端源码分析)
  - [2.1 request源码](chapter12_net/02http_client/http_request.md)
  - [2.2 response源码](chapter12_net/02http_client/http_response.md)
  - [2.3 transport源码](chapter12_net/02http_client/http_transport.md)
- 3 Tcp实现客户端及服务端(tcp底层原理分析)
- 4 Tcp黏包分析及处理(大小端介绍)
- 5 [fastHttp(源码分析)](chapter12_net/05_fasthttp/fasthttp.md)
  - 5.1 服务端
  - 5.2 客户端
- [6 优雅退出原理分析（go-zero实践）](chapter12_net/06_grateful_stop/grateful_stop.md)
- [7 URL的解析 Parse，query 数据的转义与反转义](chapter12_net/07_url/url.md)

## [第十三章 CGO调用C语言](chapter13_Go_call_C_or_C++/introduction.md)
[cgo在confluent-kafka-go源码使用](https://github.com/Danny5487401/go_grpc_example/blob/master/03_amqp/02_kafka/02_confluent-kafka/confluent_kafka_source_code.md)
- 1 Go调用自定义C函数
- 2 Go调用模块化C库
- 3 Go实现C定义函数
- 4 Go获取C函数的errno
- 5 C的void返回

## [第十四章 Context上下文](chapter14_context/introduction.md)
- [0 父类EmptyCtx](chapter14_context/00_original_code_of_context/empty.md)
- 1 Context来源
    - 1.1 问题：如何释放资源
    - 1.2 方式一：全局参数
    - 1.3 方式二: 通道channel
    - 1.4 方式三: Context
- 2 WithCancel源码及使用
- 3 WithDeadline源码及使用
- 4 WithValue源码及使用
- 5 WithTimeout源码及使用

## 第十五章 接口编程
- 1 冗余代码写法
- 2 简单优化
- 3 更优方式

## 第十六章 并发编程
- 1 简单流水线模型
- 2 FAN-IN和FAN-OUT模型

## 第十七章 数据结构及算法
- 1 queue队列
- 2 [加解密](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/encryption.md)
  - 2.1 对称式加密
    - [aes高级加密标准](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/01_symmetric_encryption/aes/aes.md)
      - Cipher FeedBack密码反馈模式
      - Cipher Block Chaining密码分组链接模式
    - [des美国数据加密标准](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/01_symmetric_encryption/des/des.md)
  - 2.2 数字签名(hmac,md5,sha1)
    - [md5](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/02_digital_signature/md5/md5.md)
    - [hmac](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/02_digital_signature/hmac/hmac.md)
    - [sha-1](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/02_digital_signature/sha1/sha1.md)
  - 2.3 [非对称加密算法rsa](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/03_rsa/rsa.md)
    - 分段与不分段加解密
- [3 随机算法（伪随机和真随机)](chapter17_dataStructure_n_algorithm/03_rand/rand.md)
- [4 排序算法分类及图解(sort包源码分析)](chapter17_dataStructure_n_algorithm/04_sort/algorithm.md)
  - [4.1 map排序 ](chapter17_dataStructure_n_algorithm/04_sort/sorted_map/map_sort.go)
  - [4.2 排序接口实现（反射方式）](chapter17_dataStructure_n_algorithm/04_sort/sortByReflect/sort.go)
  - [4.3 slice排序](chapter17_dataStructure_n_algorithm/04_sort/sorted_slice/main.go)
- [5 Jwt源码分析及中间件使用](chapter17_dataStructure_n_algorithm/05_middleware/jwt.md)
- [6 Privacy Enhanced Mail Certificate (pem文件)生成](chapter17_dataStructure_n_algorithm/06_pem_generate/main.go)
- [7 Base64编码解析](chapter17_dataStructure_n_algorithm/07_base64_encoding/base64.md)

## 第十八章 错误跟踪和panic
- 0 错误(err)和异常（exception）区别及处理方式
- 1 自定义错误类型打印错误栈
- [2 扩展包pkg.errors](chapter18_error_n_panic/02_pkg_errors/pkg_errors.md)
- [3 Gin的错误recover分析(panic和recover源码分析)](chapter18_error_n_panic/03_recover/panic.md)
- [4 errCode错误码自动化生成](chapter18_error_n_panic/04_errorCode/02generate_n_stringer/intro.md)
- 5 error如何比较

## 第十九章 nil预定义标识
- 1 不同类型为nil时的地址和大小
- 2 不同类型与nil的比较
  - interface
  - nil==nil
  - ptr,channel,func,map
  - slice
- 3 不同类型nil时的特点

## [第二十章 for-range源码分析](chapter20_for_range/for_range.md)
- 1 遍历数组和切片
- 2 遍历Goroutine(协程启动时间)
- 3 遍历Map

## [第二十一章 time标准包源码分析](chapter21_time/time.md)
- 1 比time.Now()更优雅获取时间戳（go:link技术）

## [第二十二章 数据驱动模板-kratos工具生成](chapter22_template/text_template.md)
- 1 加载多个模版
- 2 自定义模版函数
