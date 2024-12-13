<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Go 1.16 新特性](#go-116-%E6%96%B0%E7%89%B9%E6%80%A7)
  - [module deprecation 后悔药](#module-deprecation-%E5%90%8E%E6%82%94%E8%8D%AF)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Go 1.16 新特性

## module deprecation 后悔药
在使用场景上：在发现严重问题或无意发布某些版本后，模块的维护作者可以撤回该版本，支持撤回单个或多个版本.
以前没有办法解决，因此一旦出现就非常麻烦。对应两者的操作如下：

维护者：

- 删除有问题版本的 tag。
- 重新打一个新版本的 tag。
使用者：

- 发现有问题的版本 tag 丢失，需手动介入。
- 不知道有问题，由于其他库依赖，因此被动升级而踩坑。


在Go 1.16版本在go.mod中加入retract，以帮助go module作者作废自己的module。

```go
module github.com/dgraph-io/dgraph/v24

go 1.22.7

// ... 

retract v24.0.3 // should have been a minor release instead of a patch
```

对于那些使用了被废弃的module的go项目，go list、go get命令都会给出warning。



## 参考

- [Go1.16 新特性：Go mod 的后悔药](https://segmentfault.com/a/1190000039359906)