# pflag

pflag是Go的flag包的直接替代，实现了POSIX / GNU样式的–flags。pflag是Go的本机标志包的直接替代。
如果您在名称“ flag”下导入pflag，则所有代码应继续运行且无需更改。

flag和pflag都是源自于Google，工作原理甚至代码实现基本上都是一样的。
flag虽然是Golang官方的命令行参数解析库，但是pflag却得到更加广泛的应用。 因为pflag相对flag有如下优势
- 支持更加精细的参数类型：例如，flag只支持uint和uint64，而pflag额外支持uint8、uint16、int32。
- 支持更多参数类型：ip、ip mask、ip net、count、以及所有类型的slice类型，例如string slice、int slice、ip slice等。
- 兼容标准flag库的Flag和FlagSet。
- 原生支持更丰富flag功能：shorthand、deprecated、hidden等高级功能。


## shorthand 用法
与 flag 包不同，在 pflag 包中，选项名称前面的 –- 和 - 是不一样的。- 表示 shorthand，-- 表示完整的选项名称

1. 对于布尔类型的参数和设置了 NoOptDefVal 的参数可以写成下面的形式：
```shell
-o
-o=true
// 注意，下面的写法是不正确的
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

## NoOptDefVal 用法

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

## Normalize 用法:标准化参数的名称

如果我们创建了名称为 --des-detail 的参数，但是用户却在传参时写成了 --des_detail 或 --des.detail 会怎么样？
默认情况下程序会报错退出，但是我们可以通过 pflag 提供的 SetNormalizeFunc 功能轻松的解决这个问题


下面的写法也能正确设置参数了
```shell
--des_detail="person detail"
```

## deprecated 用法:把参数标记为即将废弃

在程序的不断升级中添加新的参数和废弃旧的参数都是常见的用例，pflag 包对废弃参数也提供了很好的支持。
通过 MarkDeprecated 和 MarkShorthandDeprecated 方法可以分别把参数及其 shorthand 标记为废弃：
```shell
➜  02_pflag git:(feature/flag) ✗ go run pflag.go -b test                     
Flag shorthand -b has been deprecated, please use -d instead
Flag --badflag has been deprecated, please use --des-detail instead
name= nick
age= 22
gender= male
ok= false
des= 

```

## hidden 用法:在帮助文档中隐藏参数