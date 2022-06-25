package main

import (
	"fmt"
	"go_advanced_code/chapter18_error_n_panic/04_errorCode/02generate_n_stringer/errcode"
)

func main() {
	code := errcode.ERR_CODE_TIMEOUT
	fmt.Println(code) // 输出: 1 无效参数

}
