# 建造者模式
将一个复杂对象的构造与它的表示分离，使同样的构建过程可以创建不同的表示。主要解决多变参数传递问题

## 背景：

当一个方法有多个变量的时候，我们在调用该方法的时候可能会因为参数的顺序、个数错误，而造成调用错误或者不能达到我们预期的目的

## 角色
- Product 这是我们要创建的复杂对象(一般都是很复杂的对象才需要使用建造者模式)。
- Builder 抽象的一个东西， 主要是用来规范我们的建造者的。
- ConcreteBuilder 具体的Builder实现， 这是今天的重点，主要用来根据不用的业务来完成要创建对象的组建的创建。
- Director 这个的作用是规范复杂对象的创建流程

## 案例
### 1. xorm就是使用了builder设计模式

    某业务系统, 希望使用SQLQuery类动态构造复杂SQL查询语句
    SQLQuery类的各种属性组合情况很多, 因此创建SQLQueryBuilder作为SQLQuery的建造者
### 2. k8s
```go
// k8s.io/cli-runtime/pkg/resource/builder.go
r := f.NewBuilder().
    Unstructured().
    Schema(schema).
    ContinueOnError().
    NamespaceParam(cmdNamespace).DefaultNamespace().
    // 读取文件信息，发现除了支持简单的本地文件，也支持标准输入和http/https协议访问的文件
    FilenameParam(enforceNamespace, &o.FilenameOptions).
    LabelSelectorParam(o.Selector).
    Flatten()
```
### 3. zap
```go
//开发环境下的Logger
func NewDevelopment(options ...Option) (*Logger, error) {
    return NewDevelopmentConfig().Build(options...)
}
//生产环境下的Logger
func NewProduction(options ...Option) (*Logger, error) {
    return NewProductionConfig().Build(options...)
}
//测试环境下的Logger
func NewExample(options ...Option) *Logger {
    encoderCfg := zapcore.EncoderConfig{
        MessageKey:     "msg",
        LevelKey:       "level",
        NameKey:        "logger",
        EncodeLevel:    zapcore.LowercaseLevelEncoder,
        EncodeTime:     zapcore.ISO8601TimeEncoder,
        EncodeDuration: zapcore.StringDurationEncoder,
    }
    core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, DebugLevel)
    return New(core).WithOptions(options...)
}
```
不同的构造方式，唯一不同的就是Config.
```go

//Config这个结构体每个字段都有json和yaml的标注, 也就是说这些配置不仅仅可以在代码中赋值，也可以从配置文件中直接反序列化得到
type Config struct {
	// Level是用来配置日志级别的，即日志的最低输出级别，这里的AtomicLevel虽然是个结构体，但是如果使用配置文件直接反序列化
	Level AtomicLevel `json:"level" yaml:"level"`
	// 这个字段的含义是用来标记是否为开发者模式，在开发者模式下，日志输出的一些行为会和生产环境上不同
	Development bool `json:"development" yaml:"development"`
	// 用来标记是否开启行号和文件名显示功能。
	// 默认都是标记的
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`
	// 标记是否开启调用栈追踪能力，即在打印异常日志时，是否打印调用栈. 
	//By default, stacktraces are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace"`
	// Sampling实现了日志的流控功能，或者叫采样配置，主要有两个配置参数，Initial和Thereafter，
	//实现的效果是在1s的时间单位内，如果某个日志级别下同样内容的日志输出数量超过了Initial的数量，
	//那么超过之后，每隔Thereafter的数量，才会再输出一次。是一个对日志输出的保护功能
	Sampling *SamplingConfig `json:"sampling" yaml:"sampling"`
	
	// 用来指定日志的编码器，也就是用户在调用日志打印接口时，zap内部使用什么样的编码器将日志信息编码为日志条目，日志的编码也是日志组件的一个重点。
	// 默认支持两种配置，json和console，用户可以自行实现自己需要的编码器并注册进日志组件，实现自定义编码的能力
	Encoding string `json:"encoding" yaml:"encoding"`
	

	// 定义了输出样式
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`
	
	// 用来指定日志的输出路径，不过这个路径不仅仅支持文件路径和标准输出，还支持其他的自定义协议，
	//当然如果要使用自定义协议，也需要使用RegisterSink方法先注册一个该协议对应的工厂方法，该工厂方法实现了Sink接口
	OutputPaths []string `json:"outputPaths" yaml:"outputPaths"`
	
	// 与OutputPaths类似，不过用来指定的是错误日志的输出，不过要注意，这个错误日志不是业务的错误日志，
	// 而是zap中出现的内部错误，将会被定向到这个路径下.
	ErrorOutputPaths []string `json:"errorOutputPaths" yaml:"errorOutputPaths"`
	
	// 初始化的Fields，每行日志都会爱上这些Field
	InitialFields map[string]interface{} `json:"initialFields" yaml:"initialFields"`
}
```
zapcore.EncoderConfig
```go
type EncoderConfig struct {
    //*Key：设置的是在结构化输出时，value对应的key
    MessageKey    string `json:"messageKey" yaml:"messageKey"`
    LevelKey      string `json:"levelKey" yaml:"levelKey"`
    TimeKey       string `json:"timeKey" yaml:"timeKey"`
    NameKey       string `json:"nameKey" yaml:"nameKey"`
    CallerKey     string `json:"callerKey" yaml:"callerKey"`
    StacktraceKey string `json:"stacktraceKey" yaml:"stacktraceKey"`
    //日志的结束符
    LineEnding    string `json:"lineEnding" yaml:"lineEnding"`
    
    //Level的输出样式，比如 大小写，颜色等
    EncodeLevel    LevelEncoder    `json:"levelEncoder" yaml:"levelEncoder"`
    
    //日志时间的输出样式
    EncodeTime     TimeEncoder     `json:"timeEncoder" yaml:"timeEncoder"`
    
    //消耗时间的输出样式
    EncodeDuration DurationEncoder `json:"durationEncoder" yaml:"durationEncoder"`
    
    //Caller的输出样式，比如 全名称，短名称
    EncodeCaller   CallerEncoder   `json:"callerEncoder" yaml:"callerEncoder"`
    
    // Unlike the other primitive type encoders, EncodeName is optional. The
    // zero value falls back to FullNameEncoder.
    EncodeName NameEncoder `json:"nameEncoder" yaml:"nameEncoder"`
}
```
那开发环境进行讲解
```go
// 初始化
func NewDevelopment(options ...Option) (*Logger, error) {
	return NewDevelopmentConfig().Build(options...)
}


// 初始化配置
func NewDevelopmentConfig() Config {
	return Config{
		Level:            NewAtomicLevelAt(DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    NewDevelopmentEncoderConfig(), //序列化配置
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// 序列化配置
func NewDevelopmentEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
```
开始构造
```go
func (cfg Config) Build(opts ...Option) (*Logger, error) {
	enc, err := cfg.buildEncoder()
	if err != nil {
		return nil, err
	}

	sink, errSink, err := cfg.openSinks()
	if err != nil {
		return nil, err
	}

	if cfg.Level == (AtomicLevel{}) {
		return nil, fmt.Errorf("missing Level")
	}

	log := New(
		zapcore.NewCore(enc, sink, cfg.Level),
		cfg.buildOptions(errSink)...,
	)
	if len(opts) > 0 {
		log = log.WithOptions(opts...)
	}
	return log, nil
}

//  构建编码
func (cfg Config) buildEncoder() (zapcore.Encoder, error) {
   return newEncoder(cfg.Encoding, cfg.EncoderConfig)
}

// 打开路径
func (cfg Config) openSinks() (zapcore.WriteSyncer, zapcore.WriteSyncer, error) {
	sink, closeOut, err := Open(cfg.OutputPaths...)
	if err != nil {
		return nil, nil, err
	}
	errSink, _, err := Open(cfg.ErrorOutputPaths...)
	if err != nil {
		closeOut()
		return nil, nil, err
	}
	return sink, errSink, nil
}
```

## 故事：
1. 平时去面馆吃面，有各种味道的面条（牛肉味、肥肠味等）
   - 有各种配料（香菜、葱、姜、辣椒等）
   - 第一个客人：一碗牛肉面 加葱、姜
   - 第二个客人：一碗牛肉面 加葱、香菜
2. 例如肯德基的点餐系统，汉堡，薯条，可乐，炸鸡是不变的，变化的是他们组合出的套餐，所以点餐系统使用建造者模式，非常容易拓展出不同的套餐，
    而且既定的套餐，修改薯条为大份的时候，也非常方便

## 建造者模式的优点
（1）封装性好，构建和表示分离。
（2）扩展性好，建造类之间独立，在一定程度上解耦。
（3）便于控制细节，建造者可以对创建过程逐步细化，而不对其他模块产生任何影响。
## 建造者模式的缺点
（1）需要多创建一个IBuilder对象。
（2）如果产品内部发生变化，则建造者也要同步修改，后期维护成本较大。

