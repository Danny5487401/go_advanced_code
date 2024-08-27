<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [context](#context)
  - [应用场景](#%E5%BA%94%E7%94%A8%E5%9C%BA%E6%99%AF)
  - [Context 的调用](#context-%E7%9A%84%E8%B0%83%E7%94%A8)
  - [遵循规则](#%E9%81%B5%E5%BE%AA%E8%A7%84%E5%88%99)
  - [Context 源码分析](#context-%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [1. Context接口](#1-context%E6%8E%A5%E5%8F%A3)
    - [2. canceler 接口](#2-canceler-%E6%8E%A5%E5%8F%A3)
    - [实现的类型：](#%E5%AE%9E%E7%8E%B0%E7%9A%84%E7%B1%BB%E5%9E%8B)
      - [1. emptyCtx：即空context，也是所有子context的祖先](#1-emptyctx%E5%8D%B3%E7%A9%BAcontext%E4%B9%9F%E6%98%AF%E6%89%80%E6%9C%89%E5%AD%90context%E7%9A%84%E7%A5%96%E5%85%88)
      - [2. cancelCtx](#2-cancelctx)
      - [3. timerCtx：继承自 cancelCtx 带取消功能的 Context](#3-timerctx%E7%BB%A7%E6%89%BF%E8%87%AA-cancelctx-%E5%B8%A6%E5%8F%96%E6%B6%88%E5%8A%9F%E8%83%BD%E7%9A%84-context)
      - [4. valueCtx：只能携带一个键值对，且自身要依附在上一级 Context 上](#4-valuectx%E5%8F%AA%E8%83%BD%E6%90%BA%E5%B8%A6%E4%B8%80%E4%B8%AA%E9%94%AE%E5%80%BC%E5%AF%B9%E4%B8%94%E8%87%AA%E8%BA%AB%E8%A6%81%E4%BE%9D%E9%99%84%E5%9C%A8%E4%B8%8A%E4%B8%80%E7%BA%A7-context-%E4%B8%8A)
        - [使用场景](#%E4%BD%BF%E7%94%A8%E5%9C%BA%E6%99%AF)
        - [使用注意](#%E4%BD%BF%E7%94%A8%E6%B3%A8%E6%84%8F)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# context

Context 为同一任务的多个 goroutine 之间提供了 *退出信号通知* 和 *元数据传递* 的功能。

## 应用场景
![](.introduction_images/multi_goroutine.png)

如果不用 Context，就不能在 Go 语言里实现多个 goroutine  间的信号通知和元数据传递了吗？答案是：简单场景下可以，在多层级 goroutine 的控制中就行不通了。

假如主协程中有多个任务，主协程对这些任务有超时控制；而其中任务1又有多个子任务，任务1对这些子任务也有自己的超时控制，那么这些子任务既要感知主协程的取消信号，也需要感知任务1的取消信号。

任务的 goroutine 层级越深，想要自己做退出信号感知和元数据共享就越难。

在 Go http 包的 Server 中，每一个请求在都有一个对应的goroutine去处理。请求处理函数通常会启动额外的goroutine用来访问后端服务，比如数据库和 RPC 服务。
用来处理一个请求的goroutine通常需要访问一些与请求特定的数据， 比如终端用户的身份认证信息、验证相关的 token、请求的截止时间。
当一个请求被取消或超时时，所有用来处理该请求的goroutine都应该迅速退出， 然后系统才能释放这些goroutine占用的资源.

所以我们需要一种优雅的方案来实现这样一种机制：

* 上层任务取消后，所有的下层任务都会被取消；
* 中间某一层的任务取消后，只会将当前任务的下层任务取消，而不会影响上层的任务以及同级任务；
* 可以线程安全地在 goroutine 之间共享元数据

## Context 的调用

应该是链式的，通过WithCancel，WithDeadline，WithTimeout或WithValue派生出新的 Context。当父 Context 被取消时，其派生的所有 Context 都将取消。
通过context.WithXXX都将返回新的 Context 和 CancelFunc。调用 CancelFunc 将取消子代，移除父代对子代的引用，并且停止所有定时器。
未能调用 CancelFunc 将泄漏子代，直到父代被取消或定时器触发。go vet工具检查所有流程控制路径上使用 CancelFuncs

## 遵循规则

遵循以下规则，以保持包之间的接口一致，并启用静态分析工具以检查上下文传播。

不要将 Contexts 放入结构体，相反context应该作为第一个参数传入，命名为ctx。
```go
func DoSomething(ctx context.Context，arg Arg)error { // ... user ctx ... }
```

**即使函数允许，也不要传入nil的 Context。如果不知道用哪种 Context，可以使用context.TODO()。**

使用context的Value相关方法只应该用于在程序和接口中传递的和请求相关的元数据，不要用它来传递一些可选的参数
相同的 Context 可以传递给在不同的goroutine；Context 是并发安全的。

## Context 源码分析
![](.introduction_images/context_relation.png)

### 1. Context接口
```go
type Context interface {

    Deadline() (deadline time.Time, ok bool)   // 返回截止时间和ok,
    // 截止时间没有设置，ok为false. 连续调用的结果相同.
    
    Done() <-chan struct{}    // 当 context 被取消或者到了 deadline，返回一个被关闭的 channel
    
    Err() error  // 在 channel Done 关闭后，返回 context 取消原因
    
    Value(key interface{}) interface{}   // 返回值
}
// Context 是一个接口，定义了 4 个方法，它们都是 幂等的。也就是说连续多次调用同一个方法，得到的结果都是相同的

```

上面可以看到Context是一个接口，想要使用就得实现其方法。在context包内部已经为我们实现好了两个空的Context，可以通过调用Background()和TODO()方法获取。
一般的将它们作为Context的根，往下派生。

1. Deadline方法需要返回当前Context被取消的时间，也就是完成工作的截止时间（deadline）；
2. Done方法需要返回一个Channel，这个Channel会在当前工作完成或者上下文被取消之后关闭，多次调用Done方法会返回同一个Channel；
3. Err方法会返回当前Context结束的原因，它只会在Done返回的Channel被关闭时才会返回非空的值；
   - 如果当前Context被取消就会返回Canceled错误；
   - 如果当前Context超时就会返回DeadlineExceeded错误；
```go
var Canceled = errors.New("context canceled")
var DeadlineExceeded error = deadlineExceededError{}

type deadlineExceededError struct{}

func (deadlineExceededError) Error() string   { return "context deadline exceeded" }
func (deadlineExceededError) Timeout() bool   { return true }
func (deadlineExceededError) Temporary() bool { return true }
```
4. Value方法会从Context中返回键对应的值，对于同一个上下文来说，多次调用Value 并传入相同的Key会返回相同的结果，该方法仅用于传递跨API和进程间跟请求域的数据


### 2. canceler 接口
```go

type canceler interface {
    cancel(removeFromParent bool, err error)
    
    Done() <-chan struct{}
}
```

源码中有两个类型实现了 canceler 接口：*cancelCtx 和 *timerCtx。注意是加了 * 号的，是这两个结构体的指针实现了 canceler 接口。

Context 接口设计成这个样子的原因：

1. “取消”操作应该是建议性，而非强制性

caller 不应该去关心、干涉 callee 的情况，决定如何以及何时 return 是 callee 的责任。caller 只需发送“取消”信息，callee 根据收到的信息来做进一步的决策，因此接口并没有定义 cancel 方法。

2. “取消”操作应该可传递

“取消”某个函数时，和它相关联的其他函数也应该“取消”。因此， Done() 方法返回一个只读的 channel，所有相关函数监听此 channel。
一旦 channel 关闭，通过 channel 的“广播机制”，所有监听者都能收到。





### 实现的类型：

- emptyCtx：默认初始的context使用的类型，仅实现Context接口，不做任何处理，返回默认空值。
- cancelCtx：实现了canceler接口，为context提供了可取消自身和子孙的功能。
- timerCtx：在cancelCtx的基础上，对带有定时功能的Context进行了实现。
- valueCtx：对应携带K-V值的context接口实现，携带k-v数据成员，实现了value函数的具体操作。


#### 1. emptyCtx：即空context，也是所有子context的祖先
```go
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
//这个实现只用于在包内定义两个内部实例，并提供对外访问函数。
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
```

- background 通常用在 main 函数中，作为所有 context 的根节点。

- todo 通常用在并不知道传递什么 context的情形。例如，调用一个需要传递 context 参数的函数，你手头并没有其他 context 可以传递，这时就可以传递 todo。这常常发生在重构进行中，给一些函数添加了一个 Context 参数，但不知道要传什么，就用 todo “占个位子”，最终要换成其他 context.

c设计思想：

1. 不需要再对父Context是否为空作为额外的判断，优化了代码结构，在调用时逻辑也更通顺。

2. Go语言不支持继承，而内嵌一个匿名成员，实际上达到了继承的效果，在后面可以看到，因为以一个完全实现了context接口的emptyCtx实例为起点，
    cancelCtx等实现已经继承了默认的函数，只需要再实现需要用到的函数即可，缺失的其他函数一定会被最底层的emptyCtx实例提供


#### 2. cancelCtx
![](.introduction_images/cancelCtx.png)

withCancel使用分析：WithCancel 以一个新的 Done channel 返回一个父 Context 的拷贝。
```go
// Canceled is the error returned by Context.Err when the context is canceled.
var Canceled = errors.New("context canceled")

func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
    c := newCancelCtx(parent)
    propagateCancel(parent, &c)
    return &c, func() { c.cancel(true, Canceled) }
}
```


```go

type cancelCtx struct {
    Context
    mu       sync.Mutex            // 用于保证Context的线程安全，
    done     chan struct{}         // 互斥锁通道实例done向本goroutine外发送本context已经被关闭的的消息
    children map[canceler]struct{} // children字典则存储了本context派生的所有context，key值为canceler类型.
    err      error                 //  err字段用于标记该context是否已经被取消，取消则将是非空值
}

type canceler interface {
    cancel(removeFromParent bool, err error)
    Done() <-chan struct{}
}
```

```go
// newCancelCtx returns an initialized cancelCtx.
func newCancelCtx(parent Context) cancelCtx {
    return cancelCtx{Context: parent}
}

func propagateCancel(parent Context, child canceler) {
    if parent.Done() == nil {
    	// 祖先为不可取消类型，则自己就是取消链的根，直接返回
    	// Note：emptyCtx的Done() 为Nil
        return // parent is never canceled
    }
    // 通过辅助函数parentCancelCtx向上回溯，尝试找到最近的*cancelCtx类型祖先。
    if p, ok := parentCancelCtx(parent); ok {
        p.mu.Lock()
        if p.err != nil {
            // parent has already been canceled
            child.cancel(false, p.err)
        } else {
        	//则将自身加入其children列表中。
            if p.children == nil {
                p.children = make(map[canceler]struct{})
            }
            p.children[child] = struct{}{}
        }
        p.mu.Unlock()
    } else {
    	// 单独监听其父context和自己的取消情况。
        go func() {
            select {
            case <-parent.Done():
                child.cancel(false, parent.Err())
            case <-child.Done():
            }
        }()
    }
}


func parentCancelCtx(parent Context) (*cancelCtx, bool) {
	for {
		switch c := parent.(type) {
		// 总共就三种Ctx
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


```
propagateCancel完成的主要工作：
将新建立的cancelCtx，绑定到祖先的取消广播树中， 简单来说，就是将自身存储到最近的*cancelCtx类型祖先的children列表中，接收该祖先的广播。

1、parent.Done()==nil，祖先为不可取消类型，则自己就是取消链的根，直接返回。

2、通过辅助函数parentCancelCtx向上回溯，尝试找到最近的*cancelCtx类型祖先。

3、如果成功找到，则先判断是否已经被取消，如果为否，则将自身加入其children列表中。

4、否则，就单独监听其父context和自己的取消情况。


分析：

得益于子context引用父context的设计，对于每个contest都将可以通过向上回溯得到一条引用链，
辅助函数 parentCancelCtx即通过不断向内部引用类型转换，达到回看context历史的目的，寻找最近的*cancelCtx型祖先

结束时调用
```go
// cancel closes c.done, cancels each of c's children, and, if
// removeFromParent is true, removes c from its parent's children.
func (c *cancelCtx) cancel(removeFromParent bool, err error) {
	if err == nil {
		// 调用时需要传入错误
		panic("context: internal error: missing cancel error")
	}
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return // already canceled
	}
	c.err = err
	if c.done == nil {
		c.done = closedchan
	} else {
		// 1. 将本context自身关闭
		close(c.done)
	}
	for child := range c.children {
		// NOTE: acquiring the child's lock while holding parent's lock.
		// 2. 向其在取消链中的子节点广播取消通知
		child.cancel(false, err)
	}
	c.children = nil
	c.mu.Unlock()

	if removeFromParent {
		// 解除自身在取消链中的绑定
		removeChild(c.Context, c)
	}
}
```
```go
func removeChild(parent Context, child canceler) {
	p, ok := parentCancelCtx(parent)
	if !ok {
		return
	}
	p.mu.Lock()
	if p.children != nil {
		delete(p.children, child)
	}
	p.mu.Unlock()
}
```


#### 3. timerCtx：继承自 cancelCtx 带取消功能的 Context
```go
type timerCtx struct {
    cancelCtx
    timer *time.Timer // Under cancelCtx.mu.
    deadline time.Time
}
```
```go
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if cur, ok := parent.Deadline(); ok && cur.Before(d) {
		// 时间还在范围之内
		// The current deadline is already sooner than the new one.
		return WithCancel(parent)
	}
	c := &timerCtx{
		cancelCtx: newCancelCtx(parent),
		deadline:  d,
	}
	propagateCancel(parent, c)
	dur := time.Until(d)
	if dur <= 0 {
		// 已经超时
		c.cancel(true, DeadlineExceeded) // deadline has already passed
		return c, func() { c.cancel(false, Canceled) }
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil {
		// 建立定时任务
		c.timer = time.AfterFunc(dur, func() {
			c.cancel(true, DeadlineExceeded)
		})
	}
	return c, func() { c.cancel(true, Canceled) }
}
```
本质上说，无论是WithDeadline还是WithTimeout生成的子context，都是对于time.AfterFunc函数和cancelCtx的一个封装。

#### 4. valueCtx：只能携带一个键值对，且自身要依附在上一级 Context 上
![](.introduction_images/value_context.png)
```go

type valueCtx struct {
    Context
    key, val interface{} // valueCtx:在原状态基础上添加一个键值对
}
	
func (c *valueCtx) Value(key interface{}) interface{} {
    if c.key == key {
        return c.val
    }
    // valueCtx类型真正实现了value函数，该函数是一个向上递归的查询过程，如果key不存在，将递归调用emptyCtx定义好的默认函数，返回一个nil值
    return c.Context.Value(key)  
}

```

##### 使用场景
WithValue函数能够将请求作用域的数据与 Context 对象建立关系。

每次要在Context链路上增加要携带的键值对时，都要在上级Context的基础上新建一个 valueCtx 存储键值对，切只能增加不能修改，
读取 Context 上的键值又是一个幂等的操作，所以 Context 就这样实现了线程安全的数据共享机制，且全程无锁，不会影响性能

##### 使用注意
提供的键必须是可比性和应该不是字符串类型或任何其他内置的类型以避免包使用的上下文之间的碰撞。WithValue 用户应该定义自己的键的类型。
为了避免在分配给interface{}时进行分配，上下文键通常具有具体类型struct{}。或者，导出的上下文关键变量的静态类型应该是指针或接口

使用 context.WithValue
```go

func WithValue(parent Context, key, val interface{}) Context {
    if key == nil {
        panic("nil key")
    }
    if !reflectlite.TypeOf(key).Comparable() {
        panic("key is not comparable")
    }
    return &valueCtx{parent, key, val}
}

```
valueCtx类型真正实现了value函数，该函数是一个向上递归的查询过程，如果key不存在，将递归调用emptyCtx定义好的默认函数，返回一个nil值






