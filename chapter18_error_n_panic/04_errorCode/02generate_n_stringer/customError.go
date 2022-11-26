package _2generate_n_stringer

import (
	"fmt"
	"github.com/Danny5487401/go_advanced_code/chapter18_error_n_panic/04_errorCode/02generate_n_stringer/errcode"
)

type Xerror interface {
	GetErrCode() uint32
	GetErrMsg() string
	error
}

// CustomError 1、错误时返回自定义结构r结构体，并重写Error()方法
type CustomError struct {
	Code    errcode.ErrCode // 业务码
	message string          // 业务消息
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.Code, e.message)
}

// GetErrCode 返回给前端的错误码
func (e *CustomError) GetErrCode() uint32 {
	return uint32(e.Code)
}

// GetErrMsg 返回给前端显示端错误信息
func (e *CustomError) GetErrMsg() string {
	return e.message
}

func NewErrCodeMsg(errCode errcode.ErrCode, errMsg string) Xerror {
	return &CustomError{Code: errCode, message: errMsg}
}
func NewErrCode(errCode errcode.ErrCode) Xerror {
	return &CustomError{Code: errCode, message: errCode.String()}
}

func NewErrMsg(errMsg string) Xerror {
	return &CustomError{Code: errcode.ServerCommonError, message: errMsg}
}
