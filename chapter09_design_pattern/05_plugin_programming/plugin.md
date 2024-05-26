<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [插件式编程](#%E6%8F%92%E4%BB%B6%E5%BC%8F%E7%BC%96%E7%A8%8B)
  - [案例](#%E6%A1%88%E4%BE%8B)
  - [第三方源码分析:grpc](#%E7%AC%AC%E4%B8%89%E6%96%B9%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90grpc)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 插件式编程
特点:可插拨

## 案例
sql驱动，grpc

## 第三方源码分析:grpc

```go
// google.golang.org/grpc@v1.26.0/resolver/resolver.go
// 定义插件的接口
type Builder interface {
    Build(target Target, cc ClientConn, opts BuildOptions) (Resolver, error)
    Scheme() string
}


// 全局变量存放插件
var (
    // m is a map from scheme to resolver builder.
    m = make(map[string]Builder)
    // defaultScheme is the default scheme to user.
    defaultScheme = "passthrough"
)
// 注册插件
func Register(b Builder) {
    m[b.Scheme()] = b
}
// 获取插件
func Get(scheme string) Builder {
    if b, ok := m[scheme]; ok {
        return b
    }
    return nil
}
```
grpc实现的插件有：consul,etcd，等等
```go
// 具体插件实现文件 internal/resolver/dns/dns_resolver.go
// 通过init函数，将实现注册到resolver
func init() { resolver.Register(NewBuilder()) }

// 实现resolver.Builder接口的 Build 函数（在这里进行真正的构建操作
func Build() {}
// 返回当前resolver解决的解析样式
func Scheme() string { return "dns" }


//应用获取 resolver clientconn.go
// 通过解析用户传入的target 获得scheme
cc.parsedTarget = grpcutil.ParseTarget(cc.target)

// 通过target的scheme获取对应的resolver.Builder
func (cc *ClientConn) getResolver(scheme string) resolver.Builder {
    for _, rb := range cc.dopts.resolvers {
        if scheme == rb.Scheme() {
            return rb
        }
    }
    return resolver.Get(scheme)
}
```



