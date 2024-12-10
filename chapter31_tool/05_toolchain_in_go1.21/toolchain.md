<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [toolchain 规则](#toolchain-%E8%A7%84%E5%88%99)
  - [Go 1.21版本之前的向前兼容性问题](#go-121%E7%89%88%E6%9C%AC%E4%B9%8B%E5%89%8D%E7%9A%84%E5%90%91%E5%89%8D%E5%85%BC%E5%AE%B9%E6%80%A7%E9%97%AE%E9%A2%98)
  - [Go 1.21版本后的向前兼容性策略](#go-121%E7%89%88%E6%9C%AC%E5%90%8E%E7%9A%84%E5%90%91%E5%89%8D%E5%85%BC%E5%AE%B9%E6%80%A7%E7%AD%96%E7%95%A5)
  - [GOTOOLCHAIN环境变量与toolchain版本选择](#gotoolchain%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F%E4%B8%8Etoolchain%E7%89%88%E6%9C%AC%E9%80%89%E6%8B%A9)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# toolchain 规则

Go 1.21通过增强go语句语义和添加工具链管理，大幅改进了Go语言的向前兼容性。

## Go 1.21版本之前的向前兼容性问题

```go
// go.mod
module demo1

go 1.20

```

上面go.mod文件中的go directive表示建议使用Go 1.20及以上版本编译本module代码，但并不强制禁止使用低于1.20版本的Go对module进行编译。你也可以使用Go 1.19版本，甚至是Go 1.15版本编译这个module的代码


## Go 1.21版本后的向前兼容性策略

```go
// go.mod

module demo1

go 1.21.1 // 指定最小可用版本为Go 1.21.1
```


## GOTOOLCHAIN环境变量与toolchain版本选择


是否执行自动工具链下载和缓存、Go toolchain switches(Go工具链切换)以及切换的工具链的版本取决于GOTOOLCHAIN环境变量的设置、go.mod中go和toolchain指示的版本。

当go命令捆绑的工具链与module的go.mod的go或工具链版本一样时或更新时，go命令会使用自己的捆绑工具链。例如，当在main module的go.mod包含有go 1.21.0时，如果go命令绑定的工具链是Go 1.21.3版本，那么将继续使用初始toolchain的版本，即Go 1.21.3


而如果go.mod中的go版本写着go 1.21.9，那么go命令捆绑的工具链版本1.21.3显然不能满足要求，那此时就要看GOTOOLCHAIN环境变量的配置。


```shell
$cat go.mod
module demo1

go 1.23.1
toolchain go1.23.1 

$GOTOOLCHAIN=go1.24.1+auto go build
go: downloading go1.24.1 (darwin/amd64) // 使用name指定工具链，但该工具链本地不存在，于是下载。

$GOTOOLCHAIN=go1.20.1+auto go build
go: downloading go1.23.1 (darwin/amd64) // 使用go.mod中的版本的工具链
```

当GOTOOLCHAIN设置为\<name>+auto时，go命令会根据需要选择并运行较新的Go版本。具体来说，它会查询go.mod文件中的工具链版本和go version。
如果go.mod 文件中有toolchain行，且toolchain指示的版本比默认的Go工具链(name)新，那么系统就会调用toolchain指示的工具链版本；反之会使用默认工具链。


## 参考

- [toolchain的使用规则](https://tonybai.com/2023/09/10/understand-go-forward-compatibility-and-toolchain-rule/)
- https://go.dev/doc/toolchain