<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [go tool](#go-tool)
  - [go tool objdump 命令用来根据目标文件或二进制文件反编译出汇编代码。该命令支持两个选项：](#go-tool-objdump-%E5%91%BD%E4%BB%A4%E7%94%A8%E6%9D%A5%E6%A0%B9%E6%8D%AE%E7%9B%AE%E6%A0%87%E6%96%87%E4%BB%B6%E6%88%96%E4%BA%8C%E8%BF%9B%E5%88%B6%E6%96%87%E4%BB%B6%E5%8F%8D%E7%BC%96%E8%AF%91%E5%87%BA%E6%B1%87%E7%BC%96%E4%BB%A3%E7%A0%81%E8%AF%A5%E5%91%BD%E4%BB%A4%E6%94%AF%E6%8C%81%E4%B8%A4%E4%B8%AA%E9%80%89%E9%A1%B9)
- [go build](#go-build)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## go tool


| 名称           | 含义                                                                                                                                                                                                                                                                                                                                                                                                                                             |
|--------------| -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| tool cgo     | cgo 用于支持 Go 包调用 C 代码                                                                                                                                                                                                                                                                                                                                                                                                                    |
| tool cover   | cover 是一个程序，用于创建和分析覆盖率分析信息，由`"go test -coverprofile"` 生成。                                                                                                                                                                                                                                                                                                                                                               |
| tool fix     | fix 程序找到使用语言和库的旧功能的 Go 程序，并以较新的 Go 语言重写。                                                                                                                                                                                                                                                                                                                                                                             |
| fmt          | fmt 格式化 go 包，它也可以作为一个独立的 gofmt 命令与更多的通用选项                                                                                                                                                                                                                                                                                                                                                                              |
| tool doc     | godoc 为 Go 包提取并生成文档                                                                                                                                                                                                                                                                                                                                                                                                                     |
| tool vet     | vet 检查 Go 源代码并报告可疑的构造，如参数与格式字符串不匹配的Printf调用。                                                                                                                                                                                                                                                                                                                                                                       |
| tool asm     | asm（通常被调用为 go tool asm）将源文件（汇编代码文件）组装到一个对象文件中，该对象文件以参数源文件的基名命名，后缀为.o。然后，可以将对象文件与其他对象组合到包归档文件中。                                                                                                                                                                                                                                                                      |
| tool compile | compile（通常作为“go tool compile”调用）编译多个文件组成的Go包。然后，写入一个对象文件，以第一个源码文件名的对象文件（ .o）。然后，可以将对象文件与其他对象组合到包归档文件中，或者直接传递到链接器（“go tool link”）。如果使用-pack调用，编译器将跳过中间对象文件，直接写入归档文件。生成的文件包含包导出的符号的类型信息，以及包从其他包导入的符号使用的类型信息。因此，在编译P包的客户端C时，不必读取P的依赖项文件，只需读取P的编译输出。 |
| tool link    | link（通常作为“go tool link”调用）读取包主目录的Go归档文件或对象及其依赖项，并将其组合到可执行二进制文件中。                                                                                                                                                                                                                                                                                                                                   |
| tool pack    | pack 是传统 Unix ar 工具的简单版本。                                                                                                                                                                                                                                                                                                                                                                                                             |
| tool pprof   | pprof 解释和显示 Go 程序的概要文件。                                                                                                                                                                                                                                                                                                                                                                                                             |
| tool test2json | Test2json 将 go test 输出转换为机器可读的JSON流。                                                                                                                                                                                                                                                                                                                                                                                                |
| tool trace   | trace 是查看跟踪文件的工具。                                                                                                                                                                                                                                                                                                                                                                                                                     |
| tool nm      | 查看Go 二进制文件中符号表信息。                                                                                                                                                                                                                                                                                                                                                                                                                     |
| tool  objdump        | 根据目标文件或二进制文件反编译出汇编代码                                                                                                                                                                                                                                                                                                                                                                                                                    |


### go tool objdump 命令用来根据目标文件或二进制文件反编译出汇编代码。该命令支持两个选项：

-S选项指示打印汇编代码
-s选项指示搜索相关的汇编代码
```go
go tool compile -N -l main.go # 生成main.o
go tool objdump main.o # 打印所有汇编代码
go tool objdump -s "main.(main|add)" ./test # objdump支持搜索特定字符串
```

## go build
```go
✗ go help build                                  
usage: go build [-o output] [build flags] [packages]
```

编译时的参数传递（gcflags, ldflags）
1. -x 列出 build 过程中用到的所有工具，mkdir/gcc等等
2. -n 不实际编译，仅打印出 build 过程中用到的所有工具.通过该选项我们可以查看编译器，连接器如何工作的
3. -a 全部重新构建（命令源码文件与库源码文件）
4. -race 竞争检测
5. -gcflags:给go编译器传入参数，也就是传给go tool compile的参数，因此可以用go tool compile --help查看所有可用的参数。
    * -gcflags="all=-N -l"
    * -N 取消优化
    * -l 取消内联
    * -m 逃逸分析，打印逃逸信息
    * -gcflags=-S fmt 仅打印fmt包的反汇编信息
    * -gcflags=all=-S fmt' 打印fmt以及其依赖包的反汇编信息
6. go build -buildmode - plugin 编译成.so插件，通过包 plugin 进行打开，获取符号 - c-shared：使用该参数时会生成出来两个文件，一个.so文件，一个.h头文件 ，使用起来就和使用c 生成的库文件和模块文件一样使用。
7. -ldflags: 给go链接器传入参数，实际是给go tool link的参数，可以用go tool link --help查看可用的参数
    * -X来指定版本号等编译时才决定的参数值。例如代码中定义var buildVer string，然后在编译时用go build -ldflags "-X main.buildVer=1.0" ... 来赋值。注意-X只能给string类型变量赋值






