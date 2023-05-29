<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [ants: goroutine 池](#ants-goroutine-%E6%B1%A0)
  - [两种使用方式](#%E4%B8%A4%E7%A7%8D%E4%BD%BF%E7%94%A8%E6%96%B9%E5%BC%8F)
    - [1. 第一种方式：NewPool](#1-%E7%AC%AC%E4%B8%80%E7%A7%8D%E6%96%B9%E5%BC%8Fnewpool)
      - [管理工人需要实现的接口](#%E7%AE%A1%E7%90%86%E5%B7%A5%E4%BA%BA%E9%9C%80%E8%A6%81%E5%AE%9E%E7%8E%B0%E7%9A%84%E6%8E%A5%E5%8F%A3)
  - [选项Options](#%E9%80%89%E9%A1%B9options)
  - [第二种方式：poolWithFunc结构体](#%E7%AC%AC%E4%BA%8C%E7%A7%8D%E6%96%B9%E5%BC%8Fpoolwithfunc%E7%BB%93%E6%9E%84%E4%BD%93)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# ants: goroutine 池

## 两种使用方式
p, _ := ants.NewPool(cap)：这种方式创建的池子对象需要调用p.Submit(task)提交任务，任务是一个无参数无返回值的函数；

p, _ := ants.NewPoolWithFunc(cap, func(interface{}))：这种方式创建的池子对象需要指定池函数，并且使用p.Invoke(arg)调用池函数。arg就是传给池函数func(interface{})的参数。

### 1. 第一种方式：NewPool
结构体
```go
type Pool struct {
    // 容量，负数为不限制
    capacity int32
    
    // 当前运行的工作数量
    running int32
    
    // lock for protecting the worker queue.
    lock sync.Locker
    
    // 存放一组 worker 对象，workerArray只是一个接口，表示一个 worker 容器
    workers workerArray
    
    // state is used to notice the pool to closed itself.
    state int32
    
    // cond for waiting to get a idle worker.
    cond *sync.Cond
    
    // 使用sync.Pool对象池管理和创建worker对象，提升性能；
    workerCache sync.Pool
    
    // pool.Submit时已经阻塞的数目
    blockingNum int
    
    options *Options
}
```

真正执行的工人
```go
type goWorker struct {
	// 持有 goroutine 池的引用
	pool *Pool

	// task is a job should be done.
	task chan func()

	// recycleTime will be update when putting a worker back into queue.
	recycleTime time.Time
}
```

运行时
```go
func (w *goWorker) run() {
    go func() {
        for f := range w.task {
            //f == nil为 true 时return，也是一个细节点
          if f == nil {
            return
          }
          f()
          if ok := w.pool.revertWorker(w); !ok {
            //还有一个细节，如果放回操作失败，则会调用return，这会让 goroutine 运行结束，防止 goroutine 泄漏
            return
          }
        }
    }()
}
```

#### 管理工人需要实现的接口
```go
type workerArray interface {
    len() int  //worker 数量
    isEmpty() bool  //worker 数量是否为 0
    insert(worker *goWorker) error  // goroutine 任务执行结束后，将相应的 worker 放回workerArray中
    detach() *goWorker //从workerArray中取出一个 worker
    retrieveExpiry(duration time.Duration) []*goWorker //取出所有的过期 worker
    reset() //重置容器
} 
```

workerArray在ants中有两种实现，即workerStack和loopQueue

1. 没有预分配时：workerStack
```go
type workerStack struct {
	items  []*goWorker
	expiry []*goWorker
	size   int
}
```
具体实现:获取worker符合栈后进先出
```go

func (wq *workerStack) insert(worker *goWorker) error {
	wq.items = append(wq.items, worker)
	return nil
}

func (wq *workerStack) detach() *goWorker {
	l := wq.len()
	if l == 0 {
		return nil
	}

	w := wq.items[l-1]
	wq.items[l-1] = nil // 避免内存泄漏
	wq.items = wq.items[:l-1]

	return w
}

```
2. 有预分配时：loopQueue
```go
type loopQueue struct {
	items  []*goWorker
	expiry []*goWorker
	head   int
	tail   int
	size   int
	isFull bool
}
```
具体实现
```go
func (wq *loopQueue) insert(worker *goWorker) error {
	if wq.size == 0 {
		return errQueueIsReleased
	}

	if wq.isFull {
		return errQueueIsFull
	}
	wq.items[wq.tail] = worker
	wq.tail++

	if wq.tail == wq.size {
		wq.tail = 0
	}
	if wq.tail == wq.head {
		wq.isFull = true
	}

	return nil
}

func (wq *loopQueue) detach() *goWorker {
	if wq.isEmpty() {
		return nil
	}

	w := wq.items[wq.head]
	wq.items[wq.head] = nil
	wq.head++
	if wq.head == wq.size {
		wq.head = 0
	}
	wq.isFull = false

	return w
}
```


初始化流程
```go
func NewPool(size int, options ...Option) (*Pool, error) {
    p := &Pool{
        capacity: int32(size),
        lock:     internal.NewSpinLock(), // 自旋锁
        options:  opts,
    }
	
    // 根据预分配生成不同的worker队列
	//使用预分配时，创建loopQueueType类型的结构，反之创建stackType类型
    if p.options.PreAlloc {
        if size == -1 {
            return nil, ErrInvalidPreAllocSize
        }
        p.workers = newWorkerArray(loopQueueType, size)
    } else {
            p.workers = newWorkerArray(stackType, 0)
        }

    //最后启动一个 goroutine 用于定期清理过期的 worker
    go p.purgePeriodically()
}
```


提交任务
```go
func (p *Pool) Submit(task func()) error {
  if p.IsClosed() {
    return ErrPoolClosed
  }
  var w *goWorker
  if w = p.retrieveWorker(); w == nil {
    return ErrPoolOverload
  }
  w.task <- task
  return nil
}
```


## 选项Options
```go
// src/github.com/panjf2000/ants/options.go
type Options struct {
    //表示 goroutine 空闲多长时间之后会被ants池回收
    ExpiryDuration time.Duration
    //调用NewPool()/NewPoolWithFunc()之后预分配worker（管理一个工作 goroutine 的结构体）切片
    PreAlloc bool
    //最大阻塞任务数量。即池中 goroutine 数量已到池容量，且所有 goroutine 都处理繁忙状态，这时到来的任务会在阻塞列表等待。这个选项设置的是列表的最大长度。阻塞的任务数量达到这个值后，后续任务提交直接返回失败
    MaxBlockingTasks int
    
    //池是否阻塞，默认阻塞。提交任务时，如果ants池中 goroutine 已到上限且全部繁忙，阻塞的池会将任务添加的阻塞列表等待（当然受限于阻塞列表长度，见上一个选项）。非阻塞的池直接返回失败
    Nonblocking bool
    PanicHandler func(interface{})
    Logger Logger
}
```



## 第二种方式：poolWithFunc结构体
```go
type PoolWithFunc struct {
	// 容量
	capacity int32

	// running is the number of the currently running goroutines.
	running int32

	// lock for protecting the worker queue.
	lock sync.Locker

	// 真正执行工作的worker
	workers []*goWorkerWithFunc

	// state is used to notice the pool to closed itself.
	state int32

	// cond for waiting to get a idle worker.
	cond *sync.Cond

	// 需要执行的工作
	poolFunc func(interface{})

	// 对象复用池
	workerCache sync.Pool

	// blockingNum is the number of the goroutines already been blocked on pool.Submit, protected by pool.lock
	blockingNum int

	options *Options
}
```