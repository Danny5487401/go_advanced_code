
# 工厂方法模式
使用子类的方式延迟生成对象到子类中实现。Go中不存在继承 所以使用匿名组合来实现


## 源码参考：k8s------>
```go
// k8s.io/kubectl/pkg/cmd/util/factory.go

type Factory interface {
	genericclioptions.RESTClientGetter
	DynamicClient() (dynamic.Interface, error)
	KubernetesClientSet() (*kubernetes.Clientset, error)
	RESTClient() (*restclient.RESTClient, error)
	NewBuilder() *resource.Builder
	ClientForMapping(mapping *meta.RESTMapping) (resource.RESTClient, error)
	UnstructuredClientForMapping(mapping *meta.RESTMapping) (resource.RESTClient, error)
	Validator(validate bool) (validation.Schema, error)
	OpenAPISchema() (openapi.Resources, error)
}

```

```go
// pkg/kubectl/cmd/cmd.go 生成工厂--->f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
func NewFactory(clientGetter genericclioptions.RESTClientGetter) Factory {
	f := &factoryImpl{
		clientGetter: clientGetter,
	}
	return f
}

```