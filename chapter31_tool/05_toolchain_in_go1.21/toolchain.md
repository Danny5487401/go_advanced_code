<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [toolchain 规则](#toolchain-%E8%A7%84%E5%88%99)
  - [Go 1.21版本之前的向前兼容性问题](#go-121%E7%89%88%E6%9C%AC%E4%B9%8B%E5%89%8D%E7%9A%84%E5%90%91%E5%89%8D%E5%85%BC%E5%AE%B9%E6%80%A7%E9%97%AE%E9%A2%98)
  - [Go 1.21版本后的向前兼容性策略](#go-121%E7%89%88%E6%9C%AC%E5%90%8E%E7%9A%84%E5%90%91%E5%89%8D%E5%85%BC%E5%AE%B9%E6%80%A7%E7%AD%96%E7%95%A5)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# toolchain 规则


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