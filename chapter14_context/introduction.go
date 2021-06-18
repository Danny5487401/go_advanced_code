package chapter14_context

/*
context理解：
	结合linux操作系统的cpu上下文切换/子进程好父进程进行理解

应用场景：
	在 Go http 包的 Server 中，每一个请求在都有一个对应的goroutine去处理。
	请求处理函数通常会启动额外的goroutine用来访问后端服务，比如数据库和 RPC 服务。用来处理一个请求的goroutine通常需要访问一些与请求特定的数据，
	比如终端用户的身份认证信息、验证相关的 token、请求的截止时间。当一个请求被取消或超时时，所有用来处理该请求的goroutine都应该迅速退出
	，然后系统才能释放这些goroutine占用的资源

Context 的调用
	应该是链式的，通过WithCancel，WithDeadline，WithTimeout或WithValue派生出新的 Context。当父 Context 被取消时，其派生的所有 Context 都将取消。
	通过context.WithXXX都将返回新的 Context 和 CancelFunc。调用 CancelFunc 将取消子代，移除父代对子代的引用，并且停止所有定时器。
		未能调用 CancelFunc 将泄漏子代，直到父代被取消或定时器触发。go vet工具检查所有流程控制路径上使用 CancelFuncs

遵循规则
	遵循以下规则，以保持包之间的接口一致，并启用静态分析工具以检查上下文传播。

	不要将 Contexts 放入结构体，相反context应该作为第一个参数传入，命名为ctx。 func DoSomething（ctx context.Context，arg Arg）error { // ... use ctx ... }
	即使函数允许，也不要传入nil的 Context。如果不知道用哪种 Context，可以使用context.TODO()。
	使用context的Value相关方法只应该用于在程序和接口中传递的和请求相关的元数据，不要用它来传递一些可选的参数
	相同的 Context 可以传递给在不同的goroutine；Context 是并发安全的。
Context 原理
	总结：mutex和channel的结合

源码
	type Context interface {
		// Done returns a channel that is closed when this Context is canceled
		// or times out.
		Done() <-chan struct{}

		// Err indicates why this context was canceled, after the Done channel
		// is closed.
		Err() error

		// Deadline returns the time when this Context will be canceled, if any.
		Deadline() (deadline time.Time, ok bool)

		// Value returns the value associated with key or nil if none.
		Value(key interface{}) interface{}
	}
	Done()，返回一个channel。当times out或者调用cancel方法时，将会close掉。
	Err()，返回一个错误。该context为什么被取消掉。
	Deadline()，返回截止时间和ok。
	Value()，返回值。

上面可以看到Context是一个接口，想要使用就得实现其方法。在context包内部已经为我们实现好了两个空的Context，可以通过调用Background()和TODO()方法获取。
	一般的将它们作为Context的根，往下派生。

	1. Deadline方法需要返回当前Context被取消的时间，也就是完成工作的截止时间（deadline）；
	2. Done方法需要返回一个Channel，这个Channel会在当前工作完成或者上下文被取消之后关闭，多次调用Done方法会返回同一个Channel；
	3. Err方法会返回当前Context结束的原因，它只会在Done返回的Channel被关闭时才会返回非空的值；
		如果当前Context被取消就会返回Canceled错误；
		如果当前Context超时就会返回DeadlineExceeded错误；
	4. Value方法会从Context中返回键对应的值，对于同一个上下文来说，多次调用Value 并传入相同的Key会返回相同的结果，该方法仅用于传递跨API和进程间跟请求域的数据
实现的类型：
	emptyCtx：默认初始的context使用的类型，仅实现Context接口，不做任何处理，返回默认空值。

	valueCtx：对应携带K-V值的context接口实现，携带k-v数据成员，实现了value函数的具体操作。

	cancelCtx：实现了cancler接口，为context提供了可取消自身和子孙的功能。

	timerCtx：在cancelCtx的基础上，对带有定时功能的Context进行了实现。

*/
