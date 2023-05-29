<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [单例模式](#%E5%8D%95%E4%BE%8B%E6%A8%A1%E5%BC%8F)
  - [意图](#%E6%84%8F%E5%9B%BE)
  - [主要解决](#%E4%B8%BB%E8%A6%81%E8%A7%A3%E5%86%B3)
  - [何时使用](#%E4%BD%95%E6%97%B6%E4%BD%BF%E7%94%A8)
  - [如何解决](#%E5%A6%82%E4%BD%95%E8%A7%A3%E5%86%B3)
    - [关键代码：构造函数是私有的。](#%E5%85%B3%E9%94%AE%E4%BB%A3%E7%A0%81%E6%9E%84%E9%80%A0%E5%87%BD%E6%95%B0%E6%98%AF%E7%A7%81%E6%9C%89%E7%9A%84)
  - [使用场景](#%E4%BD%BF%E7%94%A8%E5%9C%BA%E6%99%AF)
  - [优点](#%E4%BC%98%E7%82%B9)
  - [缺点：](#%E7%BC%BA%E7%82%B9)
  - [分类：](#%E5%88%86%E7%B1%BB)
  - [标准库实现](#%E6%A0%87%E5%87%86%E5%BA%93%E5%AE%9E%E7%8E%B0)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 单例模式
在程序的运行中只产生一个实例。

## 意图
保证一个类仅有一个实例，并提供一个访问它的全局访问点。

## 主要解决
一个全局使用的类频繁地创建与销毁。

## 何时使用
当您想控制实例数目，节省系统资源的时候。

## 如何解决
![](process.png)
判断系统是否已经有这个单例，如果有则返回，如果没有则创建。

### 关键代码：构造函数是私有的。

## 使用场景
1. 要求生产唯一序列号。
2. WEB 中的计数器，不用每次刷新都在数据库里加一次，用单例先缓存起来。
3. 创建的一个对象需要消耗的资源过多，比如 I/O 与数据库的连接等

## 优点
在内存里只有一个实例，减少了内存的开销，尤其是频繁的创建和销毁实例（比如管理学院首页页面缓存）。
避免对资源的多重占用（比如写文件操作）。

## 缺点：
没有接口，不能继承，与单一职责原则冲突，一个类应该只关心内部逻辑，而不关心外面怎么样来实例化。

## 分类：
有懒汉式和饿汉式


## 标准库实现
结构体定义
```go
//strings/replace.go
type Replacer struct {
    once   sync.Once // 控制 r replacer 替换算法初始化
    r      replacer
    oldnew []string
}
```
```go

//线程安全且支持规则复用
func NewReplacer(oldnew ...string) *Replacer {
        //****
    return &Replacer{oldnew: append([]string(nil), oldnew...)} //没有创建算法
}
//当我们使用 strings.NewReplacer 创建 strings.Replacer 时，这里采用惰性算法，并没有在这时进行 build 解析替换规则并创建对应算法实例，
func (r *Replacer) Replace(s string) string {
    r.once.Do(r.buildOnce) //初始化
    return r.r.Replace(s)
}

func (r *Replacer) WriteString(w io.Writer, s string) (n int, err error) {
    r.once.Do(r.buildOnce) //初始化
    return r.r.WriteString(w, s)
}
//而是在执行替换时( Replacer.Replace 和 Replacer.WriteString)进行的
//初始化算法
func (r *Replacer) buildOnce() {
    r.r = r.build()
    r.oldnew = nil
}

func (b *Replacer) build() replacer {
    //....
}

```


