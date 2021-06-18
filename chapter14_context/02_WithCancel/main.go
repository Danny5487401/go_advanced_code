package main

// 简化 对于处理单个请求的多个 goroutine 之间与请求域的数据、取消信号、截止时间等相关操

import (
	"context"
	"fmt"
)

/*
withCancel
源码:
	type cancelCtx struct {
		Context

		mu       sync.Mutex            // protects following fields
		done     chan struct{}         // created lazily, closed by first cancel call
		children map[canceler]struct{} // set to nil by the first cancel call
		err      error                 // set to non-nil by the first cancel call
	}

	type canceler interface {
		cancel(removeFromParent bool, err error)
		Done() <-chan struct{}
	}
分析：
	在cancelCtx的结构定义中，包含了互斥锁用于保证Context的线程安全，通道实例done向本runtion外发送本context已经被关闭的的消息，
		err字段用于标记该context是否已经被取消，取消则将是非空值，children字典则存储了本context派生的所有context，key值为canceler类型

*/

func main() {
	// gen generates integers in a separate goroutine and
	// sends them to the returned channel.
	// The callers of gen need to cancel the context once
	// they are done consuming generated integers not to leak
	// the internal goroutine started by gen.

	/*
		gen函数在单独的goroutine中生成整数并将它们发送到返回的通道。 gen的调用者在使用生成的整数之后需要取消上下文，以免gen启动的内部goroutine发生泄漏。
	*/

	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					return //  return结束该goroutine，防止泄露
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 当我们取完需要的整数后调用cancel

	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
}
/*
使用分析：

	func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
		c := newCancelCtx(parent)
		propagateCancel(parent, &c)
		return &c, func() { c.cancel(true, Canceled) }
	}

	func propagateCancel(parent Context, child canceler) {
	if parent.Done() == nil {
		return // parent is never canceled
	}
	if p, ok := parentCancelCtx(parent); ok {
		p.mu.Lock()
		if p.err != nil {
			// parent has already been canceled
			child.cancel(false, p.err)
		} else {
			if p.children == nil {
				p.children = make(map[canceler]struct{})
			}
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
	} else {
		go func() {
			select {
			case <-parent.Done():
				child.cancel(false, parent.Err())
			case <-child.Done():
			}
		}()
	}
}


propagateCancel完成的主要工作：
	将新建立的cancelCtx，绑定到祖先的取消广播树中，简单来说，就是将自身存储到最近的*cancelCtx类型祖先的chindren列表中，接收该祖先的广播。

	1、parent.Done()==nil，祖先为不可取消类型，则自己就是取消链的根，直接返回。

	2、通过辅助函数parentCancelCtx向上回溯，尝试找到最近的*cancelCtx类型祖先。

	3、如果成功找到，则先判断是否已经被取消，如果为否，则将自身加入其children列表中。

	4、否则，就单独监听其父context和自己的取消情况。
func parentCancelCtx(parent Context) (*cancelCtx, bool) {
	for {
		switch c := parent.(type) {
		case *cancelCtx:
			return c, true
		case *timerCtx:
			return &c.cancelCtx, true
		case *valueCtx:
			parent = c.Context
		default:
			return nil, false
		}
	}
}
分析：
	得益于子context引用父context的设计，对于每个contest都将可以通过向上回溯得到一条引用链，
	辅助函数 parentCancelCtx即通过不断向内部引用类型转换，达到回看context历史的目的，寻找最近的*cancelCtx型祖先

 */
