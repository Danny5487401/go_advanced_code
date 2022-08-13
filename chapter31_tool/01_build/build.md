# go build 

go 的编译是以 package main 的 main() 函数作为主入口，生成可执行文件。若 build 的是非 main 包，则不会生成可执行文件，只检查是否可执行编译。
可以输入 go help build 查看官方解释

go build 编译包时，会忽略“_test.go”结尾的文件（即测试文件）。
```go
➜  go_advanced_code git:(feature/map) ✗ go help build
usage: go build [-o output] [build flags] [packages]

Build compiles the packages named by the import paths,
along with their dependencies, but it does not install the results.

If the arguments to build are a list of .go files from a single directory,
build treats them as a list of source files specifying a single package.

When compiling packages, build ignores files that end in '_test.go'.

When compiling a single main package, build writes
the resulting executable to an output file named after
the first source file ('go build ed.go rx.go' writes 'ed' or 'ed.exe')
or the source code directory ('go build unix/sam' writes 'sam' or 'sam.exe').
The '.exe' suffix is added when writing a Windows executable.

When compiling multiple packages or a single non-main package,
build compiles the packages but discards the resulting object,
serving only as a check that the packages can be built.

The -o flag forces build to write the resulting executable or object
to the named output file or directory, instead of the default behavior described
in the last two paragraphs. If the named output is an existing directory or
ends with a slash or backslash, then any resulting executables
will be written to that directory.

The -i flag installs the packages that are dependencies of the target.
The -i flag is deprecated. Compiled packages are cached automatically.

The build flags are shared by the build, clean, get, install, list, run,
and test commands:
```
![](01_tags/.build_images/build_commnd.png)


## 参考资料
1. [Go 增量构建](https://mp.weixin.qq.com/s?__biz=MzIyNzM0MDk0Mg==&mid=2247491831&idx=1&sn=8eb54239e5105aed870ae931b338868e&chksm=e8600716df178e00d890d6528de47d16f843f35fd6de453d3e7219ee97eeacbf3eed2ba6f6a4&scene=178&cur_album_id=1509674724564500480#rd)