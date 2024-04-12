<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [generate](#generate)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# generate 

go generate用于一键式批量执行任何命令，创建或更新Go文件或者输出结果

```shell
$ go help generate
usage: go generate [-run regexp] [-n] [-v] [-x] [build flags] [file.go... | packages]

Generate runs commands described by directives within existing
files. Those commands can run any process but the intent is to
create or update Go source files.

Go generate is never run automatically by go build, go test,
and so on. It must be run explicitly.

Go generate scans the file for directives, which are lines of
the form,

//go:generate command argument...

```

参数说明：

* -run 正则表达式匹配命令行，仅执行匹配的命令
* -v 打印已被检索处理的文件。
* -n 打印出将被执行的命令，此时将不真实执行命令
* -x 打印已执行的命令


```shell
$GOARCH
	架构 (arm, amd64, etc.)
$GOOS
	OS (linux, windows, etc.)
$GOFILE
	当前处理中的文件名
$GOLINE
	当前命令在文件中的行号
$GOPACKAGE
    当前处理文件的包名

```


## 参考

- [深入理解Go之generate](https://darjun.github.io/2019/08/21/golang-generate/)