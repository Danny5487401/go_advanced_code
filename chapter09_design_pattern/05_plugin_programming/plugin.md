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
    // defaultScheme is the default scheme to use.
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



