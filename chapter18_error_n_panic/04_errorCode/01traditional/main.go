package main

import (
	"fmt"
	"github.com/Danny5487401/go_advanced_code/chapter18_error_n_panic/04_errorCode/01traditional/errcode"
)

func main() {
	code := errcode.ERR_CODE_INVALID_PARAMS
	fmt.Println(code.String())

	// 输出: 1 无效参数
}
