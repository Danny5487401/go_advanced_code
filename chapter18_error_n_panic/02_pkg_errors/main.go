package main

import (
	E "errors"
	"fmt"
	"github.com/pkg/errors"
)

/*
背景：
	func AuthenticateRequest(r *Request) error {
		 err := authenticate(r.User)
		 if err != nil {
			return fmt.Errorf("authenticate failed: %v", err)
		 }
		 return nil
	}
	这种做法实际上是先错误转换成字符串，再拼接另一个字符串，最后，再通过 fmt.Errorf 转换成错误。这样做破坏了相等性检测，即我们无法判断错误是否是一种预先定义好的错误了。
方案：
	github.com/pkg/errors。提供了友好的界面

%s,%v //功能一样，输出错误信息，不包含堆栈

%q //输出的错误信息带引号，不包含堆栈

%+v //输出错误信息和堆栈
*/

func main() {
	err0 := t1()
	fmt.Printf("原始方式:%+v\n", err0.Error())
	// 1。对错误消息进一步封装
	err1 := errors.Wrap(err0, "附加信息")
	if err1 != nil {
		//打印错误需要%+v才能详细输出
		fmt.Printf("err :%+v\n", err1)
	}

	fmt.Println("-------------")
	// 2。对消息打印错误栈
	err2 := errors.WithStack(err0)
	if err2 != nil {
		//打印错误需要%+v才能详细输出
		fmt.Printf("Stackerr :%+v\n", err2)
	}

	fmt.Println("-------------")
	// 3.找根本原因
	err3 := errors.Cause(err0)
	if err3 != nil {
		//打印错误需要%+v才能详细输出
		fmt.Printf("Causeerr :%+v\n", err3)
	}
}

func t1() error {
	e := E.New("err0")

	return e
}
