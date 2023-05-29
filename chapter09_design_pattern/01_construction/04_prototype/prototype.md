<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [原型模式](#%E5%8E%9F%E5%9E%8B%E6%A8%A1%E5%BC%8F)
  - [解决的问题：](#%E8%A7%A3%E5%86%B3%E7%9A%84%E9%97%AE%E9%A2%98)
  - [优点](#%E4%BC%98%E7%82%B9)
  - [缺点](#%E7%BC%BA%E7%82%B9)
  - [浅拷贝](#%E6%B5%85%E6%8B%B7%E8%B4%9D)
  - [深拷贝](#%E6%B7%B1%E6%8B%B7%E8%B4%9D)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 原型模式

用原型实例指定创建对象的种类，并且通过拷贝这些原型创建新的对象。通过实现克隆clone()操作，快速的生成和原型对象一样的实例

## 解决的问题：
1. 资源优化，类初始化要消耗非常多的资源，包括硬件，数据等资源。
2. 性能和安全要求的场景
3. 通过new产生一个对象需要非常繁琐的数据准备或访问权限，使用原型模式克隆比直接new一个对象再逐属性赋值的过程更简洁高效。。
4. 一个对象多个修改者的场景
5. 实际项目中，原型模式大多和工厂模式一起出现，通过clone的方法返回一个对象，然后由工厂模式提供给调用者

## 优点
1. 性能提高
2. 逃避构造函数的约束
## 缺点
1. clone方法位于类的内部，当对已有类进行改造的时候，需要修改代码，违背了开闭原则
2. 当实现深克隆时，需要编写较为复杂的代码，尤其当对象之间存在多重嵌套引用时，为了实现深克隆，每一层对象对应的类都必须支持深克隆。因此，深克隆、浅克隆需要运用得当。

## 浅拷贝

如果你需要拷贝的对象中没有引用类型，那么对于Golang而言使用浅拷贝就可以了。

```go
//（对于map和slice无效，依旧共享相同内存对象，其他会拷贝一份）
//（可以单独处理map和slice）
Pc2 := Pc1
```

## 深拷贝
1. 简单值
```go
func (e *Example) Clone() *Example {
    res := *e
    return &res
}
```
2. map拷贝
```go
// map操作-->github.com/confluentinc/confluent-kafka-go@v1.7.0/kafka/config.go
type ConfigValue interface{}
type ConfigMap map[string]ConfigValue
func (m ConfigMap) clone() ConfigMap {
    m2 := make(ConfigMap)
    for k, v := range m {
    m2[k] = v
    }
    return m2
}
```


## 参考
1. [Golang深拷贝浅拷贝](https://blog.csdn.net/weixin_40165163/article/details/90680466)




