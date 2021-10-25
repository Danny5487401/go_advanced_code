# ***高级Goland学习代码*** _go_advanced_code_

![高级go编程](./img/golang.jpeg)
# *目录*
## *goVersion==1.16*
## [第一章 I/O操作](chapter01_input_output/io.md)
    1. os操作系统模块
        1.1 os中FileInfo底层的文件描述符和相关信息
        1.2 os文件操作
        1.3 io包底层Reader和Writer接口
            1.3.1 os,bytes,strings包
        1.4 io断点续传
        1.5 walkPath遍历目录及文件
    2. bufio缓存读写
---
## 第二章 协程Goroutine
    1. 线程模型分类及Goroutine切换原则
        1.1 trace查看宏观调度流程(GPM模型)
    2. runtime模块和GC
        2.1 runtime核心功能及系统信息调用
        2.2 Goexit()终止线程
        2.3 资源竞争一致性问题分析
        2.4 垃圾回收机制(trace查看map垃圾回收）
        2.5 监控代码pprof
            2.5.1 标准包runtime/pprof及net/http/pprof
            2.5.2 第三方包pkg/profile
        2.6 Go内存结构
    3. 多goroutine的缓存一致性(涉及cpu伪共享)
    4. 线程池(池化技术)
        4.1 Goroutine最大数量限制(令牌桶方式)
        4.2 百万请求处理
        4.3 第三方包线程池ants
        4.4 标准库连接池sql实现
    5. goroutine泄漏分析及处理
---
## 第三章 通道Channel
    1. 基本使用
        1.0 channel初始化
        1.1 无缓存通道
        1.2 父子通信
        1.3 死锁
        1.4 优雅关闭channel
        1.5 通道遍历range
        1.6 有缓冲channel增强并发
        1.7 双向通道
        1.8 单向通道
        1.9 使用channel传递channel
        1.10 happened before
        1.11 读取关闭的通道值
    2. channel应用:TimerChan模块
        2.1 reset陷阱
        2.2 timerStop使用
        2.3 TimerAfter陷阱
    3. Select多路复用
    4. CSP理论中的Process/Channel
    5. Channel内部结构及源码分析(含PPT分析)
---
## 第四章 interface和反射
    1. interface
        1.1 源码分析
    2. 反射
        2.1 反射三大定律
        2.2 类型断言
        2.3 获取结构体字段及获取方法
        2.4 reflect.Value调用带参数的方法
        2.5 反射源码分析
        2.6 反射性能优化案例
---
## 第五章 切片和数组
    1. 值传递-数组
    2. 引用传递-指针切片和指针数组
    3. 切片和数组参数传递性能对比
    4. 切片底层结构
    5. nil切片和空切片
    6. 扩容策略
    7. 不同类型的切片间互转
    8. 带索引初始化数组和切片
---
## 第六章 指针
    1. 指针类型转换及修改值
    2. 指针分类及unsafe包使用
    3. 获取私有变量值
    4. 切片与字符串零拷贝互转
---
## 第七章 闭包  
    1. 闭包理论
    2. 匿名函数 
        2.1 匿名函数
        2.2 函数模版:定义行为
## 第八章 defer函数及汇编语言理解
    1. 注册延迟调用机制定义及使用
    2. defer陷阱
    3. 分解defer函数
    4. defer循环性能问题
## [第九章 设计模式](chapter09_design_pattern/introduction.md)
### OOP七大准则
    1. 创建型模式
        1.1 静态工厂模式-->new关键字函数实现简单工厂
        1.2 工厂方法模式-->k8s中实现
        1.3 单例模式-->标准库strings/replace实现
        1.4 原型模式-->confluent-kafka中map实现
        1.5 建造者模式-->xorm，k8s中实现
    2. 结构型模式
        2.1 组合模式
            2.1.1 修改前：使用面向对象处理
            2.1.2 修改后：使用组合模式处理
        2.2 装饰模式
            2.2.1 闭包实现
            2.2.2 grpc源码体现
            2.2.3 反射实现
        2.3 享元模式-->线程池,缓存思想
        2.4 适配器模式
        2.5 桥接模式(两个变化系统结偶)
        2.6 门面模式(外观模式)-->在gin中render应用(封装多个子服务)
        2.7 代理模式
    3. 行为型模式
        3.1  访问者模式-->k8s中实现
        3.2  迭代器-->标准库container/ring中实现
        3.3  状态模式
        3.4  责任链模式 
        3.5  模版模式
        3.6  策略模式-->if-else的另类写法(内部算法封装)
        3.7  解释器模式 
        3.8  命令模式-->go-redis中实现
        3.9  备忘录模式
        3.10 观察者模式-->etcd的watch机制
        3.11 中介者模式
    4. 函数选项:成例模式
        4.1 未使用的现状
        4.2 区分必填项和选项
        4.3 带参数的选项模式
    5. 插件式编程-->grpc中实现
    6. 同步模式(sync同步原语以及扩展原语)
        6.1 waitGroup同步等待组对象
        6.2 互斥锁（sync.Mutex）和读写锁（sync.RWMutex）性能对比
        6.3 Once单例对象
        6.4 并发安全Map(读多写少)
        6.5 Pool对象池模式( *非连接池* !）-->官方包对象池fmt
            6.5.1 未使用newFunc
            6.5.2 newFunc与GC（附源码分析）
            6.5.3 何时使用对象缓存
            6.5.4 第三方对象池object pool(bytebufferpool)
        6.6 Cond条件变量及NoCopy机制
        6.7 atomic原子操作
            6.7.0 Value的load和store
            6.7.1 add
            6.7.2 cas算法和自旋锁
            6.7.3 load和store用法
            6.7.4 swap交换
        6.8 ErrorGroup获取协程中error
        6.9 信号量Semaphore
        6.10 SingleFlight避免缓存击穿
## 第十章 函数式编程
    1. 函数介绍
    2. 高级函数
        2.1 简单实现filter,map,reduce
        2.2 简单案例
        2.3 复杂实现：泛型及类型检查
    3. 一等公民
        3.1 网络管理中问题需求
        3.2 网络管理中三种处理对比
            3.2.1 通过同享内存通信
            3.2.2 通过通信(具体数据)共享内存
            3.2.3 通过通信(函数)共享内存
## 第十一章 汇编理解go语言底层源码
[1. 汇编基本指令](chapter11_assembly_language/01asm/introduction.md)  

[2. plan9汇编](chapter11_assembly_language/02plan9/introduction.md)

    2.1  常量constant
    2.2  array类型
    2.3  bool类型
    2.4  int类型
    2.5  float类型
    2.6  slice类型
    2.7  引用类型map和channel
    2.8  函数类型
    2.9  局部变量
    2.10 流程控制
    2.11 伪SP,FP及硬件SP关系
    2.12 结构体方法
    2.13 递归函数
    2.14 闭包函数
[3. Golang底层数据结构](chapter11_assembly_language/03Golang_data_structure/data.md)                   

    3.1 Map底层结构
    3.2 String底层结构
    3.3 Struct底层结构
## 第十二章 网络编程net   
### [socket介绍](chapter12_net/socket.md)
    1. net/http高级封装演变
    2. 爬虫获取邮箱案例(http客户端源码分析)
    3. Tcp实现客户端及服务端(tcp底层原理分析)
    4. Tcp黏包分析及处理(大小端介绍)
    4. fastHttp
## [第十三章 CGO调用C语言](chapter13_Go_call_C_or_C++/introduction.md)      
    1. Go调用自定义C函数
    2. Go调用模块化C库
    3. Go实现C定义函数
    4. Go获取C函数的errno
    5. C的void返回
## [第十四章 Context上下文](chapter14_context/introduction.md)
    0. 父类EmptyCtx
    1. Context来源
        1.1 如何释放资源
        1.2 方式一：全局参数
        1.3 方式二: 通道channel
        1.4 方式三: Context
    2. WithCancel源码及使用
    3. WithDealine源码及使用
    4. WithValue源码及使用
    5. WithTimeout源码及使用
## 第十五章 接口编程
    1.1 冗余代码写法
    1.2 简单优化
    1.3 更优方式
## 第十六章 并发编程
    1. 简单流水线模型
    2. FAN-IN和FAN-OUT模型
## 第十七章 数据结构及算法
    1. queue队列
    2. 哈希函数
        2.1 hash函数分类及算法md5使用
    3. 非对称加密算法rsa
        3.1 分段与不分段加解密
    4. 排序分析
        4.1 排序算法分类及图解
        4.2 排序接口实现（反射方式）
    5. Jwt源码分析及中间件使用
## 第十八章 错误跟踪和panic
    1. 自定义错误类型打印错误栈
    2. 扩展包pkg.error
    3. Gin的错误recover分析(panic和recover源码分析)
    4. errCode自动化生成
## 第十九章 nil预定义标识
    1. 不同类型为nil时的地址和大小
    2. 不同类型与nil的比较
    3. 不同类型nil时的特点
## 第二十章 for-range源码分析
    1. 遍历数组和切片
    2. 遍历Goroutine
    3. 遍历Map