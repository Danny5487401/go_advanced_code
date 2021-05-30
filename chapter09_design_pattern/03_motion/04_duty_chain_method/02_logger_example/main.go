package main

import (
	"fmt"
)

const (
	Debug = iota
	Info
	Warning
	Error
	Fatal
)

type Logger interface {
	PrintLog(level int, msg string)
}

type BaseLogger struct {
	next Logger
}

func (b *BaseLogger) SetNext(logger Logger) {
	b.next = logger
}

type DebugLogger struct {
	BaseLogger
}

func (d *DebugLogger) PrintLog(level int, msg string) {
	if level == Debug {
		fmt.Printf("[DEBUG] %s\n", msg)
	} else {
		fmt.Printf("ignore [DEBUG]\n")
		d.next.PrintLog(level, msg)
	}
}

type InfoLogger struct {
	BaseLogger
}

func (d *InfoLogger) PrintLog(level int, msg string) {
	if level == Info {
		fmt.Printf("[INFO] %s\n", msg)
	} else {
		fmt.Printf("ignore [INFO]\n")
		d.next.PrintLog(level, msg)
	}
}

type ErrorLogger struct {
	BaseLogger
}

func (e *ErrorLogger) PrintLog(level int, msg string) {
	if level == Error {
		fmt.Printf("[ERROR] %s\n", msg)
	} else {
		fmt.Printf("ignore [ERROR]\n")
		e.next.PrintLog(level, msg)
	}
}

func main() {
	errorLogger := &ErrorLogger{}
	infoLogger := &InfoLogger{}
	debugLogger := &DebugLogger{}

	// error -> info > debug
	infoLogger.SetNext(debugLogger)
	errorLogger.SetNext(infoLogger)
	// 开始调用
	errorLogger.PrintLog(Info, "info")
	errorLogger.PrintLog(Debug, "debug")
}

// 1.分层处理
// 2.当层没处理完成的，放下一层继续处理
// 3.可以把各个处理函数(handler)串成一个链