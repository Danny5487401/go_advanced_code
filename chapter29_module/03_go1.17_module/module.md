<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Go 1.17 新特性: module依赖图修剪与延迟module加载](#go-117-%E6%96%B0%E7%89%B9%E6%80%A7-module%E4%BE%9D%E8%B5%96%E5%9B%BE%E4%BF%AE%E5%89%AA%E4%B8%8E%E5%BB%B6%E8%BF%9Fmodule%E5%8A%A0%E8%BD%BD)
  - [module依赖图修剪 module graph pruning](#module%E4%BE%9D%E8%B5%96%E5%9B%BE%E4%BF%AE%E5%89%AA-module-graph-pruning)
  - [延迟module加载(lazy module loading)](#%E5%BB%B6%E8%BF%9Fmodule%E5%8A%A0%E8%BD%BDlazy-module-loading)
  - [module deprecation注释](#module-deprecation%E6%B3%A8%E9%87%8A)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Go 1.17 新特性: module依赖图修剪与延迟module加载

## module依赖图修剪 module graph pruning

![](.module_images/module_graph.png)

main module中的lazy.go导入了module a的package x，后者则导入了module b；并且module a还有一个package y，该包导入了module c

现在问题来了！package y是因为自身是module a的一部分而被main module所依赖，它没有为main module的构建做出任何“代码级贡献”；
同理，package y所依赖的module c亦是如此。但是在Go 1.17之前的版本中，如果Go编译器找不到module c，那么main module的构建将会失败，这会让开发者们觉得不够合理！

gnet为例，Go 1.17版本之前的go.mod如下：

```go

module github.com/panjf2000/gnet

go 1.16

require (
 github.com/BurntSushi/toml v0.3.1 // indirect
 github.com/panjf2000/ants/v2 v2.4.6
 github.com/stretchr/testify v1.7.0
 github.com/valyala/bytebufferpool v1.0.0
 go.uber.org/atomic v1.8.0 // indirect
 go.uber.org/multierr v1.7.0 // indirect
 go.uber.org/zap v1.18.1
 golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c
 gopkg.in/natefinch/lumberjack.v2 v2.0.0
)
```

Go 1.17重新mod tidy后，go.mod内容如下：
```go

module github.com/panjf2000/gnet

go 1.17

require (
 github.com/BurntSushi/toml v0.3.1 // indirect
 github.com/panjf2000/ants/v2 v2.4.6
 github.com/stretchr/testify v1.7.0
 github.com/valyala/bytebufferpool v1.0.0
 go.uber.org/atomic v1.8.0 // indirect
 go.uber.org/multierr v1.7.0 // indirect
 go.uber.org/zap v1.18.1
 golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c
 gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

require (
 github.com/davecgh/go-spew v1.1.1 // indirect
 github.com/pmezard/go-difflib v1.0.0 // indirect
 gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
```

我们看到go 1.17后，go.mod中的main module的依赖分成了两个require块儿，第一个是直接依赖，第二个是间接依赖。


这种将那些“占着茅坑不拉屎”、对构建完全没有“贡献”的间接依赖module从构建时使用的依赖图中修剪掉的过程，就被称为module依赖图修剪。



## 延迟module加载(lazy module loading)


延迟module加载其含义就是那些在完整的module graph(complete module graph)中，但不在pruned module graph中的module的go.mod不会被go命令加载。



## module deprecation注释

在Go 1.16版本在go.mod中加入retract，以帮助go module作者作废自己的module。
go module作者只需在自己的go.mod中的module声明上面用**// Deprecated: comment**对module做出注释，就像下面这样：

```go
// Deprecated: user example.com/mod/v2 instead.
module example.com/mod
```

对于那些使用了被废弃的module的go项目，go list、go get命令都会给出warning。