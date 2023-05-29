<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [自定义New结构体（option模式）](#%E8%87%AA%E5%AE%9A%E4%B9%89new%E7%BB%93%E6%9E%84%E4%BD%93option%E6%A8%A1%E5%BC%8F)
  - [在zap源码中实现](#%E5%9C%A8zap%E6%BA%90%E7%A0%81%E4%B8%AD%E5%AE%9E%E7%8E%B0)
    - [结构体logger](#%E7%BB%93%E6%9E%84%E4%BD%93logger)
    - [初始化](#%E5%88%9D%E5%A7%8B%E5%8C%96)
    - [option](#option)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 自定义New结构体（option模式）

## 在zap源码中实现
### 结构体logger
```go
type Logger struct {
	core zapcore.Core  //定义了输出日志核心接口

	development bool
	addCaller   bool //是否输出调用者的信息
	onFatal     zapcore.CheckWriteAction // default is WriteThenFatal

	name        string
	errorOutput zapcore.WriteSyncer // 错误输出终端，注意区别于zapcore中的输出，这里一般是指做运行过程中，发生错误记录日志（如：参数错误，未定义错误等），默认是os.Stderr

	addStack zapcore.LevelEnabler //需要记录stack信息的日志级别

	callerSkip int //调用者的层级：用于指定记录哪个调用者信息
}
```
Note: 定义了与输出相关的基本信息，比如：name，stack，core等，我们可以看到这些属性都是不对外公开的，所以不能直接初始化结构体.
zap为我们提供了New，Build两种方式来初始化Logger。除了core以外，其他的都可以通过Option来设置。

### 初始化
```go
func New(core zapcore.Core, options ...Option) *Logger {
	if core == nil {
		return NewNop()
	}
	log := &Logger{
		core:        core,
		errorOutput: zapcore.Lock(os.Stderr),
		addStack:    zapcore.FatalLevel + 1,
	}
	return log.WithOptions(options...)
}
```

### option
```go
// An Option configures a Logger.
type Option interface {
	apply(*Logger)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*Logger)

func (f optionFunc) apply(log *Logger) {
	f(log)
}

// 以下是实现option
//设置成开发模式
func Development() Option {
    return optionFunc(func(log *Logger) {
        log.development = true
    })
}
//自定义错误输出路径
func ErrorOutput(w zapcore.WriteSyncer) Option {
    return optionFunc(func(log *Logger) {
        log.errorOutput = w
    })
}
//Logger的结构化字段，每条日志都会打印这些Filed信息
func Fields(fs ...Field) Option {
    return optionFunc(func(log *Logger) {
        log.core = log.core.With(fs)
    })
}
//日志添加调用者信息
func AddCaller() Option {
    return WithCaller(true)
}
func WithCaller(enabled bool) Option {
    return optionFunc(func(log *Logger) {
        log.addCaller = enabled
    })
}

//设置skip，用户runtime.Caller的参数
func AddCallerSkip(skip int) Option {
    return optionFunc(func(log *Logger) {
        log.callerSkip += skip
    })
}

//设置stack
func AddStacktrace(lvl zapcore.LevelEnabler) Option {
    return optionFunc(func(log *Logger) {
        log.addStack = lvl
    })
}
```
zap还为我们添加了hook，让我们在每次打印日志的时候，可以调用hook方法：比如可以统计打印日志的次数、统计打印字段等.
```go
func Hooks(hooks ...func(zapcore.Entry) error) Option {
    return optionFunc(func(log *Logger) {
        log.core = zapcore.RegisterHooks(log.core, hooks...)
    })
}
```
