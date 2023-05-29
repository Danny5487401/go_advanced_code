<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [io.FS 的抽象](#iofs-%E7%9A%84%E6%8A%BD%E8%B1%A1)
  - [Go 1.16 关于 io 有哪些改变](#go-116-%E5%85%B3%E4%BA%8E-io-%E6%9C%89%E5%93%AA%E4%BA%9B%E6%94%B9%E5%8F%98)
  - [FS 接口](#fs-%E6%8E%A5%E5%8F%A3)
  - [参考链接](#%E5%8F%82%E8%80%83%E9%93%BE%E6%8E%A5)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# io.FS 的抽象

Go 在文件 IO 的场景有个神奇的事情。打开一个文件的时候，返回的竟然不是 interface ，而是一个 os.File  结构体的指针
```go
func Open(name string) (*File, error) {
    return OpenFile(name, O_RDONLY, 0)
}
```
这个意味着，Go 的文件系统的概念和 OS 的文件系统的概念直接关联起来。你必须传入一个文件路径，并且必须真的要去打开一个操作系统的文件.
![](chapter01_input_output/.io_images/io_issue.png)
![](chapter01_input_output/.io_imagesio_issue_fix.png)

## Go 1.16 关于 io 有哪些改变

- 新增了一个 io/fs 的包，抽象了一个 FS 出来。
- embed 的 package 用了这个抽象。
- 规整 io/ioutil 里面的内容。

把之前大杂烩的 io/ioutil 里面的东西拆出来了。移到对应的 io 包和 os 包。为了兼容性，ioutil 包并没有直接删除，而是导入。比如：

- Discard 移到了 io 库实现
- ReadAll 移到了 io 库实现
- NopCloser 移到了 io 库实现
- ReadFile 移到 os 库实现
- WriteFile 移到 os 库实现
## FS 接口

```go
// 文件系统的接口
type FS interface {
    Open(name string) (File, error)
}

// 文件的接口
type File interface {
    Stat() (FileInfo, error)
    Read([]byte) (int, error)
    Close() error
}
```


Go 理解的文件系统，只要能实现一个 Open 方法，返回一个 File 的 interface ，这个 File 只需要实现 Stat，Read，Close 方法即可。

我觉得主要有两方面的原因：

- 屏蔽了fs具体的实现，我们可以按照需要替换成不同的文件系统实现，比如内存/磁盘/网络/分布式的文件系统
- 易于单元测试，实际也是上一点的体现

Go 在此 io.FS 的基础上，再去扩展接口，增加文件系统的功能。

1. 比如，加个 ReadDir 就是一个有读目录的文件系统 ReadDirFS ：

```go
type ReadDirFS interface {
    FS
    // 读目录
    ReadDir(name string) ([]DirEntry, error)
}
```

2. 加个 Glob 方法，就成为一个具备路径通配符查询的文件系统

```go
type GlobFS interface {
    FS
    // 路径通配符的功能
    Glob(pattern string) ([]string, error)
}
```

3. 加个 Stat ，就变成一个路径查询的文件系统：
```go
type StatFS interface {
    FS
    // 查询某个路径的文件信息
    Stat(name string) (FileInfo, error)

```


## 参考链接
1. [go 眼中的 Fs](https://mp.weixin.qq.com/s/bZO6kfhfdMaOkYZuGjOl_Q)