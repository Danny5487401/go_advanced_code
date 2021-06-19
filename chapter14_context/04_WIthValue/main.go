package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)


/*
context.WithValue

使用场景
	 WithValue函数能够将请求作用域的数据与 Context 对象建立关系。
使用注意：
	提供的键必须是可比性和应该不是字符串类型或任何其他内置的类型以避免包使用的上下文之间的碰撞。WithValue 用户应该定义自己的键的类型。
	为了避免在分配给interface{}时进行分配，上下文键通常具有具体类型struct{}。或者，导出的上下文关键变量的静态类型应该是指针或接口
源码分析:
	type valueCtx struct {
		Context
		key, val interface{}
	}
	valueCtx:在原状态基础上添加一个键值对
	func (c *valueCtx) Value(key interface{}) interface{} {
		if c.key == key {
			return c.val
		}
		return c.Context.Value(key)
	}
	valueCtx类型真正实现了value函数，该函数是一个向上递归的查询过程，如果key不存在，将递归调用emptyCtx定义好的默认函数，返回一个nil值


	func WithValue(parent Context, key, val interface{}) Context {
		if key == nil {
			panic("nil key")
		}
		if !reflectlite.TypeOf(key).Comparable() {
			panic("key is not comparable")
		}
		return &valueCtx{parent, key, val}
	}
 */


type TraceCode string

var wg sync.WaitGroup

func worker(ctx context.Context) {
	key := TraceCode("TRACE_CODE")
	traceCode, ok := ctx.Value(key).(string) // 在子goroutine中获取trace code
	if !ok {
		fmt.Println("invalid trace code")
	}
LOOP:
	for {
		fmt.Printf("worker, trace code:%s\n", traceCode)
		time.Sleep(time.Millisecond * 10) // 假设正常连接数据库耗时10毫秒
		select {
		case <-ctx.Done(): // 50毫秒后自动调用
			break LOOP
		default:
		}
	}
	fmt.Println("worker done!")
	wg.Done()
}

func main() {
	// 设置一个50毫秒的超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	// 在系统的入口中设置trace code传递给后续启动的goroutine实现日志数据聚合
	ctx = context.WithValue(ctx, TraceCode("TRACE_CODE"), "12512312234") // ctx多次封装
	wg.Add(1)
	go worker(ctx)
	time.Sleep(time.Second * 5)
	cancel() // 通知子goroutine结束
	wg.Wait()
	fmt.Println("over")
}
