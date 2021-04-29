package main

import (
	"fmt"
	"time"
)

// 并发控制

// 场景：同时不超过2个goroutine 执行任务

type GoPool struct {
	// 最大并发量
	MaxLimit int

	tokenChan chan struct{}
}

type GoPoolOption func(*GoPool)

func WithMaxLimit(max int) GoPoolOption {
	return func(gp *GoPool) {
		gp.MaxLimit = max
		gp.tokenChan = make(chan struct{}, gp.MaxLimit)

		for i := 0; i < gp.MaxLimit; i++ {
			gp.tokenChan <- struct{}{}
		}
	}
}

func NewGoPool(options ...GoPoolOption) *GoPool {
	p := &GoPool{}
	for _, o := range options {
		o(p)
	}

	return p
}

// Submit will wait a token, and then execute fn
func (gp *GoPool) Submit(fn func()) {
	token := <-gp.tokenChan // if there are no tokens, we'll block here

	go func() {
		fn()
		// 模拟并发执行时间
		time.Sleep(time.Second*3)

		gp.tokenChan <- token
	}()
}

// Wait will wait all the tasks executed, and then return
func (gp *GoPool) Wait() {
	for i := 0; i < gp.MaxLimit; i++ {
		<-gp.tokenChan
	}

	close(gp.tokenChan)
}

func (gp *GoPool) size() int {
	return len(gp.tokenChan)
}


/*使用
1. Submit 在令牌不足时，会阻塞当前调用(因此Go runtime会执行其他不阻塞的代码)
2. Wait() 会等到回收所有令牌之后，才返回
 */

func main()  {
	goPool := NewGoPool(WithMaxLimit(2))
	defer goPool.Wait()

	//goPool.Submit(func() {//你的代码})
	goPool.Submit(
		func() {fmt.Println(1)})
	goPool.Submit(
		func() {fmt.Println(2)})
	goPool.Submit(
		func() {fmt.Println(3)})
	goPool.Submit(
		func() {fmt.Println(4)})
	goPool.Submit(
		func() {fmt.Println(5)})

}
/* 运行结果：打印无须
5
1
3
4
2


 */
