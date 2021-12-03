# fasthttp
高性能主要源自于“复用”，通过服务协程和内存变量的复用，节省了大量资源分配的成本。

## 源码分析
优点
- 工作协程复用
- 内存变量复用： 内部大量使用了sync.Pool

### 初始化
```go
func (s *Server) Serve(ln net.Listener) error {
    // ....

	wp := &workerPool{
		WorkerFunc:      s.serveConn,  // 具体的工作逻辑
		MaxWorkersCount: maxWorkersCount,
		LogAllErrors:    s.LogAllErrors,
		Logger:          s.logger(),
		connState:       s.setState,
	}
	wp.Start()

	// Count our waiting to accept a connection as an open connection.
	// This way we can't get into any weird state where just after accepting
	// a connection Shutdown is called which reads open as 0 because it isn't
	// incremented yet.
	atomic.AddInt32(&s.open, 1)
	defer atomic.AddInt32(&s.open, -1)

	for {
		if c, err = acceptConn(s, ln, &lastPerIPErrorTime); err != nil {
			wp.Stop()
			if err == io.EOF {
				return nil
			}
			return err
		}
		s.setState(c, StateNew)
		atomic.AddInt32(&s.open, 1)
		// 开始服务
		if !wp.Serve(c) {
			atomic.AddInt32(&s.open, -1)
			s.writeFastError(c, StatusServiceUnavailable,
				"The connection cannot be served because Server.Concurrency limit exceeded")
			c.Close()
			s.setState(c, StateClosed)
            // ...
		}
		c = nil
	}
}

func (wp *workerPool) Serve(c net.Conn) bool {
	// 获取工作的channel
    ch := wp.getCh()
    if ch == nil {
        return false
    }
    ch.ch <- c
    return true
}
```
### 具体获取过程: 工作协程的复用
```go
func (wp *workerPool) getCh() *workerChan {
	var ch *workerChan
	createWorker := false

	wp.lock.Lock()
	ready := wp.ready
	n := len(ready) - 1
	if n < 0 {
		if wp.workersCount < wp.MaxWorkersCount {
			createWorker = true
			wp.workersCount++
		}
	} else {
		// 从准备好的worker池中获取
		ch = ready[n]
		ready[n] = nil
		wp.ready = ready[:n]
	}
	wp.lock.Unlock()

	if ch == nil {
		if !createWorker {
			return nil
		}
		vch := wp.workerChanPool.Get()
		ch = vch.(*workerChan)
		go func() {
			wp.workerFunc(ch)
			wp.workerChanPool.Put(vch)
		}()
	}
	return ch
}

```

工作池结构池 workerPool
```go 
// /Users/python/go/pkg/mod/github.com/valyala/fasthttp@v1.29.0/workerpool.go
type workerPool struct {
	// Function for serving server connections.
	// It must leave c unclosed.
	WorkerFunc ServeHandler

	MaxWorkersCount int

	LogAllErrors bool

	MaxIdleWorkerDuration time.Duration

	Logger Logger

	lock         sync.Mutex
	workersCount int
	mustStop     bool

	ready []*workerChan

	stopCh chan struct{}

	workerChanPool sync.Pool

	connState func(net.Conn, ConnState)
}
```

每个工作协程
```go
func (wp *workerPool) workerFunc(ch *workerChan) {
	var c net.Conn

	var err error
	for c = range ch.ch {
		if c == nil {
			break
		}
        // 业务执行流程
		if err = wp.WorkerFunc(c); err != nil && err != errHijacked {
			errStr := err.Error()
			if wp.LogAllErrors || !(strings.Contains(errStr, "broken pipe") ||
				strings.Contains(errStr, "reset by peer") ||
				strings.Contains(errStr, "request headers: small read buffer") ||
				strings.Contains(errStr, "unexpected EOF") ||
				strings.Contains(errStr, "i/o timeout")) {
				wp.Logger.Printf("error when serving connection %q<->%q: %s", c.LocalAddr(), c.RemoteAddr(), err)
			}
		}
		if err == errHijacked {
			wp.connState(c, StateHijacked)
		} else {
			_ = c.Close()
			wp.connState(c, StateClosed)
		}
		c = nil
        
		if !wp.release(ch) {
			break
		}
	}

	wp.lock.Lock()
	wp.workersCount--
	wp.lock.Unlock()
}
```

### 内存变量复用
具体工作时
```go
func (s *Server) serveConn(c net.Conn) (err error) {
    // ....
    // 获取上下文
	ctx := s.acquireCtx(c)
	// ....
	if ctx != nil {
		// 释放上下文
        s.releaseCtx(ctx)
    }
	return
}
```
获取及释放上下文
```go
func (s *Server) acquireCtx(c net.Conn) (ctx *RequestCtx) {
	v := s.ctxPool.Get()
	if v == nil {
		ctx = &RequestCtx{
			s: s,
		}
		keepBodyBuffer := !s.ReduceMemoryUsage
		ctx.Request.keepBodyBuffer = keepBodyBuffer
		ctx.Response.keepBodyBuffer = keepBodyBuffer
	} else {
		ctx = v.(*RequestCtx)
	}
	ctx.c = c
	return
}
func (s *Server) releaseCtx(ctx *RequestCtx) {
    if ctx.timeoutResponse != nil {
        panic("BUG: cannot release timed out RequestCtx")
    }
    ctx.c = nil
    ctx.remoteAddr = nil
    ctx.fbr.c = nil
    ctx.userValues.Reset()
    s.ctxPool.Put(ctx)
}

```
