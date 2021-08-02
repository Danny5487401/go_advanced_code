package main

import (
	"fmt"
	"go_advenced_code/chapter18_error/04_errorCode/01traditional/errcode"
)

func main() {
	code := errcode.ERR_CODE_INVALID_PARAMS
	fmt.Println(code.String())

	// 输出: 1 无效参数
}
