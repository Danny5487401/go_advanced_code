#sync 的扩展包errgroup

##需求：

	一般在golang 中想要并发运行业务时会直接开goroutine，关键字go ,但是直接go的话函数是无法对返回数据进行处理error的。
##解决方案

### 初级版本：
	一般是直接在出错的地方打入log日志,将出的错误记录到日志文件中，也可以集合日志收集系统直接将该错误用邮箱或者办公软件发送给你如：钉钉机器人+graylog.

### 中级版本
	当然你也可以自己在log包里封装好可以接受channel。
	利用channel通道，将go中出现的error传入到封装好的带有channel接受器的log包中，进行错误收集或者通知通道接受return出来即可

### 终极版本
	errgroup

##源码分析
```go



// 结构体
type Group struct {
  cancel  func()             //context cancel()
    wg      sync.WaitGroup
    errOnce sync.Once          //只会传递第一个出现错的协程的 error
    err     error              //传递子协程错误
}


func (g *Group) Go(f func() error) {
    g.wg.Add(1)

    go func() {
        defer g.wg.Done()
        if err := f(); err != nil {
            g.errOnce.Do(func() {
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