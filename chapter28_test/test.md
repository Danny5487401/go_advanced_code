<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Golang Testing](#golang-testing)
  - [背景](#%E8%83%8C%E6%99%AF)
  - [基本概念](#%E5%9F%BA%E6%9C%AC%E6%A6%82%E5%BF%B5)
    - [TestDouble 简单理解就是测试替身](#testdouble-%E7%AE%80%E5%8D%95%E7%90%86%E8%A7%A3%E5%B0%B1%E6%98%AF%E6%B5%8B%E8%AF%95%E6%9B%BF%E8%BA%AB)
    - [Unit Test 单元测试](#unit-test-%E5%8D%95%E5%85%83%E6%B5%8B%E8%AF%95)
  - [覆盖 cover 测试](#%E8%A6%86%E7%9B%96-cover-%E6%B5%8B%E8%AF%95)
  - [go test 命令行参数](#go-test-%E5%91%BD%E4%BB%A4%E8%A1%8C%E5%8F%82%E6%95%B0)
    - [1 常规语法](#1-%E5%B8%B8%E8%A7%84%E8%AF%AD%E6%B3%95)
    - [2 执行特定的测试用例](#2-%E6%89%A7%E8%A1%8C%E7%89%B9%E5%AE%9A%E7%9A%84%E6%B5%8B%E8%AF%95%E7%94%A8%E4%BE%8B)
    - [3 执行覆盖测试](#3-%E6%89%A7%E8%A1%8C%E8%A6%86%E7%9B%96%E6%B5%8B%E8%AF%95)
    - [4 在测试时检测数据竞争问题](#4-%E5%9C%A8%E6%B5%8B%E8%AF%95%E6%97%B6%E6%A3%80%E6%B5%8B%E6%95%B0%E6%8D%AE%E7%AB%9E%E4%BA%89%E9%97%AE%E9%A2%98)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Golang Testing


## 背景

不写单元测试（Unit Test,UT)也许能使代码交付更快，但是我们无法保证写出来的代码真的能够正确地执行。写UT可以减少后期解决bug的时间，也能让我们放心地使用自己写出来的代码。从长远来看，后者更能有效地节省开发时间。

既然UT这么重要，是什么原因在阻止开发人员写UT呢？这是因为除了开发人员的惰性习惯之外，编写UT代码同样存在难点。

- 代码耦合度高，缺少必要的抽象与拆分，以至于不知道如何写UT。
- 存在第三方依赖，例如依赖数据库连接、HTTP请求、数据缓存等。



## 基本概念
| 术语 | 解释 |
|--|--|
| SUT（System Under Test） | 被测系统 |
| DOC（depended-on component） | 第三方依赖组件 |
| BDD - Behavior-Driven Development） | 行为驱动开发 |
| Test-Driven Development (TDD) | 测试驱动开发，侧重于系统的实现|


### TestDouble 简单理解就是测试替身

TestDouble 简单理解就是测试替身，在多数情况下，我们的系统能够正常运行，不仅仅依托系统本身，还需要依赖一些外部服务，比如其他系统提供的 http、rpc 服务，系统自身以来的像 redis 缓存服务或者 mysql 这类数据库服务。
在微服务场景下，业务按照业务领域将一个系统拆分为多个系统，系统之间的交互不仅仅是简单的 A->B，可能是 A ->B -> C ->D，对于编写单元测试的开发者来说，当我需要编写系统A 的测试用例时，不可能去构建完整的调用链路，
那么在测试工程中，通常会以 “测试替身” 来解决外部依赖所带来的测试复杂性问题

测试替身主要包括以下几种类型
- Dummy Object：虚拟对象，本质上不会对测试产生任何影响，实际上只作为类似参数填充类角色存在
- Test Stub：SUT 内部打的一个桩，可以按照我们的要求返回特定的内容给 SUT，Test Stub 的交互完全在 SUT 内部
- Test Spy：专门负责将 SUT 内部的间接输出(indirect outputs)传到外部。它的特点是将内部的间接输出返回给测试案例，由测试案例进行验证，Test Spy 只负责获取内部情报，并把情报发出去，不负责验证情报的正确性。
- Mock Object：Mock Object 还负责对情报(indirect outputs)进行验证，总部(外部的测试案例)信任 Mock Object 的验证结果。
- Fake Object：并不关注 SUT 内部的间接输入(indirect inputs)或间接输出(indirect outputs)，它仅仅是用来替代一个实际的对象，并且拥有几乎和实际对象一样的功能，保证 SUT 能够正常工作


### Unit Test 单元测试
单元测试是针对任意一个具体的函数而言，无论是一个已导出的函数接口，或者是一个并不导出的内部工具函数，你可以针对这个函数做一组测试，目的在于证明该函数的功用与其所宣称的相同

对于 Golang 来说，编写单元测试很容易：

- 在一个包例如 yy 之中新建一个 go 源文件，确保文件名以 _test.go 结尾，例如 yy_test.go
- 在这个文件中可以使用 yy 或者 yy_test 作为包名
- 编写一个测试函数入口，其签名必以 Test开头，参数必须是 t *testing.T （对于性能测试函数来说是 b *benchmark.B）
- 在函数体中编写测试代码，如果认为测试不通过，采用 t.Fatal("...") 的方式抛出异常；如果没有异常正常地结束了函数体的运行，则被视作测试已通过。
- 执行过程中可以使用 t.Log(...) 等方式输出日志文本。类似地 t.Fatal 也会输出日志文件，以报错的形式



## 覆盖 cover 测试
覆盖测试是单元测试的一种，我们期待的是对代码的测试覆盖率越高越好。

## go test 命令行参数

1 go help test 命令
```shell
✗ go help test                                                                             
usage: go test [build/test flags] [packages] [build/test flags & test binary flags]

'Go test' automates testing the packages named by the import paths.
It prints a summary of the test results in the format:

        ok   archive/tar   0.011s
        FAIL archive/zip   0.022s
        ok   compress/gzip 0.033s
        ...

followed by detailed output for each failed package.

'Go test' recompiles each package along with any files with names matching
the file pattern "*_test.go".
These additional files can contain test functions, benchmark functions, fuzz
tests and example functions. See 'go help testfunc' for more.
Each listed package causes the execution of a separate test binary.
Files whose names begin with "_" (including "_test.go") or "." are ignored.

Test files that declare a package with the suffix "_test" will be compiled as a
separate package, and then linked and run with the main test binary.

The go tool will ignore a directory named "testdata", making it available
to hold ancillary data needed by the tests.

As part of building a test binary, go test runs go vet on the package
and its test source files to identify significant problems. If go vet
finds any problems, go test reports those and does not run the test
binary. Only a high-confidence subset of the default go vet checks are
used. That subset is: 'atomic', 'bool', 'buildtags', 'errorsas',
'ifaceassert', 'nilfunc', 'printf', and 'stringintconv'. You can see
the documentation for these and other vet tests via "go doc cmd/vet".
To disable the running of go vet, user the -vet=off flag. To run all
checks, user the -vet=all flag.

All test output and summary lines are printed to the go command's
standard output, even if the test printed them to its own standard
error. (The go command's standard error is reserved for printing
errors building the tests.)

Go test runs in two different modes:

The first, called local directory mode, occurs when go test is
invoked with no package arguments (for example, 'go test' or 'go
test -v'). In this mode, go test compiles the package sources and
tests found in the current directory and then runs the resulting
test binary. In this mode, caching (discussed below) is disabled.
After the package test finishes, go test prints a summary line
showing the test status ('ok' or 'FAIL'), package name, and elapsed
time.

The second, called package list mode, occurs when go test is invoked
with explicit package arguments (for example 'go test math', 'go
test ./...', and even 'go test .'). In this mode, go test compiles
and tests each of the packages listed on the command line. If a
package test passes, go test prints only the final 'ok' summary
line. If a package test fails, go test prints the full test output.
If invoked with the -bench or -v flag, go test prints the full
output even for passing package tests, in order to display the
requested benchmark results or verbose logging. After the package
tests for all of the listed packages finish, and their output is
printed, go test prints a final 'FAIL' status if any package test
has failed.

In package list mode only, go test caches successful package test
results to avoid unnecessary repeated running of tests. When the
result of a test can be recovered from the cache, go test will
redisplay the previous output instead of running the test binary
again. When this happens, go test prints '(cached)' in place of the
elapsed time in the summary line.

The rule for a match in the cache is that the run involves the same
test binary and the flags on the command line come entirely from a
restricted set of 'cacheable' test flags, defined as -benchtime, -cpu,
-list, -parallel, -run, -short, -timeout, -failfast, and -v.
If a run of go test has any test or non-test flags outside this set,
the result is not cached. To disable test caching, user any test flag
or argument other than the cacheable flags. The idiomatic way to disable
test caching explicitly is to user -count=1. Tests that open files within
the package's source root (usually $GOPATH) or that consult environment
variables only match future runs in which the files and environment
variables are unchanged. A cached test result is treated as executing
in no time at all,so a successful package test result will be cached and
reused regardless of -timeout setting.

In addition to the build flags, the flags handled by 'go test' itself are:

        -args
           -args 该参数会被原封不动 (uninterpreted and unchanged) 的传递给测试二进制文件
            Because this flag consumes the remainder of the command line,
            the package list (if present) must appear before this flag.

        -c
            只编译不执行
            (where pkg is the last element of the package's import path).
            The file name can be changed with the -o flag.

        -exec xprog
            使用xprog来执行测试二进制程序. The behavior is the same as
            in 'go run'. See 'go help run' for details.

        -i
            安装测试依赖的包，不执行测试
            The -i flag is deprecated. Compiled packages are cached automatically.

        -json
            将输出格式化为json格式，用于自动化处理结果
            See 'go doc test2json' for the encoding details.

        -o file
            Compile the test binary to the named file.
            The test still runs (unless -c or -i is specified).

The test binary also accepts flags that control execution of the test; these
flags are also accessible by 'go test'. See 'go help testflag' for details.

For more about build flags, see 'go help build'.
For more about specifying packages, see 'go help packages'.

See also: go build, go vet.
```


2 go help testflag
```shell
✗ go help testflag
The 'go test' command takes both flags that apply to 'go test' itself
and flags that apply to the resulting test binary.

Several of the flags control profiling and write an execution profile
suitable for "go tool pprof"; run "go tool pprof -h" for more
information. The --alloc_space, --alloc_objects, and --show_bytes
options of pprof control how the information is presented.

The following flags are recognized by the 'go test' command and
control the execution of any test:

```

go test运行时，跟据是否指定package分为两种模式，即本地目录模式和包列表模式。

### 1 常规语法
```shell
# 在当前项目当前包文件夹下执行全部测试用例，但不递归子目录
go test .
# 在当前项目当前文件夹下执行全部测试用例并显示测试过程中的日志内容，不递归子目录
go test -v .

# 和 go test . 相似，但也递归子目录中的一切测试用例
go test ./...
go test -v ./...
```


### 2 执行特定的测试用例
```shell
1
go test -v . -test.run '^TestOne$'
```

### 3 执行覆盖测试
```shell
# 以下两句连用以生成覆盖测试报告 cover.html
go test -v . -coverprofile=coverage.txt -covermode=atomic
go tool cover -html=coverage.txt -o cover.html

# 也可以执行最长的用例执行时间，超出时则判为测试失败
go test -v . -coverprofile=coverage.txt -covermode=atomic -timeout=20m
```

### 4 在测试时检测数据竞争问题

```shell

go test -v -race .
```


## 参考

- 深入探討 Test Double、Dummy、Fake、Stub 、Mock 與 Spy: https://old-oomusou.goodjack.tw/jasmine/jasmine-test-double/



