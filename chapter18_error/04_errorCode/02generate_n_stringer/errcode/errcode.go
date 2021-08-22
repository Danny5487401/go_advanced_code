//go:generate stringer -type ErrCode -linecomment -output code_string.go

/*
注意点
	go:generate前面只能使用//注释，注释必须在行首，前面不能有空格且//与go:generate之间不能有空格！！！
	go:generate可以在任何 Go 源文件中，最好在类型定义的地方
*/

package errcode

type ErrCode int

const (
	ERR_CODE_OK             ErrCode = 0 // OK
	ERR_CODE_INVALID_PARAMS ErrCode = 1 // 无效参数
	ERR_CODE_TIMEOUT        ErrCode = 2 // 超时
)
