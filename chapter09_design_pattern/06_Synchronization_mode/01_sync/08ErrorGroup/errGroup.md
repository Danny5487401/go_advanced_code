# sync 的扩展包errgroup

## 需求：
	一般在golang 中想要并发运行业务时会直接开goroutine，关键字go,但是直接go的话函数是无法对返回数据进行处理error的。
## 解决方案

### 初级版本：
	一般是直接在出错的地方打入log日志,将出的错误记录到日志文件中，也可以集合日志收集系统直接将该错误用邮箱或者办公软件发送给你如：钉钉机器人+graylog.

### 中级版本
	当然你也可以自己在log包里封装好可以接受channel。
	利用channel通道，将go中出现的error传入到封装好的带有channel接受器的log包中，进行错误收集或者通知通道接受return出来即可

### 终极版本
	errgroup

## 源码分析
```go
// 结构体
type Group struct {
    cancel  func()             // 这个存的是context的cancel方法
    wg      sync.WaitGroup  // 封装sync.WaitGroup
    errOnce sync.Once          // 保证只接受一次错误
    err     error             // 保存第一个返回的错误
}

func WithContext(ctx context.Context) (*Group, context.Context){
    ctx, cancel := context.WithCancel(ctx)
    return &Group{cancel: cancel}, ctx
}

func (g *Group) Go(f func() error) {
	// 增加一个计数器
    g.wg.Add(1)

    go func() {
		// 控制是否结束
        defer g.wg.Done()
        if err := f(); err != nil {
			// 如果有一个函数f运行出错了，我们把它保存起来
            g.errOnce.Do(func() {
				// 这里的目的就是保证获取到第一个出错的信息，避免被后面的Goroutine的错误覆盖
                g.err = err             //记录子协程中的错误
                if g.cancel != nil {
                    g.cancel()
                }
            })
        }
    }()
}

// errgroup 可以捕获和记录子协程的错误(只能记录最先出错的协程的错误
// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
    g.wg.Wait()
    if g.cancel != nil {
        g.cancel()
    }
    return g.err
}

// errgroup 可以控制协程并发顺序。确保子协程执行完成后再执行主协程
```