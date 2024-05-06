<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [构建约束（build constraint）](#%E6%9E%84%E5%BB%BA%E7%BA%A6%E6%9D%9Fbuild-constraint)
  - [golang支持的两种构建约束](#golang%E6%94%AF%E6%8C%81%E7%9A%84%E4%B8%A4%E7%A7%8D%E6%9E%84%E5%BB%BA%E7%BA%A6%E6%9D%9F)
    - [1 使用文件后缀进行约束](#1-%E4%BD%BF%E7%94%A8%E6%96%87%E4%BB%B6%E5%90%8E%E7%BC%80%E8%BF%9B%E8%A1%8C%E7%BA%A6%E6%9D%9F)
    - [2 通过代码文件中添加注释进行约束](#2-%E9%80%9A%E8%BF%87%E4%BB%A3%E7%A0%81%E6%96%87%E4%BB%B6%E4%B8%AD%E6%B7%BB%E5%8A%A0%E6%B3%A8%E9%87%8A%E8%BF%9B%E8%A1%8C%E7%BA%A6%E6%9D%9F)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 构建约束（build constraint）
构建约束（build constraint），也叫做构建标记（build tag），是在 Go 源文件最开始的注释行，比如

```go
// +build linux
```

构建约束也称之为条件编译，就是可以对某些源代码文件指定在特定的平台，架构，编译器甚至Go版本下面进行编译，在其他环境中会自动忽略这些文件

几个注意点：

- 约束可以出现在任何源文件中，比如 .go、.s 等；
- 必须在文件顶部附近，它的前面只能有空行或其他注释行；可见包子句也在约束之后；
- 约束可以有多行；
- 为了区别约束和包文档，在约束之后必须有空行


## golang支持的两种构建约束


### 1 使用文件后缀进行约束

这种方式就是通过文件的后缀名来对要指定平台的编译的文件进行约束，文件格式如下：

```shell
sourcefilename_GOOS_GOARCH.go


```

GOOS和GOARCH可以通过go env看到

- user_windows_amd64.go  // 在 windows 中 amd64 架构下才会编译，其他的环境中会自动忽略
- user_linux_arm.go     // 在 linux 中的 arm 架构下才会编译，其他环境中会自动忽略


### 2 通过代码文件中添加注释进行约束


旧版构建约束: 构建约束的语法是 // +build 这种形式，如果多个条件组合，通过空格、逗号或多行构建约束表示
```go
// +build linux,amd64

package user

func User(){}
```

注释编译首先就是在文件头部添加注释 // +build 然后后面添加编译条件，该注释要和下面的代码空一行，否则就会被作为普通的注释对待；

当然，也可以Go1.17 后新版的构建约束，也使用了 //go: 开头：
```go
//go:build linux,amd64

package user

func User(){}
```

注释支持五种已有的标签，有：

- 操作系统，环境变量中GOOS的值，如：linux、darwin、windows等等
- 操作系统的架构，环境变量中GOARCH的值，如：arch64、x86、i386等等
- 使用的编译器，如：gc或者gccgo，是否开启CGO,cgo
- golang版本号，如：Go Version 1.5, Go Version 1.13，以此类推
- 自定义标签，这种标签就需要在编译时指定，如：通过 go build -tags 自定义tag名称 指定tag值
- // +build ignore，编译时自动忽略该文件

标签之间有如下几种运算关系：
* 空格表示：OR
* 逗号表示：AND
* !表示：NOT
* 换行表示：AND
```go
// +build linux,386 darwin,!cgo 

// 运算表达式为：(linux && 386) || (darwin && !cgo)

```

有些时候，多个约束分成多行书写，会更易读些：
```shell
// +build linux darwin
// +build amd64

// 这相当于 (linux OR darwin) AND amd64 
```