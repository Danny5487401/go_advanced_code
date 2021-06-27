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

// 函数选项模式
type GoPoolOption func(*GoPool)

// 初始化选项
func WithMaxLimit(max int) GoPoolOption {
	return func(gp *GoPool) {
		gp.MaxLimit = max
		gp.tokenChan = make(chan struct{}, gp.MaxLimit) //令牌数等于最大并发量

		for i := 0; i < gp.MaxLimit; i++ {
			gp.tokenChan <- struct{}{} // 把通道晒满
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
	// 等待取出一个令牌
	token := <-gp.tokenChan // if there are no tokens, we'll block here

	go func() {
		fn()
		// 模拟并发执行时间
		time.Sleep(time.Second * 3)
		// 归还令牌
		gp.tokenChan <- token
	}()
}

// Wait will wait all the tasks executed, and then return
// 等到所有chan取出来后关闭通道
func (gp *GoPool) Wait() {
	for i := 0; i < gp.MaxLimit; i++ {
		<-gp.tokenChan
	}

	close(gp.tokenChan)
}

// 返回令牌个数
func (gp *GoPool) size() int {
	return len(gp.tokenChan)
}

/*使用
1. Submit 在令牌不足时，会阻塞当前调用(因此Go runtime会执行其他不阻塞的代码)
2. Wait() 会等到回收所有令牌之后，才返回
*/

func main() {
	// 初始化线程池
	goPool := NewGoPool(WithMaxLimit(2))
	defer goPool.Wait()

	//goPool.Submit(func() {//你的代码})
	goPool.Submit(
		func() { fmt.Println(1) })
	goPool.Submit(
		func() { fmt.Println(2) })
	goPool.Submit(
		func() { fmt.Println(3) })
	goPool.Submit(
		func() { fmt.Println(4) })
	goPool.Submit(
		func() { fmt.Println(5) })

}

/* 运行结果：打印无序
5
1
3
4
2


*/
