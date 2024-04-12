//go:generate stringer -type ErrCode -linecomment -output code_gen.go

package errcode

type ErrCode uint32

/**(前3位代表业务,后三位代表具体功能)**/

const (
	OK ErrCode = 200 // 成功返回
)

//全局错误码

const (
	ServerCommonError   ErrCode = 1e6 // OK
	ServerInvalidParams               // 无效参数
	ServerTimout                      // 超时
)

// 工单模块300开始
const (
	TicketNotExit     ErrCode = iota + 3e6 // 工单不存在
	TicketStatusNotOK                      // 工单状态不合理
	TicketUpdateFail                       // 工单更新失败
)

// 读书模块400开始
const (
	BookNotFoundError        ErrCode = iota + 4e6 // 书未找到
	BookHasBeenBorrowedError                      // 书已经被借走了
)
