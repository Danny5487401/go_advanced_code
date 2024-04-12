package errcode

import "fmt"

// 传统方式
// 1. 定义错误码

type ErrCode int

const (
	ERR_CODE_OK             ErrCode = 0 // OK
	ERR_CODE_INVALID_PARAMS ErrCode = 1 // 无效参数
	ERR_CODE_TIMEOUT        ErrCode = 2 // 超时
)

// 2. 定义错误码与描述信息的映射
var mapErrDesc = map[ErrCode]string{
	ERR_CODE_OK:             "OK",
	ERR_CODE_INVALID_PARAMS: "无效参数",
	ERR_CODE_TIMEOUT:        "超时",
	// ...
}

// 3.  根据错误码返回描述信息
// 不可导出 使用string()导出
func getDescription(errCode ErrCode) string {
	if desc, exist := mapErrDesc[errCode]; exist {
		return desc
	}

	return fmt.Sprintf("error code: %d", errCode)
}

func (e ErrCode) String() string {
	return getDescription(e)
}

/*
缺点
	每次增加错误码时，都需要修改mapErrDesc，有时候可能会忘。另外，错误描述在注释和mapErrDesc都出现了
*/
