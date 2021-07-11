package _1_customized_error

import "runtime"

/*
自定义错误
	接口
	type error interface {
		Error() string
	}
*/
// 1。简单版
func New(text string) error {
	return &errorString{text}
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

/*
以上存在问题：
	Go语言对错误的设计非常简洁，但是对于我们开发者来说，很明显是不足的，比如我们需要知道出错的更多信息，在什么文件的，哪一行代码？
	只有这样我们才更容易的定位问题
*/

// 2。复杂版本 ：以下可参考源码pkg/error
type stack []uintptr
type errorStackString struct {
	s string
	*stack
}

func (e errorStackString) Error() string {
	return e.s
}

// 存储堆栈信息的stack字段，我们在生成错误的时候，就可以把调用的堆栈信息存储在这个字段里
func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

func NewStack(text string) error {
	return &errorStackString{
		s:     text,
		stack: callers(),
	}
}

// 3. 封装错误的消息
type withMessage struct {
	cause error
	msg   string
}

func (w withMessage) Error() string {
	return w.msg

}

func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   message,
	}
}
