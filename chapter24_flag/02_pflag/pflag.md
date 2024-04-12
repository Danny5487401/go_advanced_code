<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [pflag](#pflag)
  - [应用](#%E5%BA%94%E7%94%A8)
  - [注意项](#%E6%B3%A8%E6%84%8F%E9%A1%B9)
    - [shorthand 短选项用法](#shorthand-%E7%9F%AD%E9%80%89%E9%A1%B9%E7%94%A8%E6%B3%95)
    - [NoOptDefVal 用法：指定了选项但是没有指定选项值时的默认值](#nooptdefval-%E7%94%A8%E6%B3%95%E6%8C%87%E5%AE%9A%E4%BA%86%E9%80%89%E9%A1%B9%E4%BD%86%E6%98%AF%E6%B2%A1%E6%9C%89%E6%8C%87%E5%AE%9A%E9%80%89%E9%A1%B9%E5%80%BC%E6%97%B6%E7%9A%84%E9%BB%98%E8%AE%A4%E5%80%BC)
  - [Normalize 用法:标准化参数的名称](#normalize-%E7%94%A8%E6%B3%95%E6%A0%87%E5%87%86%E5%8C%96%E5%8F%82%E6%95%B0%E7%9A%84%E5%90%8D%E7%A7%B0)
    - [deprecated 用法:把参数标记为即将废弃](#deprecated-%E7%94%A8%E6%B3%95%E6%8A%8A%E5%8F%82%E6%95%B0%E6%A0%87%E8%AE%B0%E4%B8%BA%E5%8D%B3%E5%B0%86%E5%BA%9F%E5%BC%83)
    - [hidden 用法:在帮助文档中隐藏参数](#hidden-%E7%94%A8%E6%B3%95%E5%9C%A8%E5%B8%AE%E5%8A%A9%E6%96%87%E6%A1%A3%E4%B8%AD%E9%9A%90%E8%97%8F%E5%8F%82%E6%95%B0)
  - [使用](#%E4%BD%BF%E7%94%A8)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [Flag](#flag)
    - [FLagSet](#flagset)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# pflag

pflag是Go的flag包的直接替代，实现了POSIX / GNU样式的–flags。pflag是Go的本机标志包的直接替代。
如果您在名称“ flag”下导入pflag，则所有代码应继续运行且无需更改。

flag和pflag都是源自于Google，工作原理甚至代码实现基本上都是一样的。

flag虽然是Golang官方的命令行参数解析库，但是pflag却得到更加广泛的应用。 因为pflag相对flag有如下优势
- 支持更加精细的参数类型：例如，flag只支持uint和uint64，而pflag额外支持uint8、uint16、int32。
- 支持更多参数类型：ip、ip mask、ip net、count、以及所有类型的slice类型，例如string slice、int slice、ip slice等。
- 兼容标准flag库的Flag和FlagSet。
- 原生支持更丰富flag功能：shorthand、deprecated、hidden等高级功能。

## 应用
Kubernetes、Istio、Helm、Docker、Etcd 等。

## 注意项
### shorthand 短选项用法
与 flag 包不同，在 pflag 包中，选项名称前面的 –- 和 - 是不一样的。- 表示 shorthand，-- 表示完整的选项名称

1. 对于布尔类型的参数和设置了 NoOptDefVal 的参数可以写成下面的形式：
```shell
-o
-o=true
# 注意，下面的写法是不正确的
-o true

```

2. 非布尔类型的参数和 *没有*设置 NoOptDefVal 的参数的写法如下：
```shell
-g female
-g=female
-gfemale
```

3. 注意 –- 后面的参数不会被解析：
```shell
-- -gfemale
```

运行验证
```shell
➜  02_pflag git:(feature/flag) ✗ go run pflag.go -n newName -o false -a 30 -- -gfemale
name= newName
age= 25
gender= male
ok= true
des= 
```

### NoOptDefVal 用法：指定了选项但是没有指定选项值时的默认值

pflag 包支持通过简便的方式为参数设置默认值之外的值，实现方式为设置参数的 NoOptDefVal 属性
```go
var cliAge = flag.IntP("age", "a",22, "Input Your Age")
flag.Lookup("age").NoOptDefVal = "25"
```
下面是传递参数的方式和参数最终的取值：
```shell
Parsed Arguments     Resulting Value
--age=30             cliAge=30
--age                cliAge=25
[nothing]            cliAge=22

```
对于设置了NoOptDefVal的参数， -a 30，这样使用是不正确的

##  Normalize 用法:标准化参数的名称

如果我们创建了名称为 --des-detail 的参数，但是用户却在传参时写成了 --des_detail 或 --des.detail 会怎么样？
默认情况下程序会报错退出，但是我们可以通过 pflag 提供的 SetNormalizeFunc 功能轻松的解决这个问题


下面的写法也能正确设置参数了
```shell
--des_detail="person detail"
```

### deprecated 用法:把参数标记为即将废弃

在程序的不断升级中添加新的参数和废弃旧的参数都是常见的用例，pflag 包对废弃参数也提供了很好的支持。
通过 MarkDeprecated 和 MarkShorthandDeprecated 方法可以分别把参数及其 shorthand 标记为废弃：
```go
	flag.CommandLine.MarkDeprecated("badflag", "please user --des-detail instead")
```

```shell
➜  02_pflag git:(feature/flag) ✗ go run pflag.go -b test                     
Flag shorthand -b has been deprecated, please user -d instead
Flag --badflag has been deprecated, please user --des-detail instead
name= nick
age= 22
gender= male
ok= false
des= 

```

### hidden 用法:在帮助文档中隐藏参数
可以将 Flag 标记为隐藏的，这意味着它仍将正常运行，但不会显示在 usage/help 文本中。
例如：隐藏名为 secretFlag 的标志，只在内部使用，并且不希望它显示在帮助文本或者使用文本中。代码如下
```go

// hide a flag by specifying its name
pflag.CommandLine.MarkHidden("secretFlag")
```

## 使用
1. Pflag 支持以下 4 种命令行参数定义方式：

- 支持长选项、默认值和使用文本，并将标志的值存储在指针中。
```go
var name = pflag.String("name", "colin", "Input Your Name")
```

- 支持长选项、短选项、默认值和使用文本，并将标志的值存储在指针中。

```go
var name = pflag.StringP("name", "n", "colin", "Input Your Name")
```

- 支持长选项、默认值和使用文本，并将标志的值绑定到变量。
```go

var name string
pflag.StringVar(&name, "name", "colin", "Input Your Name")
```

- 支持长选项、短选项、默认值和使用文本，并将标志的值绑定到变量
```go

var name string
pflag.StringVarP(&name, "name", "n","colin", "Input Your Name")
```

上面的函数命名是有规则的：

- 函数名带Var说明是将标志的值绑定到变量，否则是将标志的值存储在指针中。
- 函数名带P说明支持短选项，否则不支持短选项。


2. 使用Get获取参数的值。

可以使用Get来获取标志的值，代表 Pflag 所支持的类型。例如：有一个 pflag.FlagSet，带有一个名为 flagname 的 int 类型的标志，可以使用GetInt()来获取 int 值。需要注意 flagname 必须存在且必须是 int，例如
```go

i, err := flagset.GetInt("flagname")
```

3. 获取非选项参数
```go

package main

import (
    "fmt"

    "github.com/spf13/pflag"
)

var (
    flagvar = pflag.Int("flagname", 1234, "help message for flagname")
)

func main() {
    pflag.Parse()

    fmt.Printf("argument number is: %v\n", pflag.NArg())
    fmt.Printf("argument list is: %v\n", pflag.Args())
    fmt.Printf("the first argument is: %v\n", pflag.Arg(0))
}
```

```shell

$ go run example1.go arg1 arg2
argument number is: 2
argument list is: [arg1 arg2]
the first argument is: arg1
```


## 源码分析

### Flag
Pflag 可以对命令行参数进行处理，一个命令行参数在 Pflag 包中会解析为一个 Flag 类型的变量。Flag 是一个结构体，

```go
// /Users/python/go/pkg/mod/github.com/spf13/pflag@v1.0.5/flag.go
type Flag struct {
    Name                string // flag长选项的名称
    Shorthand           string // flag短选项的名称，一个缩写的字符
    Usage               string // flag的使用文本
    Value               Value  // flag的值
    DefValue            string // flag的默认值
    Changed             bool // 记录flag的值是否有被设置过
    NoOptDefVal         string // 当flag出现在命令行，但是没有指定选项值时的默认值
    Deprecated          string // 记录该flag是否被放弃
    Hidden              bool // 如果值为true，则从help/usage输出信息中隐藏该flag
    ShorthandDeprecated string // 如果flag的短选项被废弃，当使用flag的短选项时打印该信息
    Annotations         map[string][]string // 给flag设置注解
}
```

Flag 的值是一个 Value 类型的接口，Value 定义如下
```go

type Value interface {
    String() string // 将flag类型的值转换为string类型的值，并返回string的内容
    Set(string) error // 将string类型的值转换为flag类型的值，转换失败报错
    Type() string // 返回flag的类型，例如：string、int、ip等
}
```

通过将 Flag 的值抽象成一个 interface 接口，我们就可以自定义 Flag 的类型了。只要实现了 Value 接口的结构体，就是一个新类型。

### FLagSet

Pflag 除了支持单个的 Flag 之外，还支持 FlagSet。FlagSet 是一些预先定义好的 Flag 的集合，几乎所有的 Pflag 操作，都需要借助 FlagSet 提供的方法来完成。
在实际开发中，我们可以使用两种方法来获取并使用 FlagSet

- 方法一，调用 NewFlagSet 创建一个 FlagSet。
```go

var version bool
flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
flagSet.BoolVar(&version, "version", true, "Print version information and quit.")
```
- 方法二，使用 Pflag 包定义的全局 FlagSet：CommandLine。实际上 CommandLine 也是由 NewFlagSet 函数创建的

```go
var CommandLine = NewFlagSet(os.Args[0], ExitOnError)
```
