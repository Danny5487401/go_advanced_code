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
```go
//（对于map和slice无效，依旧共享相同内存对象，其他会拷贝一份）
//（可以单独处理map和slice）
Pc2:=Pc1
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




