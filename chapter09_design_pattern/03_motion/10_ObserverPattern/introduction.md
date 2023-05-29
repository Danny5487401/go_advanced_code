<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [观察者模式](#%E8%A7%82%E5%AF%9F%E8%80%85%E6%A8%A1%E5%BC%8F)
  - [角色](#%E8%A7%92%E8%89%B2)
  - [场景](#%E5%9C%BA%E6%99%AF)
  - [设计](#%E8%AE%BE%E8%AE%A1)
  - [观察者模式的优点](#%E8%A7%82%E5%AF%9F%E8%80%85%E6%A8%A1%E5%BC%8F%E7%9A%84%E4%BC%98%E7%82%B9)
  - [观察者模式的缺点](#%E8%A7%82%E5%AF%9F%E8%80%85%E6%A8%A1%E5%BC%8F%E7%9A%84%E7%BC%BA%E7%82%B9)
  - [Go内部包](#go%E5%86%85%E9%83%A8%E5%8C%85)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 观察者模式

观察者模式（Observer Pattern）又叫作发布-订阅（Publish/Subscribe）模式、模型-视图（Model/View）模式、
源-监听器（Source/Listener）模式,或从属者（Dependent）模式 。
定义一种一对多的依赖关系，一个主题对象可被多个观察者对象同时监听，使得每当主题对象状态变化时，所有依赖它的对象都会得到通知并被自动更新，

## 角色
- 定义订阅者，并提供事件接收方法。这里的事件应该是系统事件。
- 定义发布者，这里就是系统，并提供订阅和取消订阅的方式，以及事件触发方法。
- 定义客户端，来声明订阅者和发布者，并注册订阅者

## 场景
- 某智能app, 需添加自定义闹铃的功能
- 闹铃可设定时间, 以及是否每日重复
- 可设定多个闹铃
- 根据观察者模式, 每个闹铃对象, 都是时间服务的观察者, 监听时间变化的事件.

## 设计
- ITimeService: 定义时间服务的接口, 接受观察者的注册和注销
- ITimeObserver: 定义时间观察者接口, 接收时间变化事件的通知
- tMockTimeService: 虚拟的时间服务, 自定义时间倍率以方便时钟相关的测试
- AlarmClock: 闹铃的实现类, 实现ITimeObserver接口以订阅时间变化通知

## 观察者模式的优点
1. 观察者和被观察者是松耦合（抽象耦合）的，符合依赖倒置原则。
2. 分离了表示层（观察者）和数据逻辑层（被观察者）， 并且建立了一套触发机制，使得数据的变化可以响应到多个表示层上。
3.实现了一对多的通信机制，支持事件注册机制，支持兴趣分发机制， 当被观察者触发事件时，只有感兴趣的观察者可以接收到通知。

## 观察者模式的缺点
1. 如果观察者数量过多，则事件通知会耗时较长。
2. 事件通知呈线性关系，如果其中一个观察者处理事件卡壳，则会影响后续的观察者接收该事件。
3. 如果观察者和被观察者之间存在循环依赖，则可能造成两者之间的循环调用，导致系统崩溃。


## Go内部包

在Go语言中，更适合观察者模式的场景主要是对于系统信号的处理，也就是Signal包。


