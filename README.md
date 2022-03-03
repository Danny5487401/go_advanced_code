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
Note:目录同级为 代码展示，可在Goland中运行
## *goVersion==1.16*

## [第一章 I/O操作](chapter01_input_output/io.md)
- 1 os操作系统模块
    - [1.1 os中FileInfo底层的文件描述符和相关信息](chapter01_input_output/01_OS_module/01FileInfo/main.go)   
    - [1.2 os文件操作](chapter01_input_output/01_OS_module/02FileOperation/main.go)   
    - 1.3 io包底层Reader和Writer接口   
        - 1.3.1 os,bytes,strings包   
    - [1.4 io断点续传 ](chapter01_input_output/01_OS_module/04seeker/02resume_from_break-point/main.go) 
    - [1.5 FilePath包 ](chapter01_input_output/01_OS_module/05filePath/walk.go)    
        - [1.5.1 walkPath遍历目录及文件(匹配文件名)](chapter01_input_output/01_OS_module/05filePath/walk.go)
- [2 bufio缓存读写](chapter01_input_output/02_bufio/bufio.md)
  - [2.1 reader](chapter01_input_output/02_bufio/01reader/main.go)
  - [2.2 writer](chapter01_input_output/02_bufio/02writer/main.go)

---
## 第二章 协程Goroutine
- 1 [线程模型分类及Goroutine切换原则(GPM模型)](chapter02_goroutine/01_GPM/GPM.md)
    - 1.1 [trace查看宏观调度流程(Goroutine启动时长)](chapter02_goroutine/01_GPM/trace/trace.md)
- 2 [runtime模块和GC](chapter02_goroutine/02_runtime/runtime.md)
    - [2.1 runtime核心功能及系统信息调用](chapter02_goroutine/02_runtime/01basic_use/main.go)
    - [2.2 Goexit()终止线程](chapter02_goroutine/02_runtime/02GoExit/main.go)
    - [2.3 资源竞争一致性问题分析](chapter02_goroutine/02_runtime/03ResourceCompetition/01problem/resource_competion.md)
      - [2.3.1 问题产生](chapter02_goroutine/02_runtime/03ResourceCompetition/01problem/main.go)
      - [2.3.2 问题解决](chapter02_goroutine/02_runtime/03ResourceCompetition/02Fix_Resource_data_consistency/main.go)
    - [2.4 GC垃圾回收机制(trace查看map垃圾回收)](chapter02_goroutine/02_runtime/04GC/gc.md)
      - [2.4.1 下次GC的时机](chapter02_goroutine/02_runtime/04GC/01_next_gc_stage/main.go)
      - [2.4.2 删除Map元素查看GC回收流程](chapter02_goroutine/02_runtime/04GC/02_map_GC/main.go)
    - [2.5 监控代码性能pprof](chapter02_goroutine/02_runtime/05pprof/intro.md)
      - [2.5.1 标准包runtime/pprof及net/http/pprof使用](chapter02_goroutine/02_runtime/05pprof/01_pprof/main.go)
      - [2.5.2 第三方包pkg/profile](chapter02_goroutine/02_runtime/05pprof/02_pkg_profile/cpu.go)
    - 2.6 Linux内存及Go内存结构管理
      - [2.6.1 Linux内存管理](chapter02_goroutine/02_runtime/06memory/linux_mem.md)
      - [2.6.2 Go内存结构管理](chapter02_goroutine/02_runtime/06memory/go_mem.md)
    - [2.7 prometheus监控程序](chapter02_goroutine/02_runtime/07prometheus/prometheus.md)
- 3 [多goroutine的缓存一致性(涉及cpu伪共享)](chapter02_goroutine/03_cache/cache.md)
- 4 [线程池(池化技术)](chapter02_goroutine/04_concurrent_pool/pool.md)
    - [4.1 Goroutine最大数量限制(令牌桶方式)](chapter02_goroutine/04_concurrent_pool/01_goroutine_max_control/main.go)
    - [4.2 百万请求处理](chapter02_goroutine/04_concurrent_pool/02_millionRequests/main.go)
    - [4.3 第三方包线程池ants](chapter02_goroutine/04_concurrent_pool/03_antsPool/ants.md)
    - [4.4 标准库连接池sql实现](chapter02_goroutine/04_concurrent_pool/04_database_sql/sql.md)
- [5 goroutine泄漏分析及处理](chapter02_goroutine/05_goroutine_leaks/goroutine_leak.md)
---

## 第三章 通道Channel
- 1 [Channel内部结构及源码分析(含PPT分析)](chapter03_channel/01_channel_use/channel.md)
    - [1.0 channel初始化](chapter03_channel/01_channel_use/00introdution/main.go)
    - [1.1 无缓存通道](chapter03_channel/01_channel_use/01unbufferd_channel/main.go)
    - [1.2 父子通信](chapter03_channel/01_channel_use/02ParentChildrenCommunication/main.go)
    - [1.3 死锁](chapter03_channel/01_channel_use/03deadlock/main.go)
    - [1.4 优雅关闭channel](chapter03_channel/01_channel_use/04channelClose/ChanClose.md)
    - [1.5 通道遍历range](chapter03_channel/01_channel_use/05ChannelRange/main.go)
    - [1.6 有缓冲channel增强并发](chapter03_channel/01_channel_use/06bufferChan/main.go)
    - [1.7 双向通道](chapter03_channel/01_channel_use/07two-wayChan/main.go)
    - [1.8 单向通道](chapter03_channel/01_channel_use/08one-wayChan/main.go)
    - [1.9 使用channel传递channel](chapter03_channel/01_channel_use/09ChanPassChan/main.go)
    - [1.10 happened before](chapter03_channel/01_channel_use/10_happened_before/happened_before.md)
    - [1.11 循环读取关闭的通道值是否阻塞](chapter03_channel/01_channel_use/11_read_closed_chan/readCloseChan.go)
    - [1.12 select中实现channel优先级-->k8s中实现](chapter03_channel/01_channel_use/12_priority_channel/priority_chan.md)
- 2 [channel应用:TimerChan模块](chapter03_channel/02_TimerChan/timer.md)
    - [2.1 reset陷阱](chapter03_channel/02_TimerChan/01_TimerReset/timer_reset.md)
    - [2.2 timerStop使用](chapter03_channel/02_TimerChan/02_TimerStop/timer_stop.md)
    - [2.3 TimerAfter陷阱](chapter03_channel/02_TimerChan/03_TimeAfter/main.go)
- 3 [Select多路复用](chapter03_channel/03_select/03Select_DataStructure/select.md)
- 4 [CSP理论中的Process/Channel](chapter03_channel/04_CSP/CSP.md)
---

## 第四章 interface和反射 
- [1 interface源码分析](chapter04_interface_n_reflect/01_interface/interface.md)
    - [1.1 汇编分析不含方法eface和带方法iface](chapter04_interface_n_reflect/01_interface/01_interface_in_asm/main.go)
    - [1.2 接口值的零值是指动态类型和动态值都为 nil](chapter04_interface_n_reflect/01_interface/02_interface_compare_with_nil/main.go)
    - [1.3 打印出接口的动态类型和值](chapter04_interface_n_reflect/01_interface/03_print_dynamic_value_n_type/main.go)
- 2 [反射](chapter04_interface_n_reflect/02_reflect/reflect.md)
    - [2.1 反射三大定律](chapter04_interface_n_reflect/02_reflect/01three_laws/threeLaw.md)
    - [2.2 类型断言](chapter04_interface_n_reflect/02_reflect/02TypeAssert/type_assertion.md)
    - [2.3 获取结构体字段及获取方法](chapter04_interface_n_reflect/02_reflect/03StructField_n_method/main.go)
    - [2.4 reflect.Value修改值，调用结构体带方法，调用普通方法](chapter04_interface_n_reflect/02_reflect/04reflectValue/main.go)
    - [2.5 反射性能优化演变案例](chapter04_interface_n_reflect/02_reflect/05PerformanceInprove/main.go)
---

## 第五章 切片和数组
- [1 值传递-数组](chapter05_slice_n_array/01passByValue_array/main.go)
- [2 引用传递-指针切片和指针数组](chapter05_slice_n_array/02passByReference/main.go)
- [3 切片和数组参数传递性能对比](chapter05_slice_n_array/03Array_n_slice_performance/main_test.go)
- 4 底层数据结构
  - [切片](chapter05_slice_n_array/04structure_of_array_n_slice/slice/sliceStructure.md)
  - [数组](chapter05_slice_n_array/04structure_of_array_n_slice/array/arrayStructure.md)
- [5 nil切片和空切片](chapter05_slice_n_array/05nilSlice_n_NoneSlice/nil_n_empty_slice.md)
- [6 扩容策略](chapter05_slice_n_array/06GrowSlice/grow_size_policy.md)
- [7 不同类型的切片间互转](chapter05_slice_n_array/07Transfer_slice_in_different_type/main.go)
- [8 带索引初始化数组和切片](chapter05_slice_n_array/08_make_slice_with_index/make_slice_with_index.go)
- [9 使用copy替代原切片进行切片](chapter05_slice_n_array/09_reslice/reslice.go)
---

## 第六章 指针
- [1 指针类型转换及修改值](chapter06_pointer/01ptrOperation/main.go)
- [2 指针分类及unsafe包使用](chapter06_pointer/02unsafe/unsafe.md)
- [3 获取并修改结构体私有变量值](chapter06_pointer/03PointerSetPrivateValue/main.go)
- [4 切片与字符串零拷贝互转(指针和反射方式)](chapter06_pointer/04SliceToString/sliceToString.go)
---

## [第七章 系统调用](chapter07_system_call/Syscall.md)
- [1 自定义kqueue服务器（涉及各种linux系统调用）](chapter07_system_call/01_kqueue_server/main.go)
- [2 使用strace工具追踪系统调用](chapter07_system_call/02_ptrace/ptrace.md)

## [第八章 defer函数及汇编语言理解](chapter08_defer/defer.md)
- [1 注册延迟调用机制定义及使用](chapter08_defer/01_defer_definiton/main.go)
- [2 defer陷阱](chapter08_defer/02_defer_common_mistakes/main.go)
- [3 分阶段解析defer函数](chapter08_defer/03_defer_params_n_return/main.go)
- [4 defer循环性能问题](chapter08_defer/04_defer_loop_performance/main.go)
- [5 汇编理解defer函数(AMD)](chapter08_defer/05_defer_assembly/defer_amd.s)

## [第九章 设计模式-OOP七大准则](chapter09_design_pattern/introduction.md)
- 1 创建型模式
    - [1.1 静态工厂模式-->new关键字函数实现简单工厂](chapter09_design_pattern/01_construction/01_StaticFactoryMethod/static_factory.md)
    - [1.2 工厂方法模式-->k8s中实现](chapter09_design_pattern/01_construction/02_factory_mode/factory.md)
    - [1.3 单例模式-->标准库strings/replace实现](chapter09_design_pattern/01_construction/03_singleton/singleton.md)
    - [1.4 原型模式-->confluent-kafka中map实现](chapter09_design_pattern/01_construction/04_prototype/prototype.md)
    - [1.5 建造者模式-->xorm,k8s,zap中实现](chapter09_design_pattern/01_construction/05_builder/builder_info.md)
- 2 结构型模式
    - 2.1 组合模式
        - [2.1.1 修改前：使用面向对象处理](chapter09_design_pattern/02_structure/01_Composite/01_modify_before/composite.go)
        - [2.1.2 修改后：使用组合模式处理](chapter09_design_pattern/02_structure/01_Composite/02_modify_after/conposite.go)
    - 2.2 [装饰模式-->grpc源码体现](chapter09_design_pattern/02_structure/02_Decorate/decorate.md)
        - [2.2.1 闭包实现--多个装饰器同时使用](chapter09_design_pattern/02_structure/02_Decorate/01_closure_decorate/main.go)
        - [2.2.2 结构体装饰](chapter09_design_pattern/02_structure/02_Decorate/02_struct_decorate_inGrpc/main.go)
        - [2.2.3 反射实现--泛型装饰器](chapter09_design_pattern/02_structure/02_Decorate/03_reflect_decorate/decorate.go)
    - 2.3 [享元模式-->线程池,缓存思想](chapter09_design_pattern/02_structure/03_FlyweightPattern/flyWeightPattern.md)
    - [2.4 适配器模式](chapter09_design_pattern/02_structure/04_adopter/adopter.md)
    - [2.5 桥接模式](chapter09_design_pattern/02_structure/05_bridgeMethod/bridge_method.md)
    - 2.6 [门面模式(外观模式)-->在gin中render应用(封装多个子服务)](chapter09_design_pattern/02_structure/06_facade_pattern/facade.md)
    - 2.7 [代理模式](chapter09_design_pattern/02_structure/07_proxy/proxy.md)
- 3 行为型模式
    - [3.1  访问者模式-->k8s中kubectl实现](chapter09_design_pattern/03_motion/01_visitor/vistor.md)
    - [3.2  迭代器-->标准库container/ring中实现](chapter09_design_pattern/03_motion/02_Iterator/main.go)
    - [3.3  状态模式](chapter09_design_pattern/03_motion/03_State/introduction.md)
    - [3.4  责任链模式](chapter09_design_pattern/03_motion/04_duty_chain_method/duty_chain.md)
    - [3.5  模版模式](chapter09_design_pattern/03_motion/05_templateMethod/template.go)
    - [3.6  策略模式-->if-else的另类写法(内部算法封装)](chapter09_design_pattern/03_motion/06_strategyMethod/strategy.md)
    - [3.7  解释器模式](chapter09_design_pattern/03_motion/07_InterpreterMethod/interpreter.md)
    - [3.8  命令模式-->go-redis中实现](chapter09_design_pattern/03_motion/08_CommandMethod/command.md)
    - [3.9  备忘录模式](chapter09_design_pattern/03_motion/09_memento/introduction.md)
    - [3.10 观察者模式-->官方Signal包及etcd的watch机制](chapter09_design_pattern/03_motion/10_ObserverPattern/introduction.md)
    - [3.11 中介者模式](chapter09_design_pattern/03_motion/11_mediator/inctroduction.md)
- [4 函数选项:成例模式-->在日志库zap中实现](chapter09_design_pattern/04_fuctional_option/option.md)
    - [4.1 未使用函数选项初始化结构体的现状](chapter09_design_pattern/04_fuctional_option/01_problem/ServerConfig.md)
    - [4.2 区分必填项和选项](chapter09_design_pattern/04_fuctional_option/02_method_splitConfig/SplitConfig.go)
    - 4.3 带参数的选项模式
      - [不返回error](chapter09_design_pattern/04_fuctional_option/03_FunctionalOption/01_simple_solution/main.go)
      - [返回error](chapter09_design_pattern/04_fuctional_option/03_FunctionalOption/02_complexed_with_error/main.go)
- 5 [插件式编程-->grpc中实现](chapter09_design_pattern/05_plugin_programming/plugin.md)
- 6 [同步模式(sync同步原语以及扩展原语)](chapter09_design_pattern/06_Synchronization_mode/01_sync/sync.md)
    - [6.1 waitGroup同步等待组对象](chapter09_design_pattern/06_Synchronization_mode/01_sync/01waitGroup/waitGroup.md)
    - 6.2 使用互斥锁（sync.Mutex）实现读写功能和直接使用读写锁（sync.RWMutex）性能对比
      - [6.2.1 使用互斥锁（sync.Mutex）实现读写功能](chapter09_design_pattern/06_Synchronization_mode/01_sync/02RWMutex_vs_mutex/Mutex/main.go)
      - [6.2.2 直接使用读写锁（sync.RWMutex）实现读写功能](chapter09_design_pattern/06_Synchronization_mode/01_sync/02RWMutex_vs_mutex/RWMutex/main.go)
      - [Mutex和RWMutex源码分析](chapter09_design_pattern/06_Synchronization_mode/01_sync/02RWMutex_vs_mutex/mutex.md)
    - [6.3 Once源码分析](chapter09_design_pattern/06_Synchronization_mode/01_sync/03Once/once.md)
    - [6.4 并发安全Map(读多写少)](chapter09_design_pattern/06_Synchronization_mode/01_sync/04map/sync_map.md)
    - [6.5 Pool对象池模式( *非连接池* !）-->官方包对象池fmt](chapter09_design_pattern/06_Synchronization_mode/01_sync/05Pool/pool.md)
        - [6.5.1 错误使用：未使用newFunc](chapter09_design_pattern/06_Synchronization_mode/01_sync/05Pool/01Without_newFunc/main.go)
        - [6.5.2 newFunc与GC前后Get对比](chapter09_design_pattern/06_Synchronization_mode/01_sync/05Pool/02NewFunc/newFunc.go)
        - [6.5.3 何时使用对象缓存](chapter09_design_pattern/06_Synchronization_mode/01_sync/05Pool/03When2Use_object_pool/main.go)
        - [6.5.4 第三方对象池object pool(bytebufferpool)](chapter09_design_pattern/06_Synchronization_mode/01_sync/05Pool/04_byteBufferPool/main.go)
    - 6.6 [Cond条件变量通知所有协程及NoCopy机制-->熔断框架hystrix-go实现](chapter09_design_pattern/06_Synchronization_mode/01_sync/06Cond/Cond.md)
    - 6.7 [atomic原子操作源码分析-->zerolog源码中实现](chapter09_design_pattern/06_Synchronization_mode/01_sync/07Atomic/atomic.md)
        - [6.7.0 Value的load和store](chapter09_design_pattern/06_Synchronization_mode/01_sync/07Atomic/00_value/main.go)
        - [6.7.1 add及补码减](chapter09_design_pattern/06_Synchronization_mode/01_sync/07Atomic/01_add/main.go)
        - [6.7.2 cas算法和自旋锁](chapter09_design_pattern/06_Synchronization_mode/01_sync/07Atomic/02_CompareAndSwap/main.go)
        - [6.7.3 load和store用法](chapter09_design_pattern/06_Synchronization_mode/01_sync/07Atomic/03_load_n_store/main.go)
        - [6.7.4 swap交换](chapter09_design_pattern/06_Synchronization_mode/01_sync/07Atomic/04_swap/main.go)
    - [6.8 ErrorGroup获取协程中error](chapter09_design_pattern/06_Synchronization_mode/01_sync/08ErrorGroup/errGroup.md)
        - [6.8.1 不带context](chapter09_design_pattern/06_Synchronization_mode/01_sync/08ErrorGroup/01WithoutContext/main.go)
        - [6.8.2 带context](chapter09_design_pattern/06_Synchronization_mode/01_sync/08ErrorGroup/02WithContext/main.go)
    - [6.9 信号量Semaphore](chapter09_design_pattern/06_Synchronization_mode/01_sync/09Semaphore/semaphore.md)
    - [6.10 SingleFlight避免缓存击穿](chapter09_design_pattern/06_Synchronization_mode/01_sync/10SingleFlight/singleFlight.md)
        - [6.10.1 Do方法](chapter09_design_pattern/06_Synchronization_mode/01_sync/10SingleFlight/01_do/main.go)
        - [6.10.2 DoChan方法](chapter09_design_pattern/06_Synchronization_mode/01_sync/10SingleFlight/02_do_chan/main.go)

## [第十章 函数式编程](chapter10_function/func.md)
- 1 函数
    - [1.1 闭包基本使用](chapter10_function/01_func_application/01_closure/main.go)
    - [1.2 匿名函数应用:回调函数](chapter10_function/01_func_application/02_anonymousFunc/main.go)
    - [1.3 函数模版:定义行为](chapter10_function/01_func_application/03_func_template/main.go)
- 2 [高级函数](chapter10_function/02_advanced_function/introduction.md)
    - 2.1 简单实现filter,map,reduce
    - 2.2 简单案例
    - 2.3 复杂实现：泛型及类型检查
    - [2.4 一个应用首页可能依托于很多服务,在没有强依赖关系下,优雅地实现并发编排任务](chapter10_function/02_advanced_function/04_mapReduce/main.go)
- 3 一等公民案例
    - [网络管理中问题需求](chapter10_function/03_Firstclassfunction/problem_desc.md)
    - 网络管理中三种处理对比
        - [3.1 通过同享内存通信](chapter10_function/03_Firstclassfunction/01_communicate_by_sharing_memory/main.go)
        - 3.2 通过通信(具体数据)共享内存
        - 3.3 通过通信(函数)共享内存

## 第十一章 汇编理解go语言底层源码(AMD芯片运行代码)
- 1 [汇编基本指令](chapter11_assembly_language/01asm/introduction.md)
- 2 [ plan9汇编](chapter11_assembly_language/02plan9/introduction.md)
    - [2.1  常量constant](chapter11_assembly_language/02plan9/01_pkg_constant_string/main.go)
    - [2.2  array数组](chapter11_assembly_language/02plan9/02_pkg_array/main.go)
    - [2.3  bool类型](chapter11_assembly_language/02plan9/03_pkg_bool/main.go)
    - [2.4  int,int32,uint32类型](chapter11_assembly_language/02plan9/04_pkg_int/main.go)
    - [2.5  float32，float64类型](chapter11_assembly_language/02plan9/05_pkg_float/main.go)
    - [2.6  slice切片([]byte)](chapter11_assembly_language/02plan9/06_pkg_slice/main.go)
    - [2.7  引用类型map和channel](chapter11_assembly_language/02plan9/07_pkg_channel_n_map/main.go)
    - [2.8  函数类型](chapter11_assembly_language/02plan9/08_pkg_func/main.go)
    - [2.9  局部变量](chapter11_assembly_language/02plan9/09_local_param/local_params.md)
    - [2.10 流程控制](chapter11_assembly_language/02plan9/10_control_process/main.go)
    - [2.11 伪SP,FP及硬件SP关系](chapter11_assembly_language/02plan9/11_FalseSP_fp_SoftwareSP_relation/main.go)
    - [2.12 结构体方法](chapter11_assembly_language/02plan9/12_struct_method/main.go)
    - [2.13 递归函数](chapter11_assembly_language/02plan9/13_recursive_func/main.go)
    - [2.14 闭包函数](chapter11_assembly_language/02plan9/14_closure/main.go)
    - [2.15 GoroutineId获取](chapter11_assembly_language/02plan9/15_GoroutineId/main.go)
    - [2.16 汇编调用非汇编Go函数](chapter11_assembly_language/02plan9/16_assembly_call_NonassemblyFunc/main.go)
- 3 [ Golang底层数据结构-涉及数值类型占用的bit](chapter11_assembly_language/03Golang_data_structure/data.md)
    - [3.1 Map底层结构](chapter11_assembly_language/03Golang_data_structure/map_structure/map_intro.md)
    - [3.2 String底层结构,字符集和字符编码,性能分析及内存泄漏分析](chapter11_assembly_language/03Golang_data_structure/string_structure/str.md)
    - [3.3 Struct底层结构,内存布局,空结构体内存对齐](chapter11_assembly_language/03Golang_data_structure/struct_structure/struct.md)

## 第十二章 网络编程net
- [socket介绍](chapter12_net/socket.md)
- [tcp介绍](chapter12_net/tcp.md)
- [多路复用](chapter12_net/io_multiplexing.md)

- 1 http服务端高级封装演变
  - [1.1 使用DefaultServeMux](chapter12_net/01_http_server/01_use_DefaultServeMux/main.go)
  - [1.2 使用内置serveMux生成函数](chapter12_net/01_http_server/02_use_http_NewServeMux/main.go)
  - [1.3 自定义实现serveMux](chapter12_net/01_http_server/03_use_cutomized_mux/main.go)
- 2 爬虫获取邮箱案例(http客户端源码分析)
  - [2.1 request源码](chapter12_net/02_http_client/http_request.md)
  - [2.2 response源码](chapter12_net/02_http_client/http_response.md)
  - [2.3 transport源码](chapter12_net/02_http_client/http_transport.md)
- 3 Tcp实现客户端及服务端(tcp底层原理分析)
  - [客户端](chapter12_net/03_tcp/client/main.go)
  - [服务端](chapter12_net/03_tcp/server/main.go)
- 4 [Tcp黏包分析及处理(大小端介绍)](chapter12_net/04_tcp_sticky_problem/big_n_small_endian.md)
  - [4.1 TCP粘包问题](chapter12_net/04_tcp_sticky_problem/01_problem)
  - [TCP粘包解决方式](chapter12_net/04_tcp_sticky_problem/02_solution)
- 5 [fastHttp(源码分析)](chapter12_net/05_fasthttp/fasthttp.md)
  - [5.1 服务端](chapter12_net/05_fasthttp/server/main.go)
  - [5.2 客户端](chapter12_net/05_fasthttp/client/client.go)
- [6 优雅退出原理分析-涉及linux信号介绍（go-zero实践）](chapter12_net/06_grateful_stop/grateful_stop.md)
  - [6.1 信号监听处理](chapter12_net/06_grateful_stop/signal.go)
- [7 URL的解析 Parse，query 数据的转义与反转义](chapter12_net/07_url/url.md)

## [第十三章 CGO调用C语言](chapter13_Go_call_C_or_C++/introduction.md)
[cgo在confluent-kafka-go源码使用](https://github.com/Danny5487401/go_grpc_example/blob/master/03_amqp/02_kafka/02_confluent-kafka/confluent_kafka_source_code.md)

**_Note:内部运行代码待优化_**

- [1 Go调用自定义C函数-未模块化](chapter13_Go_call_C_or_C++/01_call_C_func/main.go)
- 2 Go调用自定义C函数-模块化
- 3 Go重写C定义函数
- 4 cgo引入其他包的变量错误
- 5 #Cgo语句
- [6 Go获取C函数的errno](chapter13_Go_call_C_or_C++/06_return_err/main.go)
- [7 C的void返回](chapter13_Go_call_C_or_C++/07_void_return/main.go)
- 8 Go调用模块化C++库
- [9 调用静态C库](chapter13_Go_call_C_or_C++/09_static_c_lib/main.go)
- 10 调用动态C库
- [11 Golang使用pkg-config自动获取头文件和链接库的方法](chapter13_Go_call_C_or_C++/11_pkg_config/pkg_config.md)

## [第十四章 Context上下文-源码分析涉及父类EmptyCtx](chapter14_context/introduction.md)
- 1 Context使用背景
    - [1.1 问题：如何通过父进程优雅释放子goroutine资源](chapter14_context/01_Reason_To_Use_Context/01_problem/main.go)
    - [1.2 方式一：全局参数方式解决的优缺点](chapter14_context/01_Reason_To_Use_Context/02_Method1_Global_Param/main.go)
    - [1.3 方式二: 通道channel方式解决的优缺点](chapter14_context/01_Reason_To_Use_Context/03_Method2_Channel/main.go)
    - [1.4 方式三: 最优方式Context](chapter14_context/01_Reason_To_Use_Context/04_Method3_Context/main.go)
- 2 WithCancel使用
- 3 WithDeadline使用
- 4 WithValue使用
- 5 WithTimeout使用

## 第十五章 接口嵌套编程
- [1 常见冗余代码写法](chapter15_interfaceProgramming/01_problem/main.go)
- [2 简单优化](chapter15_interfaceProgramming/02_simple_method/main.go)
- [3 更优方式](chapter15_interfaceProgramming/03_better_solution/main.go)

## 第十六章 并发编程
- [1 简单流水线模型](chapter16_concurrentProgramming/01_pipeline/pipeline.md)
- [2 FAN-IN和FAN-OUT模型](chapter16_concurrentProgramming/02_fanin_fanout/fanin_fanout.md)

## 第十七章 数据结构及算法
- 1 [queue双端单向队列(泛型)](chapter17_dataStructure_n_algorithm/01_queue/queue_test.go)
- 2 [加解密](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/encryption.md)
  - 2.1 对称式加密
    - [aes高级加密标准](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/01_symmetric_encryption/aes/aes.md)
      - [Cipher FeedBack密码反馈模式](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/01_symmetric_encryption/aes/01_cfb/aes_cfb.go)
      - [Cipher Block Chaining密码分组链接模式](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/01_symmetric_encryption/aes/02_cbc/aes_cbc.go)
    - [des美国数据加密标准](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/01_symmetric_encryption/des/des.md)
  - 2.2 数字签名(hmac,md5,sha1,sha256)
    - [MD5信息摘要算法](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/02_digital_signature/md5/md5.md)
    - [hmac](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/02_digital_signature/hmac/hmac.md)
    - [SHA安全散列算法secure Hash Algorithm](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/02_digital_signature/sha1_n_sha256/sha.md)
  - 2.3 [非对称加密算法rsa](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/03_rsa/rsa.md)
    - [分段与不分段进行加解密](chapter17_dataStructure_n_algorithm/02_encrypt_n_decript_algorithm/03_rsa/main.go)
- [3 随机算法（伪随机和真随机)](chapter17_dataStructure_n_algorithm/03_rand/rand.md)
  - [3.1 math_rand使用](chapter17_dataStructure_n_algorithm/03_rand/01_math_rand/main.go)
  - [3.2 crypto_rand使用](chapter17_dataStructure_n_algorithm/03_rand/02_crypto_rand/main.go)
  - [3.3 fastrand优化使用](chapter17_dataStructure_n_algorithm/03_rand/03_fastrand/main.go)
- [4 排序算法分类及图解(sort包源码分析)](chapter17_dataStructure_n_algorithm/04_sort/algorithm.md)
  - [4.1 map排序 ](chapter17_dataStructure_n_algorithm/04_sort/sorted_map/map_sort.go)
  - [4.2 排序接口实现（反射方式）](chapter17_dataStructure_n_algorithm/04_sort/sortByReflect/sort.go)
  - [4.3 slice排序](chapter17_dataStructure_n_algorithm/04_sort/sorted_slice/main.go)
- [5 Jwt源码分析及中间件使用](chapter17_dataStructure_n_algorithm/05_middleware/jwt.md)
- [6 Privacy Enhanced Mail Certificate (pem文件)生成](chapter17_dataStructure_n_algorithm/06_pem_generate/main.go)
- [7 Base64编码解析](chapter17_dataStructure_n_algorithm/07_base64_encoding/base64.md)
- [8 trie前缀树](chapter17_dataStructure_n_algorithm/08_trie/trie.md)

## 第十八章 错误跟踪和panic
- [0 错误(err)和异常（exception）区别及处理方式](chapter18_error_n_panic/00_diff_between_err_n_exception/main.go)
- 1 自定义错误类型打印错误栈
- [2 扩展包pkg.errors](chapter18_error_n_panic/02_pkg_errors/pkg_errors.md)
- [3 Gin的错误recover分析(panic和recover源码分析)](chapter18_error_n_panic/03_recover/panic.md)
- [4 errCode错误码自动化生成](chapter18_error_n_panic/04_errorCode/02generate_n_stringer/intro.md)
- [5 error如何正确比较](chapter18_error_n_panic/05_err_comparision/main.go)
- [6 收集多个errors-->go-multierror实现](chapter18_error_n_panic/06_multi_error/01_one_goroutine_n_errors/main.go)
## 第十九章 nil预定义标识
- 1 不同类型为nil时的地址和大小
- 2 不同类型与nil的比较
  - [interface为nil时:数据段和类型](chapter19_nil/02_comparison/interface/interface.go)
  - [nil==nil不可以比较](chapter19_nil/02_comparison/nil/main.go)
  - [ptr,channel,func,map为nil必须地址未分配](chapter19_nil/02_comparison/ptr_chan_func_map/main.go)
  - [slice的长度和容量不决定nil](chapter19_nil/02_comparison/slice/slice.go)
- 3 不同类型nil时的特点
  - [channel为nil时的接收，发送，关闭及select](chapter19_nil/03_Attribute/channel/chan.go)
  - [map为nil时可读不可写](chapter19_nil/03_Attribute/map/map.go)
  - [结构体指针为nil时是否可以调用方法](chapter19_nil/03_Attribute/ptr/ptr.go)

## [第二十章 for-range源码分析](chapter20_for_range/for_range.md)
- [1 遍历数组,切片,结构体数组](chapter20_for_range/01_for_range_slice_n_array/main.go)
- [2 正确遍历Goroutine(解析协程启动时间)](chapter20_for_range/02_for_range_goroutine/main.go)
- [3 遍历Map(增加或删除map元素时)](chapter20_for_range/03_for_range_map/main.go)

## [第二十一章 time标准包源码分析](chapter21_time/time.md)
- [1 比time.Now()更优雅获取时间戳（go:link技术）](chapter21_time/time.go)
## [第二十二章 数据驱动模板-kratos工具生成](chapter22_template/text_template.md)
- [1 加载多个模版](chapter22_template/01_multi_template/main.go)
- [2 自定义模版函数](chapter22_template/02_template_func/main.go)

## 第二十三章 调试内部对象
- [1 fmt打印结构体中含有指针对象, 数组或者map中是指针对象, 循环结构时的困难](chapter23_debug_program/01_fmt_problem/main.go)
- [2 go-spew优化调试](chapter23_debug_program/02_go_spew/main.go)
