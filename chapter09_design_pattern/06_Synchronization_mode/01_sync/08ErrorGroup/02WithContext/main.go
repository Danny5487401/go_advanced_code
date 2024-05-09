package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background()) // 外部主动执行cancel
	// 1。首先传递 context 初始化 errgroup 对象
	group, errCtx := errgroup.WithContext(ctx)

	for index := 0; index < 3; index++ {
		indexTemp := index // 子协程中若直接访问index，则可能是同一个变量，所以要用临时变量
		// 2。每一个 group.Go() 都会新启一个协程, Go()函数接受一个 func() error 函数类型
		// 新建子协程
		group.Go(func() error {
			fmt.Printf("indexTemp=%d \n", indexTemp)
			if indexTemp == 0 {
				fmt.Println("indexTemp == 0 start ")
				fmt.Println("indexTemp == 0 end")
			} else if indexTemp == 1 {
				fmt.Println("indexTemp == 1 start")
				//这里一般都是某个协程发生异常之后，调用cancel()
				//这样别的协程就可以通过errCtx获取到err信息，以便决定是否需要取消后续操作
				cancel()
				fmt.Println("indexTemp == 1 err ")
			} else if indexTemp == 2 {
				fmt.Println("indexTemp == 2 begin")

				// 休眠1秒，用于捕获子协程2的出错
				time.Sleep(1 * time.Second)

				//检查 其他协程已经发生错误，如果已经发生异常，则不再执行下面的代码
				err := checkGoroutineErr(errCtx)
				if err != nil {
					return err
				}
				fmt.Println("indexTemp == 2 end ")
			}
			return nil
		})
	}

	// 	3。使用 Wait()方法阻塞主协程,直到所有子协程执行完成
	// 捕获err
	err := group.Wait()
	if err == nil {
		fmt.Println("都完成了")
	} else {
		fmt.Printf("get error:%v", err)
	}
}

// 校验是否有协程已发生错误
func checkGoroutineErr(errContext context.Context) error {
	select {
	// errgroup 可以使用 context 实现协程撤销。或者超时撤销。子协程中使用 ctx.Done()来获取撤销信号
	case <-errContext.Done():
		return errContext.Err()
	default:
		return nil
	}
}
