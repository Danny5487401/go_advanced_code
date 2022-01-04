package main

import (
	"fmt"
	"github.com/pkg/errors"
)

func wrapNewPointerError() error {
	// 	Go1.13版本为fmt.Errorf函数新加了一个%w占位符用来生成一个可以包裹Error的Wrapping Error。
	return fmt.Errorf("wrap err0:%w", fmt.Errorf("i am a error0"))
}

func wrapConstantPointerError() error {
	return fmt.Errorf("wrap err0:%w", constantErr)
}

func newPointerErrorWithoutWrap() error {
	return fmt.Errorf("wrap err0:%v", fmt.Errorf("i am a error0"))
}
func NewError() error {
	return errors.New("i am a error0")
}

var constantErr = fmt.Errorf("i am a error0")

func main() {
	//fmt.Println("第一个结果", errors.Is(wrapNewPointerError(), fmt.Errorf("i am a error0"))) //
	fmt.Println("第二个结果", errors.Is(wrapConstantPointerError(), constantErr)) //
	//fmt.Println("第三个结果", errors.Is(newPointerErrorWithoutWrap(), fmt.Errorf("i am a error0"))) //
	fmt.Println("第四个结果", errors.Is(NewError(), errors.New("i am a error0")))

}
