package _0_emptyCtx

/*
emptyCtx：即空context，也是所有子context的祖先
源码：
	type emptyCtx int

	func (*emptyCtx) Deadline() (deadline time.Time, ok bool) {
		return
	}

	func (*emptyCtx) Done() <-chan struct{} {
		return nil
	}

	func (*emptyCtx) Err() error {
		return nil
	}

	func (*emptyCtx) Value(key interface{}) interface{} {
		return nil
	}

	func (e *emptyCtx) String() string {
		switch e {
		case background:
			return "context.Background"
		case todo:
			return "context.TODO"
		}
		return "unknown empty Context"
	}
这个实现只用于在包内定义两个内部实例，并提供对外访问函数。
	var (
		background = new(emptyCtx)
		todo       = new(emptyCtx)
	)

	func Background() Context {
		return background
	}
	func TODO() Context {
		return todo
	}
设计思想：
	1、不需要再对父Context是否为空作为额外的判断，优化了代码结构，在调用时逻辑也更通顺。

	2、Go语言不支持继承，而内嵌一个匿名成员，实际上达到了继承的效果，在后面可以看到，因为以一个完全实现了context接口的emptyCtx实例为起点，
		cancelCtx等实现已经继承了默认的函数，只需要再实现需要用到的函数即可，缺失的其他函数一定会被最底层的emptyCtx实例提供

 */
