package _2generate_n_stringer

import (
	"fmt"
	"github.com/Danny5487401/go_advanced_code/chapter18_error_n_panic/04_errorCode/02generate_n_stringer/errcode"
	"github.com/pkg/errors"
	"testing"
)

func TestText(t *testing.T) {
	books := []string{
		"Book1",
		"Book222222",
		"Book3333333333",
		"Book333333333344444444",
	}

	for _, bookName := range books {
		err := searchBook(bookName)

		// 特殊业务场景：如果发现书被借走了，下次再来就行了，不需要作为错误处理
		if err != nil {
			// 提取error这个interface底层的错误码，一般在API的返回前才提取
			// As - 获取错误的具体实现
			var customErr = new(CustomError)
			// As - 解析错误内容
			if errors.As(err, &customErr) {
				//fmt.Printf("AS中的信息：当前书为: %s ,error Code is %d, message is %s\n", bookName, customErr.Code, customErr.Message)
				if customErr.Code == errcode.BookHasBeenBorrowedError {
					fmt.Printf("IS中的info信息：%s 已经被借走了, 只需按Info处理!\n", bookName)
				} else {
					// 如果已有堆栈信息，应调用WithMessage方法
					newErr := errors.WithMessage(err, "WithMessage err1")
					// 使用%+v可以打印完整的堆栈信息
					fmt.Printf("IS中的error信息：%s 未找到，应该按Error处理! ,newErr is: %+v\n", bookName, newErr)
				}
			}
		}
	}
}

func searchBook(bookName string) error {
	// 1 发现图书馆不存在这本书 - 认为是错误，需要打印详细的错误信息
	if len(bookName) > 15 {
		return NewErrCode(errcode.BookNotFoundError)
	} else if len(bookName) > 10 {
		// 2 发现书被借走了 - 打印一下被接走的提示即可，不认为是错误
		return NewErrCodeMsg(errcode.BookHasBeenBorrowedError, "借书失败")
	} else if len(bookName) > 6 {
		return NewErrMsg("呵呵")
	}
	// 3 找到书 - 不需要任何处理
	return nil
}
