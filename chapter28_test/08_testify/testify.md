<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [testify](#testify)
  - [第三方使用](#%E7%AC%AC%E4%B8%89%E6%96%B9%E4%BD%BF%E7%94%A8)
  - [1. assert](#1-assert)
    - [Assertions 对象](#assertions-%E5%AF%B9%E8%B1%A1)
  - [2. mock](#2-mock)
  - [3. suite](#3-suite)
  - [参考资料](#%E5%8F%82%E8%80%83%E8%B5%84%E6%96%99)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# testify
三大套件：assert断言，mock测试替身，suite测试套件

## 第三方使用

gin框架使用

## 1. assert
assert子库提供了便捷的断言函数，可以大大简化测试代码的编写。总的来说，它将之前需要判断 + 信息输出的模式：
```go
if got != expected {
  t.Errorf("Xxx failed expect:%d got:%d", got, expected)
}
```
简化后
```go
assert.Equal(t, got, expected, "they should be equal")
```

### Assertions 对象

观察到上面的断言都是以TestingT为第一个参数，需要大量使用时比较麻烦。testify提供了一种方便的方式。
先以*testing.T创建一个*Assertions对象，Assertions定义了前面所有的断言方法，只是不需要再传入TestingT参数了
```go
func TestEqual(t *testing.T) {
  assertions := assert.New(t)
  assertion.Equal(a, b, "")
  // ...
}
```




## 2. mock
Mock 简单来说就是构造一个仿对象，仿对象提供和原对象一样的接口，在测试中用仿对象来替换原对象。这样我们可以在原对象很难构造，特别是涉及外部资源（数据库，访问网络等）。
    
例如，我们现在要编写一个从一个站点拉取用户列表信息的程序，拉取完成之后程序显示和分析。
如果每次都去访问网络会带来极大的不确定性，甚至每次返回不同的列表，这就给测试带来了极大的困难。我们可以使用 Mock 技术。

## 3. suite

testify提供了测试套件的功能（TestSuite），testify测试套件只是一个结构体，内嵌一个匿名的suite.Suite结构。
测试套件中可以包含多个测试，它们可以共享状态，还可以定义钩子方法执行初始化和清理操作。钩子都是通过接口来定义的，实现了这些接口的测试套件结构在运行到指定节点时会调用对应的方法。

如果定义了SetupSuite()方法（即实现了SetupAllSuite接口），在套件中所有测试开始运行前调用这个方法
```go
type SetupAllSuite interface {
  SetupSuite()
}
```
如果定义了TearDonwSuite()方法（即实现了TearDownSuite接口），在套件中所有测试运行完成后调用这个方法
```go
type TearDownAllSuite interface {
  TearDownSuite()
}
```


## 参考资料
1. [每日一库 testify](https://darjun.github.io/2021/08/11/godailylib/testify/)
2. [testify 官网](https://github.com/stretchr/testify)