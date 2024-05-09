<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [中介者模式（Mediator Pattern）](#%E4%B8%AD%E4%BB%8B%E8%80%85%E6%A8%A1%E5%BC%8Fmediator-pattern)
  - [应用场景](#%E5%BA%94%E7%94%A8%E5%9C%BA%E6%99%AF)
  - [中介者模式包含以下四个角色](#%E4%B8%AD%E4%BB%8B%E8%80%85%E6%A8%A1%E5%BC%8F%E5%8C%85%E5%90%AB%E4%BB%A5%E4%B8%8B%E5%9B%9B%E4%B8%AA%E8%A7%92%E8%89%B2)
  - [优点](#%E4%BC%98%E7%82%B9)
  - [缺点](#%E7%BC%BA%E7%82%B9)
  - [案例](#%E6%A1%88%E4%BE%8B)
    - [代码设计](#%E4%BB%A3%E7%A0%81%E8%AE%BE%E8%AE%A1)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 中介者模式（Mediator Pattern）
又叫作调解者模式或调停者模式。


中介者模式是用来降低多个对象和类之间的通信复杂性。这种模式提供了一个中介类，该类通常处理不同类之间的通信，并支持送耦合，使代码易于维护

## 应用场景

（1）系统中对象之间存在复杂的引用关系，产生的相互依赖关系结构混乱且难以理解。

（2）交互的公共行为，如果需要改变行为，则可以增加新的中介者类。

## 中介者模式包含以下四个角色

Mediator(抽象中介者)：它定义了一个接口，该接口用于与各同事对象之间进行通信。

ConcreteMediator(具体中介者)：它实现了接口，通过协调各个同事对象来实现协作行为，维持各个同事对象的引用

Colleague(抽象同事类)：它定义了各个同事类公有的方法，并声明了一些抽象方法来供子类实现，同时维持了一个对抽象中介者类的引用，
其子类可以通过该引用来与中介者通信。

ConcreteColleague(具体同事类)：抽象同事类的子类，每一个同事对象需要和其他对象通信时，都需要先与中介者对象通信，通过中介者来间接完成与其他同事类的通信



## 优点

（1）减少类间依赖，将多对多依赖转化成一对多，降低了类间耦合。

（2）类间各司其职，符合迪米特法则。

## 缺点

	中介者模式将原本多个对象直接的相互依赖变成了中介者和多个同事类的依赖关系。
	当同事类越多时，中介者就会越臃肿，变得复杂且难以维护
	
## 案例

某物联网企业, 研发各种智能家居产品, 并配套手机app以便用户集中控制,一开始的设计是手机app通过本地局域网的广播协议, 主动发现/注册/控制各种智能设备,
后来智能设备的种类越来越多, 通信协议多种多样, 导致手机app需要频繁升级, 集成过多驱动导致代码膨胀.

研发部门痛定思痛, 决定采用中介者模式重新设计整个系统架构

老架构: app -> 智能设备*N
新架构: app -> 云中心 -> 智能设备

通过引入"云中心" 作为中介, 将app与设备驱动解耦,app与云中心采用RESTFul协议通信,极大提升开发运维的效率

### [代码设计](mediator_pattern_test.go)

MockPhoneApp: 虚拟的手机app, 用于跟云中心通信, 控制智能设备
(抽象中介者)ICloudMediator: 云中心面向手机app的接口
ICloudCenter: 云中心面向智能设备的注册接口
(抽象同事类)ISmartDevice: 智能设备接口
(具体中介者)tMockCloudMediator: 虚拟的云中心服务类, 面向手机app实现ICloudMediator接口, 面向智能设备实现ICloudCenter接口
tMockSmartLight: 虚拟的智能灯设备, 实现ISmartDevice接口
