<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [池化技术](#%E6%B1%A0%E5%8C%96%E6%8A%80%E6%9C%AF)
  - [协程池：](#%E5%8D%8F%E7%A8%8B%E6%B1%A0)
  - [案例：](#%E6%A1%88%E4%BE%8B)
  - [分类：](#%E5%88%86%E7%B1%BB)
  - [优点](#%E4%BC%98%E7%82%B9)
  - [缺点](#%E7%BC%BA%E7%82%B9)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 池化技术
    
池化技术 (Pool) 是一种很常见的编程技巧，在请求量大时能明显优化应用性能，降低系统频繁建连的资源开销。
我们日常工作中常见的有数据库连接池、线程池、对象池等，它们的特点都是将 “昂贵的”、“费时的” 的资源维护在一个特定的 “池子” 中，
规定其最小连接数、最大连接数、阻塞队列等配置，方便进行统一管理和复用，通常还会附带一些探活机制、强制回收、监控一类的配套功能

## 协程池：
能够达到协程资源复用。

## 案例：
有基于链表实现的Tidb，有基于环形队列实现的Jaeger，有基于数组栈实现的FastHTTP等
## 分类：
1. 提前创建协程：Jaeger，Istio，Tars等。
2. 按需创建协程：Tidb，FastHTTP，Ants等。

## 优点
- 减少内存碎片的产生
- 提高内存的使用频率
## 缺点
造成内存的浪费