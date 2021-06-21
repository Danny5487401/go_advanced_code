package main

import "fmt"

/*
插件式编程
	特点可插拨
源码分析
	google.golang.org/grpc@v1.26.0/resolver/resolver.go

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

	具体插件实现文件 internal/resolver/dns/dns_resolver.go
	// 通过init函数，将实现注册到resolver
	func init() { resolver.Register(NewBuilder()) }

	// 实现resolver.Builder接口的 Build 函数（在这里进行真正的构建操作
	func Build() {}
	// 返回当前resolver解决的解析样式
	func Scheme() string { return "dns" }


	应用获取 resolver clientconn.go
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
*/
// 定义一个接口，里面有两个方法

type pluginFunc interface {
	hello()
	world()
}

// 定义一个类，来存放我们的插件
type plugins struct {
	plist map[string]pluginFunc
}

// 初始化插件
func (p *plugins) init() {
	p.plist = make(map[string]pluginFunc)
}

// 注册插件
func (p *plugins) register(name string, plugin pluginunc) {
	p.plist[name] = plugin
	//p.plist = append(p.plist, a)
}

// 我们定义几个插件
//plugin1
type plugin1 struct{}

func (p *plugin1) hello() {
	fmt.Println("plugin1 hello")
}
func (p *plugin1) world() {
	fmt.Println("plugin1 world")
}

//plugin2
type plugin2 struct{}

func (p *plugin2) hello() {
	fmt.Println("plugin2 hello")
}
func (p *plugin2) world() {
	fmt.Println("plugin2 world")
}

// 开始调用
func main() {
	plugin := new(plugins)
	plugin.init()

	plugin1 := new(plugin1)
	plugin2 := new(plugin2)

	plugin.register("plugin1", plugin1)
	plugin.register("plugin2", plugin2)
	for _, plugin := range plugin.plist {
		plugin.hello()
		plugin.world()
	}
}
