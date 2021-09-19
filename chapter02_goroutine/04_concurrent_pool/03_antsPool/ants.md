# ants: goroutine 池

##选项Options
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

##NewPool()部分源码
```go
if p.options.PreAlloc {
  if size == -1 {
    return nil, ErrInvalidPreAllocSize
  }
  p.workers = newWorkerArray(loopQueueType, size)
} else {
  p.workers = newWorkerArray(stackType, 0)
}

//使用预分配时，创建loopQueueType类型的结构，反之创建stackType类型
```