# singleFlight 

## 作用
在一个服务中抑制对下游的多次重复请求。一个比较常见的使用场景是：我们在使用 Redis 对数据库中的数据进行缓存，发生缓存击穿时，大量的流量都会打到数据库上进而影响服务的尾延时

## 缓存击穿
缓存在某个时间点过期的时候，恰好在这个时间点对这个Key有大量的并发请求过来，这些请求发现缓存过期一般都会从后端DB加载数据并回设到缓存，这个时候大并发的请求可能会瞬间把后端DB压垮。


## 主要使用方法
Do方法
```go
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool)
```
- shared：表示是否有其他协程得到了这个结果v。
- key：同一个key，同时只有一个协程执行。
- fn：被包装的函数。
- v：返回值，即执行的结果。其他等待的协程都会拿到。

```go
func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result
```
与Do方法一样，只是返回的是一个channel，执行结果会发送到channel中，其他等待的协程都可以从channel中拿到结果
## 源码分析

Group结构：代表一类工作，同一个group中，同样的key同时只能被执行一次。
```go
type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call // lazily initialized  保存key对应的函数执行过程和结果的变量。
}
```
Group的结构非常简单，一个锁来保证并发安全，另一个map用来保存key对应的函数执行过程和结果的变量。

```go
// call is an in-flight or completed singleflight.Do call
type call struct {
	wg sync.WaitGroup  // //用WaitGroup实现只有一个协程执行函数

	// These fields are written once before the WaitGroup is done
	// and are only read after the WaitGroup is done.
	val interface{} // 函数执行结果
	err error

	// forgotten indicates whether Forget was called with this call's key
	// while the call was still in flight.
	forgotten bool

	// These fields are read and written with the singleflight
	// mutex held before the WaitGroup is done, and are read but
	// not written after the WaitGroup is done.
	dups  int  //含义是duplications，即同时执行同一个key的协程数量
	chans []chan<- Result
}

```

## Do方法
```go
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	// 写Group的m字段时，加锁保证写安全。
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
        //如果key已经存在，说明已经有协程在执行，则dups++，并等待其执行完毕后，返回其执行结果，执行结果保存在对应的call的val字段里
		c.dups++
		g.mu.Unlock()
		c.wg.Wait()

		if e, ok := c.err.(*panicError); ok {
			panic(e)
		} else if c.err == errGoexit {
			runtime.Goexit()
		}
		return c.val, c.err, true
	}
	// 如果key不存在，则新建一个call，并使用WaitGroup来阻塞其他协程，同时在m字段里写入key和对应的call
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	// 第一个进来的协程来执行这个函数
	g.doCall(c, key, fn)
	return c.val, c.err, c.dups > 0
}
```

## DoChan方法
```go
func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result {
	ch := make(chan Result, 1)
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		// 可以看到，每个等待的协程，都有一个结果channel。
		// 从之前的g.doCall里也可以看到，每个channel都给塞了结果。
		// 为什么不所有协程共用一个channel？因为那样就得在channel里塞至少与协程数量一样的结果数量，但是你却无法保证用户一个协程只读取一次。
		c.chans = append(c.chans, ch)
		g.mu.Unlock()
		return ch
	}
	c := &call{chans: []chan<- Result{ch}}
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	go g.doCall(c, key, fn)

	return ch
}
```

## 具体执行流程
```go
func (g *Group) doCall(c *call, key string, fn func() (interface{}, error)) {
	normalReturn := false
	recovered := false

	// use double-defer to distinguish panic from runtime.Goexit,
	// more details see https://golang.org/cl/134395
	defer func() {
		// the given function invoked runtime.Goexit
		if !normalReturn && !recovered {
			c.err = errGoexit
		}
        //执行完毕后，就可以通知其他协程可以拿结果了
		c.wg.Done()
		g.mu.Lock()
		defer g.mu.Unlock()
		if !c.forgotten {
			// 其实这里是为了保证执行完毕之后，对应的key被删除，Group有一个方法Forget（key string），可以用来主动删除key，
			// 这里是判断那个方法是否被调用过，被调用过则字段forgotten会置为true，如果没有被调用过，则在这里把key删除
			delete(g.m, key)
		}

		if e, ok := c.err.(*panicError); ok {
			// In order to prevent the waiting channels from being blocked forever,
			// needs to ensure that this panic cannot be recovered.
			if len(c.chans) > 0 {
				go panic(e)
				select {} // Keep this goroutine around so that it will appear in the crash dump.
			} else {
				panic(e)
			}
		} else if c.err == errGoexit {
			// Already in the process of goexit, no need to call again
		} else {
			// Normal return
			for _, ch := range c.chans {
				// 将执行结果发送到channel里，这里是给DoChan方法使用的
				ch <- Result{c.val, c.err, c.dups > 0}
			}
		}
	}()

	func() {
		defer func() {
			if !normalReturn {
				// Ideally, we would wait to take a stack trace until we've determined
				// whether this is a panic or a runtime.Goexit.
				//
				// Unfortunately, the only way we can distinguish the two is to see
				// whether the recover stopped the goroutine from terminating, and by
				// the time we know that, the part of the stack trace relevant to the
				// panic has been discarded.
				if r := recover(); r != nil {
					c.err = newPanicError(r)
				}
			}
		}()
        //执行被包装的函数
		c.val, c.err = fn()
		normalReturn = true
	}()

	if !normalReturn {
		recovered = true
	}
}
```