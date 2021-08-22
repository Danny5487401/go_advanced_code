package main

/*
sync 的扩展包errgroup
需求：
	一般在golang 中想要并发运行业务时会直接开goroutine，关键字go ,但是直接go的话函数是无法对返回数据进行处理error的。
解决方案：
## 初级版本：
	一般是直接在出错的地方打入log日志,将出的错误记录到日志文件中，也可以集合日志收集系统直接将该错误用邮箱或者办公软件发送给你如：钉钉机器人+graylog.

## 中级版本
	当然你也可以自己在log包里封装好可以接受channel。
	利用channel通道，将go中出现的error传入到封装好的带有channel接受器的log包中，进行错误收集或者通知通道接受return出来即可

## 终极版本
	errgroup
*/

import (
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func main() {
	group := new(errgroup.Group)
	nums := []int{-1, 0, 1}

	for _, num := range nums {
		tempNum := num // 子协程中若直接访问num，则可能是同一个变量，所以要用临时变量

		// 子协程
		group.Go(func() error {
			if tempNum < 0 {
				return errors.New("tempNum < 0 ")
			}
			fmt.Println("tempNum:", tempNum)
			return nil
		})
	}

	// 捕获err
	if err := group.Wait(); err != nil {
		fmt.Println("Get errors: ", err)
	} else {
		fmt.Println("Get all num successfully!")
	}
}
