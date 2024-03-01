<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [go开发套件](#go%E5%BC%80%E5%8F%91%E5%A5%97%E4%BB%B6)
  - [指令](#%E6%8C%87%E4%BB%A4)
    - [go get 和 go install](#go-get-%E5%92%8C-go-install)
    - [环境变量](#%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F)
  - [参考资料](#%E5%8F%82%E8%80%83%E8%B5%84%E6%96%99)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# go开发套件

go开发套件包含一组编译和构建源码的指令和对应的程序。

通常有三种方式来运行这些程序：

1. 最普遍的，通过go subcommand的形式，如go fmt。

go在调用这些子命令对应的程序时， 会传递用于处理package层次的参数， 来使命令运行在所有的package上。

2. 另一种形式是go tool subcommand，被称为 独立运行 (stand-alone) 。

对于大多数指令来说，这种形式仅用来调试； 对于pprof等某些指令，只能用这种形式来运行。

3. 最后，因为gofmt和godoc经常被使用，单独构建了二进制包

## 指令


| col1  | col2                                                                                                                                                         |
|-------| ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| go    | 管理go源码和运行其他的指令                                                                                                                                   |
| asm   | Asm, typically invoked as “go tool asm”, assembles the source file into an object file named for the basename of the argument source file with a .o suffix. |
| fmt   | 格式化源码                                                                                                                                                   |
| godoc | 导出并生成go代码中的文档                                                                                                                                     |
| fix   | 用于将使用了语言或lib的旧特性的程序，改写成新特性                                                                                                            |
| list  | 提供指定代码包的更深层次的信息                                                                                                       |
| clean | 删除掉执行其它命令时产生的一些文件和目录                                                                                                       |



### go get 和 go install

go get与go install的职责分别是“管理依赖”和“安装模块”。

```shell

go install [build flags] [packages]
```
go install也会将可执行文件安装到GOBIN目录下。


### go clean 删除掉执行其它命令时产生的一些文件和目录



### 环境变量

查看环境变量NAME go env <NAME>

修改环境变量NAME go env -w <NAME>=<VALUE>

```shell
✗ go env                         
GO111MODULE="on"
GOARCH="arm64"
GOBIN=""
GOCACHE="/Users/python/Library/Caches/go-build"
GOENV="/Users/python/Library/Application Support/go/env"
GOEXE=""
GOEXPERIMENT=""
GOFLAGS=""
GOHOSTARCH="arm64"
GOHOSTOS="darwin"
GOINSECURE=""
GOMODCACHE="/Users/python/go/pkg/mod"
GONOPROXY="gitlab.xxx.com"
GONOSUMDB="gitlab.xxx.com"
GOOS="darwin"
GOPATH="/Users/python/go"
GOPRIVATE="gitlab.xxx.com"
GOPROXY="https://goproxy.cn,direct"
GOROOT="/Users/python/go/go1.18"
GOSUMDB="sum.golang.org"
GOTMPDIR=""
GOTOOLDIR="/Users/python/go/go1.18/pkg/tool/darwin_arm64"
GOVCS=""
GOVERSION="go1.18"
GCCGO="gccgo"
AR="ar"
CC="clang"
CXX="clang++"
CGO_ENABLED="1"
GOMOD="/Users/python/Desktop/go_advanced_code/go.mod"
GOWORK=""
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
PKG_CONFIG="pkg-config"
GOGCCFLAGS="-fPIC -arch arm64 -pthread -fno-caret-diagnostics -Qunused-arguments -fmessage-length=0 -fdebug-prefix-map=/var/folders/sk/m49vysmj3ss_y50cv9cvcn800000gn/T/go-build1225457076=/tmp/go-build -gno-record-gcc-switches -fno-common"

```
通用的

- GO111MODULE 控制是否运行在module-aware模式下。
- GCCGO go build -compiler=gccgo使用的指令。
- GOARCH 目标指令集。
- GOBIN go install的目标路径。
- GOCACHE 构建缓存存放的文件夹。
- GOMODCACHE 下载的模块的缓存的文件夹。
- GODEBUG 激活多种调试机制。？
- GOENV 存储go环境变量的文件。
- GOFLAGS 空格分割的-flag=value列表， 当要执行的go指令支持这些flag时，会传递给go指令。 优先级低于直接在命令中给出的flag。
- GOINSECURE 一组逗号分割的模块通配符， 符合的模块会被使用不安全的方法来获取。 只在直接获取的模块上生效。
- GOOS 编译目标的操作系统。
- GOPATH 指明需要从哪些地方来获取go代码， 在使用模块时不再用于解析引入的包。
- unix下，是冒号分割的字符串。
- windows下，是分号分割的字符串。
- GOPROXY go模块代理的地址。
- GOPRIVAGE,GONOPROXY,GONOSUMDB 逗号分割的模块前缀通配模式， 符合模式的模块会直接获取 (be fetched directly) ， 也不会进行校验和校验。
- GOROOT go树的根
- GOSUMDB 需要使用的校验和数据库。
- GOTMPDIR go指令使用的临时文件夹。
- GOVCS 会用来尝试匹配服务器的版本控制指令

## 参考资料

1. [go cmd](https://pkg.go.dev/cmd)
