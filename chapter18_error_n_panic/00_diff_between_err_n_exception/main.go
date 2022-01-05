package main

import (
	"fmt"
	"strconv"
)

// 错误和异常处理

// go语言认为错误就要自己处理，

func div(a, b int) (int, error) {
	if b == 0 {
		panic("被除数不能为零")
	}
	return a / b, nil
}

func main() {
	// 1. 错误就是能遇到可能出现的情况，这些情况可能导致你的代码出问题，例如参数检查，数据库转换
	/* go 认为这个itoa函数不可能出错，没有必要返回error，
	内部代码出错这个时候应该抛出异常Panic,对应python中的raise,Java中的throw
	*/

	i, err := strconv.Atoi("12") //返回error
	if err != nil {
		fmt.Println(err.Error())
		return

	}
	fmt.Println(i)

	// 2. 异常情况
	panicSituation()

}

/*
异常情况 ：异常处理的作用域（场景）：
	1. 空指针引用
	2. 下标越界
	3. 除数为0
	4. 不应该出现的分支，比如default
	5. 输入不应该引起函数错误
*/
func panicSituation() {

	// 抛出异常和异常捕捉
	a := 12
	b := 0
	// 不想函数被停止,需定义一个函数
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("异常捕获到:%v\n", err)
		}
		fmt.Println("恢复")
	}()
	res, err := div(a, b)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res)
}
