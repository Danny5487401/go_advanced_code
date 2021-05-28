# go_advanced_code
# goVersion==1.16
![高级go编程](https://gimg2.baidu.com/image_search/src=http%3A%2F%2Ft.ki4.cn%2F2020%2F1%2FvIVv6v.jpg&refer=http%3A%2F%2Ft.ki4.cn&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=jpeg?sec=1621305693&t=1a817e6e6ecf0e1ec1890212636f0c19)
# *目录*
## 第一章 文件操作
    1. os中FileInfo底层的文件描述符和相关信息
    2. os文件操作
    3. io包底层Reader和Writer接口
    4. io断点续传
    5. bufio缓存读写
---
## 第二章 协程
    1. 两级线程模型
    2. runtime模块和GC
        2.1 runtime核心功能及系统信息调用
        2.2 Goexit()终止线程
        2.3 资源竞争一致性问题分析
        2.4 垃圾回收机制
    3. sync模块
    4. 并发数量控制
    5. goroutine泄漏分析及处理
    6. 线程池
---
## 第三章 通道
    1. 基本使用
        1.1 阻塞
        1.2 父子通信
        1.3 死锁
        1.4 close(ch)关闭所有下游协程
        1.5 通道遍历range
        1.6 缓冲channel增强并发
        1.7 双向通道
        1.8 单向通道
        1.9 使用channel传递channel
    2. TimerChan模块
    3. Select多路复用
    4. CSP理论中的Process/Channel
    5. Channel通道结构
---
## 第四章 反射
    1. 反射原理及应用场景
    2. 反射使用：已知与未知类型
---
## 第五章 切片和数组
    1. 值传递
    2. 引用传递
    3. 切片和数组参数传递性能对比
    4. 切片底层结构
    5. nil切片和空切片
    6. 扩容策略
---
## 第六章 指针
    1. 指针分类
    2. unsafe包使用
    3. 获取私有变量值
---
## 第七章 闭包  
    1. 闭包理论
    2. 闭包应用   
## 第八章 defer函数及汇编语言理解
    1. 注册延迟调用机制定义及使用
    2. defer陷阱
    3. defer传参和返回设置
    4. defer循环性能问题
## 第九章 设计模式
    1. 简单工厂模式
    2. 装饰模式
    3. 单例模式
    4. 责任链模式 
    5. 策略模式
    6. 解释器模式
## 第十章 编写高级函数map,reduce
    1. 简单写法
    2. 泛型及类型检查
## 第十一章 汇编理解go语言
    1. 汇编的基本概念
    2. plan9汇编
    3. Golang底层数据结构
## 第十二章 网络编程net
    1. net/http高级封装
    2. 爬虫案例
    3. Tcp使用
    4. Tcp黏包分析及处理
## 第十三章 CGO调用C语言
    1. Go调用C函数
    2. Go调用C库
## 第十四章 Context上下文
    1. Context来源
        1.1 如何释放资源
        1.2 方式一：全局参数
        1.3 方式二: 通道channel
        1.4 方式三: Context
    2. WithCancel
    3. WithDealine
    4. WithValue
    5. WithTimeout
