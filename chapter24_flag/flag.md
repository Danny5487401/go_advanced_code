<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [flag命令行](#flag%E5%91%BD%E4%BB%A4%E8%A1%8C)
  - [基本概念](#%E5%9F%BA%E6%9C%AC%E6%A6%82%E5%BF%B5)
  - [标准包flag包源码分析](#%E6%A0%87%E5%87%86%E5%8C%85flag%E5%8C%85%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [有两重要的数据结构来表示flag](#%E6%9C%89%E4%B8%A4%E9%87%8D%E8%A6%81%E7%9A%84%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84%E6%9D%A5%E8%A1%A8%E7%A4%BAflag)
    - [参数解析](#%E5%8F%82%E6%95%B0%E8%A7%A3%E6%9E%90)
    - [错误处理方式](#%E9%94%99%E8%AF%AF%E5%A4%84%E7%90%86%E6%96%B9%E5%BC%8F)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# flag命令行

命令行参数用于向应用程序传递一些定制参数，使得程序的功能更加丰富和多样化。
命令行标志是一类特殊的命令行参数，通常以减号（-）或双减号（–）连接标志名称，非bool类型的标志后面还会有取值。

以git log命令为例，例如我们要观察最近的10条commit记录，且要显示每条记录修改的文件信息：
```go
git log --stat -n 10
```
其中的--stat和-n 10就是两个标志，前者是bool类型，后者是int类型。
* --stat告诉git log输出每条commit记录的统计信息，它的取值只有True或False，因此标志后不需要其它参数值。
* -n 10则需要通过后面的数值10告诉git log命令我们需要显示的commit记录条数。

非bool类型的标志值还可以通过等号的形式提供，比如下面这条命令也用于显示最近的10条commit记录


## 基本概念

- 命令行参数（或参数）：是指运行程序提供的参数
- 已定义命令行参数：是指程序中通过flag.Xxx等这种形式定义了的参数
- 非flag（non-flag）命令行参数（或保留的命令行参数）
  - ./nginx - -c 或 ./nginx build -c,这两种情况，-c 都不会被正确解析。像该例子中的"-"或build（以及之后的参数），我们称之为 non-flag 参数。

## 标准包flag包源码分析
flag 用法有如下三种形式
```shell
-flag // 只支持bool类型
-flag=x
-flag x // 只支持非bool类型
```

Golang的命令行参数解析使用的是flag包，支持布尔、整型、字符串，以及时间格式的标识解析.
Flag原生支持8种命令行参数类型：bool、int、int64、uint、uint64、string、float64和duration。 这些类型都需要实现Value接口：
```go
type Value interface {
	String() string // 将该类型变量的值，转化为string
	Set(string) error //将string类型的value解析出来，并赋值给该类型的变量；
}
```

下面我们以一个echoflag程序为例，演示flag包的用法。这个程序接收来自命令行的输入，并回显命令行标识的值，程序的执行效果如下：
```go
var bval = flag.Bool("bool", false, "bool value for test")
var ival = flag.Int("int", 100, "integer value for test")
var sval = flag.String("string", "null", "string value for test")
var tval = flag.Duration("time", 10*time.Second, "time duration for test")

func main() {

	flag.Parse()

	fmt.Println("bool:\t", *bval)
	fmt.Println("int:\t", *ival)
	fmt.Println("string:\t", *sval)
	fmt.Println("time:\t", *tval)
}
```
```shell
$ go build -o main main.go
$ ./main -bool -int 10 --string "string for test" --time=100s argv
bool:    true
int:     10
string:  string for test
time:    1m40s
```
- 布尔类型的参数仅有标记没有取值，指定bool表示标记为True，不指定就是False
- 整型参数的取值为10
- 字符串参数的取值为string for test”，注意我们这里使用了双连接线，这与-string效果是一样的
- 标记与标记值之间可以用空格或等号分隔，如时间参数我们则使用了等号
- 最后一个参数argv不带连接线，因此被当做普通参数处理，flag包对此不做解析

### 有两重要的数据结构来表示flag
1. FlagSet
```go
type FlagSet struct {
    Usage func() // 帮助函数，在命令行标志输入不符合预期时被调用，并提示用户正确的输入方式

    name          string //程序的名称，在CommandLine被初始化时赋值为os.Args[0]，也就是应用程序的名称
    parsed        bool
    actual        map[string]*Flag // 解析时， 存放实际传递了的参数（即命令行参数）
    formal        map[string]*Flag // 初始化，存放所有已定义命令行参数 
    args          []string //  开始存放所有参数，最后保留 非flag（non-flag）参数
    errorHandling ErrorHandling // 当解析出错时，处理错误的方式
    output        io.Writer // nil means stderr; use out() accessor
}

//  预定义的 FlagSet 实例 CommandLine 的定义方式，可见，默认的 FlagSet 实例在解析出错时会退出程序。
var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

// NewFlagSet() 用于实例化 FlagSet
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
    f := &FlagSet{
        name:          name,
        errorHandling: errorHandling,
    }
    f.Usage = f.defaultUsage
    return f
}
```

2. Flag
```go
type Flag struct {
	Name     string // Name就是标记名称，也就是命令行输入的类似-int 10中的int
	Usage    string // help message
	Value    Value  // value as set, flag支持的标记值类型
	DefValue string // 默认值 (as text); for usage message
}
type Value interface {
    String() string // 显示命令行标志的名
    Set(string) error //Set则用于记录标志的值
}
```

Flag 类型代表一个 flag 的状态，比如，对于命令：./nginx -c /etc/nginx.conf，相应代码是：
```go
flag.StringVar(&c, "c", "conf/nginx.conf", "set configuration `file`")
```
则该 Flag 实例（可以通过 flag.Lookup(“c”) 获得）相应各个字段的值为：
```go
&Flag{
    Name: c,
    Usage: set configuration file,
    Value: /etc/nginx.conf,
    DefValue: conf/nginx.conf,
}

```


value的其中一个实现:stringValue
```go
type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Get() interface{} { return string(*s) }

func (s *stringValue) String() string { return string(*s) }
```

初始化时    
```go
func StringVar(p *string, name string, value string, usage string) {
	CommandLine.Var(newStringValue(value, p), name, usage)
}

var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	f := &FlagSet{
		name:          name,
		errorHandling: errorHandling,
	}
	f.Usage = f.defaultUsage
	return f
}
```
定义的每个flag，都会添加到CommandLine这个全局的FlagSet中。

```go
func (f *FlagSet) String(name string, value string, usage string) *string {
	p := new(string)
	f.StringVar(p, name, value, usage)
	return p
}

func (f *FlagSet) StringVar(p *string, name string, value string, usage string) {
	f.Var(newStringValue(value, p), name, usage)
}

func (f *FlagSet) Var(value Value, name string, usage string) {
	// Remember the default value as a string; it won't change.
	flag := &Flag{name, usage, value, value.String()}
	_, alreadythere := f.formal[name]
	if alreadythere {
		// 若已存在则panic报错
		var msg string
		if f.name == "" {
			msg = fmt.Sprintf("flag redefined: %s", name)
		} else {
			msg = fmt.Sprintf("%s flag redefined: %s", f.name, name)
		}
		fmt.Fprintln(f.Output(), msg)
		panic(msg) // Happens only if flags are declared with identical names
	}
	if f.formal == nil {
		f.formal = make(map[string]*Flag)
	}
	f.formal[name] = flag
}
```
最终在调用到FlagSet的Var方法时，字符串类型的标志被记录到了CommandLine的formal里面了

### 参数解析
```go
func Parse() {
	// Ignore errors; CommandLine is set for ExitOnError.
	CommandLine.Parse(os.Args[1:])
}


func (f *FlagSet) Parse(arguments []string) error {
	f.parsed = true
	f.args = arguments
	for {
		// 逐个读取标志
		seen, err := f.parseOne()
		if seen {
			continue
		}
        // 最终无错直接退出
		if err == nil {
			break
		}
		switch f.errorHandling {
		case ContinueOnError:
			return err
		case ExitOnError:
			if err == ErrHelp {
				os.Exit(0)
			}
			os.Exit(2)
		case PanicOnError:
			panic(err)
		}
	}
	return nil
}
```

具体解析过程
```go
func (f *FlagSet) parseOne() (bool, error) {
	if len(f.args) == 0 {
		return false, nil
	}
	s := f.args[0]
	if len(s) < 2 || s[0] != '-' {
        // key必须包含’-'，且长度必须大于1
		// 也就是，当遇到单独的一个"-"或不是"-"开始时，会停止解析。比如：./nginx - -c 或 ./nginx build -c
		return false, nil
	}
	numMinuses := 1
	if s[1] == '-' {
		numMinuses++
		if len(s) == 2 { // "--" terminates the flags,比如./nginx -c --
			f.args = f.args[1:]
			return false, nil
		}
	}
	name := s[numMinuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return false, f.failf("bad flag syntax: %s", s)
	}

	// it's a flag. does it have an argument?
	//   每执行成功一次 parseOne，f.args 会少一个
	f.args = f.args[1:]
	
	
	hasValue := false
	value := ""
	for i := 1; i < len(name); i++ { // equals cannot be first
		if name[i] == '=' {
            // 处理中包含'='的问题，以'='为界限拆分为key、value两部分
			value = name[i+1:]
			hasValue = true
			name = name[0:i]
			break
		}
	}
	m := f.formal
	flag, alreadythere := m[name] // BUG
	if !alreadythere {
		// 解析到不存在的参数
		if name == "help" || name == "h" { // special case for nice help message.
			f.usage()
			return false, ErrHelp
		}
		return false, f.failf("flag provided but not defined: -%s", name)
	}

    // 1。接下来是处理-flag=x这种形式，
    // 2。然后是-flag这种形式（bool类型）（这里对bool进行了特殊处理
    // 3。接着是-flag x这种形式，最后将解析成功的Flag实例存入FlagSet的actual map中
    
	if fv, ok := flag.Value.(boolFlag); ok && fv.IsBoolFlag() { // special case: doesn't need an arg
		if hasValue {
			if err := fv.Set(value); err != nil {
				return false, f.failf("invalid boolean value %q for -%s: %v", value, name, err)
			}
		} else {
			if err := fv.Set("true"); err != nil {
				return false, f.failf("invalid boolean flag %s: %v", name, err)
			}
		}
	} else {
		// It must have a value, which might be the next argument.
		if !hasValue && len(f.args) > 0 {
			// value is the next arg
			hasValue = true
			value, f.args = f.args[0], f.args[1:]
		}
		if !hasValue {
			return false, f.failf("flag needs an argument: -%s", name)
		}
        // 设置值
		if err := flag.Value.Set(value); err != nil {
			return false, f.failf("invalid value %q for flag -%s: %v", value, name, err)
		}
	}
	if f.actual == nil {
		f.actual = make(map[string]*Flag)
	}
    // 存储
	f.actual[name] = flag
	return true, nil
}
```

### 错误处理方式
```go
type ErrorHandling int

// These constants cause FlagSet.Parse to behave as described if the parse fails.
const (
	ContinueOnError ErrorHandling = iota // Return a descriptive error.
	ExitOnError                          // Call os.Exit(2) or for -h/-help Exit(0).
	PanicOnError                         // Call panic with a descriptive error.
)
```
三个常量在源码的 FlagSet 的方法 parseOne() 中使用了。
