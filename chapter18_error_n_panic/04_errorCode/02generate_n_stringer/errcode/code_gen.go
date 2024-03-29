// Code generated by "stringer -type ErrCode -linecomment -output code_gen.go"; DO NOT EDIT.

package errcode

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OK-200]
	_ = x[ServerCommonError-1000000]
	_ = x[ServerInvalidParams-1000000]
	_ = x[ServerTimout-1000000]
	_ = x[TicketNotExit-3000000]
	_ = x[TicketStatusNotOK-3000001]
	_ = x[TicketUpdateFail-3000002]
	_ = x[BookNotFoundError-4000000]
	_ = x[BookHasBeenBorrowedError-4000001]
}

const (
	_ErrCode_name_0 = "成功返回"
	_ErrCode_name_1 = "OK"
	_ErrCode_name_2 = "工单不存在工单状态不合理工单更新失败"
	_ErrCode_name_3 = "书未找到书已经被借走了"
)

var (
	_ErrCode_index_2 = [...]uint8{0, 15, 36, 54}
	_ErrCode_index_3 = [...]uint8{0, 12, 33}
)

func (i ErrCode) String() string {
	switch {
	case i == 200:
		return _ErrCode_name_0
	case i == 1000000:
		return _ErrCode_name_1
	case 3000000 <= i && i <= 3000002:
		i -= 3000000
		return _ErrCode_name_2[_ErrCode_index_2[i]:_ErrCode_index_2[i+1]]
	case 4000000 <= i && i <= 4000001:
		i -= 4000000
		return _ErrCode_name_3[_ErrCode_index_3[i]:_ErrCode_index_3[i+1]]
	default:
		return "ErrCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
