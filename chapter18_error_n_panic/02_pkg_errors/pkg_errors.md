# error
error 是源代码内嵌的接口类型。根据导出原则，只有大写的才能被其它源码包引用，但是 error 属于 predeclared identifiers 预定义的，并不是关键字

```go
// The error built-in interface type is the conventional interface for
// representing an error condition, with the nil value representing no error.
type error interface {
	Error() string
}
```
error 只有一个方法 Error() string 返回错误消息

```go
// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
```
## 背景：
```go
func AuthenticateRequest(r *Request) error {
    err := authenticate(r.User)
    if err != nil {
        return fmt.Errorf("authenticate failed: %v", err)
    }
    return nil
}
```

这种做法实际上是先错误转换成字符串，再拼接另一个字符串，最后，再通过 fmt.Errorf 转换成错误。
这样做破坏了相等性检测，即我们无法判断错误是否是一种预先定义好的错误了。

当前 error 的问题有两点：

- 无法 wrap 更多的信息，比如调用栈，比如层层封装的 error 消息
- 无法很好的处理类型信息，比如我想知道错误是 io 类型的，还是 net 类型的

%s,%v //功能一样，输出错误信息，不包含堆栈

%q //输出的错误信息带引号，不包含堆栈

%+v //输出错误信息和堆栈

## 扩展包Pkg.errors
### 1. Wrap 更多的消息
- Wrap 封装底层 error, 增加更多消息，提供调用栈信息，这是原生 error 缺少的
- WithMessage 封装底层 error, 增加更多消息，但不提供调用栈信息
- Cause 返回最底层的 error, 剥去层层的 wrap

```go
import (
   "database/sql"
   "fmt"

   "github.com/pkg/errors"
)

func foo() error {
   return errors.Wrap(sql.ErrNoRows, "foo failed")
}

func bar() error {
   return errors.WithMessage(foo(), "bar failed")
}

func main() {
   err := bar()
   if errors.Cause(err) == sql.ErrNoRows {
      fmt.Printf("data not found, %v\n", err)
      fmt.Printf("%+v\n", err)
      return
   }
   if err != nil {
      // unknown error
   }
}
/*Output:
data not found, bar failed: foo failed: sql: no rows in result set
sql: no rows in result set
foo failed
main.foo
    /usr/three/main.go:11
main.bar
    /usr/three/main.go:15
main.main
    /usr/three/main.go:19
runtime.main
    ...
*/
```

源码：
```go
// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   message,
	}
	return &withStack{
		err,
		callers(),
	}
}
```
```go
type withStack struct {
	error
	*stack
}

func (w *withStack) Cause() error { return w.error }

func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}
```
Cause 递归调用，如果没有实现 causer 接口，那么就返回这个 err


### 2. 对比
errors.Is 会递归的 Unwrap err, 判断错误是不是 sql.ErrNoRows，这里个小问题，Is 是做的指针地址判断，如果错误 Error() 内容一样，但是根 error 是不同实例，那么 Is 判断也是 false, 这点就很扯
```go
import (
   "database/sql"
   "errors"
   "fmt"
)

func bar() error {
   if err := foo(); err != nil {
      return fmt.Errorf("bar failed: %w", foo())
   }
   return nil
}

func foo() error {
   return fmt.Errorf("foo failed: %w", sql.ErrNoRows)
}

func main() {
   err := bar()
   if errors.Is(err, sql.ErrNoRows) {
      fmt.Printf("data not found,  %+v\n", err)
      return
   }
   if err != nil {
      // unknown error
   }
}
/* Outputs:
data not found,  bar failed: foo failed: sql: no rows in result set
*/
```
