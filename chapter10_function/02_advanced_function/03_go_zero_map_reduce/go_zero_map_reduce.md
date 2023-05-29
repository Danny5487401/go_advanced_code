<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [go-zero map-reduce](#go-zero-map-reduce)
  - [分类](#%E5%88%86%E7%B1%BB)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# go-zero map-reduce


## 分类
面向用户的方法比较多，方法主要分为两大类：

1. 无返回
- 执行过程发生错误立即终止
- 执行过程不关注错误

2. 有返回值
- 手动写入 source，手动读取聚合数据 channel
- 手动写入 source，自动读取聚合数据 channel

## 源码分析

```go
const (
	defaultWorkers = 16
	minWorkers     = 1
)

var (
	// ErrCancelWithNil is an error that mapreduce was canceled with nil.
	ErrCancelWithNil = errors.New("mapreduce canceled with nil")
	// ErrReduceNoOutput is an error that reduce did not output a value.
	ErrReduceNoOutput = errors.New("reduce not writing value")
)

type (
	// ForEachFunc is used to do element processing, but no output.
	ForEachFunc func(item interface{})
	// GenerateFunc is used to let callers send elements into source.
	// 1. 数据生产func
	// source - 数据被生产后写入source
	GenerateFunc func(source chan<- interface{})
	// MapFunc is used to do element processing and write the output to writer.
	MapFunc func(item interface{}, writer Writer)
	// MapperFunc is used to do element processing and write the output to writer,
	// use cancel func to cancel the processing.
	// 2. 数据加工func
	// item - 生产出来的数据
	// writer - 调用writer.Write()可以将加工后的向后传递至reducer
	// cancel - 终止流程func
	MapperFunc func(item interface{}, writer Writer, cancel func(error))
	// ReducerFunc is used to reduce all the mapping output and write to writer,
	// use cancel func to cancel the processing.
	// 3. 数据聚合func
	// pipe - 加工出来的数据
	// writer - 调用writer.Write()可以将聚合后的数据返回给用户
	// cancel - 终止流程func
	ReducerFunc func(pipe <-chan interface{}, writer Writer, cancel func(error))
	// VoidReducerFunc is used to reduce all the mapping output, but no output.
	// Use cancel func to cancel the processing.
	VoidReducerFunc func(pipe <-chan interface{}, cancel func(error))
	// Option defines the method to customize the mapreduce.
	Option func(opts *mapReduceOptions)

	mapperContext struct {
		ctx       context.Context
		mapper    MapFunc
		source    <-chan interface{}
		panicChan *onceChan
		collector chan<- interface{}
		doneChan  <-chan lang.PlaceholderType
		workers   int
	}

	mapReduceOptions struct {
		ctx     context.Context
		workers int
	}

	// Writer interface wraps Write method.
	Writer interface {
		Write(v interface{})
	}
)

// Finish runs fns parallelly, canceled on any error.
// 1. 并发执行func，发生任何错误将会立即终止流程
func Finish(fns ...func() error) error {
	if len(fns) == 0 {
		return nil
	}

	return MapReduceVoid(func(source chan<- interface{}) {
		for _, fn := range fns {
			source <- fn
		}
	}, func(item interface{}, writer Writer, cancel func(error)) {
		fn := item.(func() error)
		if err := fn(); err != nil {
			cancel(err)
		}
	}, func(pipe <-chan interface{}, cancel func(error)) {
	}, WithWorkers(len(fns)))
}

// FinishVoid runs fns parallelly.
// 2.并发执行func，即使发生错误也不会终止流程
func FinishVoid(fns ...func()) {
	if len(fns) == 0 {
		return
	}

	ForEach(func(source chan<- interface{}) {
		for _, fn := range fns {
			source <- fn
		}
	}, func(item interface{}) {
		fn := item.(func())
		fn()
	}, WithWorkers(len(fns)))
}

// ForEach maps all elements from given generate but no output.
func ForEach(generate GenerateFunc, mapper ForEachFunc, opts ...Option) {
	options := buildOptions(opts...)
	panicChan := &onceChan{channel: make(chan interface{})}
	source := buildSource(generate, panicChan)
	collector := make(chan interface{}, options.workers)
	done := make(chan lang.PlaceholderType)

	go executeMappers(mapperContext{
		ctx: options.ctx,
		mapper: func(item interface{}, writer Writer) {
			mapper(item)
		},
		source:    source,
		panicChan: panicChan,
		collector: collector,
		doneChan:  done,
		workers:   options.workers,
	})

	for {
		select {
		case v := <-panicChan.channel:
			panic(v)
		case _, ok := <-collector:
			if !ok {
				return
			}
		}
	}
}

// MapReduce maps all elements generated from given generate func,
// and reduces the output elements with given reducer.
func MapReduce(generate GenerateFunc, mapper MapperFunc, reducer ReducerFunc,
	opts ...Option) (interface{}, error) {
	panicChan := &onceChan{channel: make(chan interface{})}
	source := buildSource(generate, panicChan)
	return mapReduceWithPanicChan(source, panicChan, mapper, reducer, opts...)
}

// MapReduceChan maps all elements from source, and reduce the output elements with given reducer.
// 支持传入数据源channel，并返回聚合后的数据
// source - 数据源channel
// mapper - 读取source内容并处理
// reducer - 数据处理完毕发送至reducer聚合
func MapReduceChan(source <-chan interface{}, mapper MapperFunc, reducer ReducerFunc,
	opts ...Option) (interface{}, error) {
	panicChan := &onceChan{channel: make(chan interface{})}
	return mapReduceWithPanicChan(source, panicChan, mapper, reducer, opts...)
}

// MapReduceChan maps all elements from source, and reduce the output elements with given reducer.
func mapReduceWithPanicChan(source <-chan interface{}, panicChan *onceChan, mapper MapperFunc,
	reducer ReducerFunc, opts ...Option) (interface{}, error) {
	// 可选参数设置--指数据加工阶段协程数量
	options := buildOptions(opts...)
	// 聚合数据channel，需要手动调用write方法写入到output中
	output := make(chan interface{})
	defer func() {
		//   如果有多次写入的话则会造成阻塞从而导致协程泄漏
		// 这里用 for range检测是否可以读出数据，读出数据说明多次写入了
		// 为什么这里使用panic呢？显示的提醒用户用法错了会比自动修复掉好一些
		for range output {
			panic("more than one element written in reducer")
		}
	}()

	// collector is used to collect data from mapper, and consume in reducer
	// 创建有缓冲的chan，容量为workers，创建有缓冲的chan，容量为workers
	collector := make(chan interface{}, options.workers)
	// if done is closed, all mappers and reducer should stop processing
	done := make(chan lang.PlaceholderType)
	// 支持阻塞写入chan的writer
	writer := newGuardedWriter(options.ctx, output, done)
	var closeOnce sync.Once
	// use atomic.Value to avoid data race
	var retErr errorx.AtomicError
	// 数据聚合任务已结束，发送完成标志
	finish := func() {
		closeOnce.Do(func() {
			close(done)
			close(output)
		})
	}
	// 取消操作
	cancel := once(func(err error) {
		// 设置error
		if err != nil {
			retErr.Set(err)
		} else {
			retErr.Set(ErrCancelWithNil)
		}
		// 清空source channel
		drain(source)
		finish()
	})

	go func() {
		defer func() {
			// 清空聚合任务channel
			drain(collector)
			if r := recover(); r != nil {
				panicChan.write(r)
			}
			// 正常结束
			finish()
		}()
		// 执行数据加工
		// 注意writer.write将加工后数据写入了output
		reducer(collector, writer, cancel)
	}()

	//  真正从生成器通道取数据执行Mapper
	// 异步执行数据加工
	// source - 数据生产
	// collector - 数据收集
	// done - 结束标志
	// workers - 并发数
	go executeMappers(mapperContext{
		ctx: options.ctx,
		mapper: func(item interface{}, w Writer) {
			mapper(item, w, cancel)
		},
		source:    source,
		panicChan: panicChan,
		collector: collector,
		doneChan:  done,
		workers:   options.workers,
	})

	select {
	case <-options.ctx.Done():
		cancel(context.DeadlineExceeded)
		return nil, context.DeadlineExceeded
	case v := <-panicChan.channel:
		panic(v)
	case v, ok := <-output:
		// reducer将加工后的数据写入了output，
		// 需要数据返回时读取output即可
		if err := retErr.Load(); err != nil {
			return nil, err
		} else if ok {
			return v, nil
		} else {
			return nil, ErrReduceNoOutput
		}
	}
}

// MapReduceVoid maps all elements generated from given generate,
// and reduce the output elements with given reducer.
// 无返回值，关注错误
// GenerateFunc 用于生产数据, MapperFunc 读取生产出的数据，进行处理, 这里表示不对 mapper 后的数据做聚合返回
func MapReduceVoid(generate GenerateFunc, mapper MapperFunc, reducer VoidReducerFunc, opts ...Option) error {
	_, err := MapReduce(generate, mapper, func(input <-chan interface{}, writer Writer, cancel func(error)) {
		reducer(input, cancel)
	}, opts...)
	if errors.Is(err, ErrReduceNoOutput) {
		return nil
	}

	return err
}

// WithContext customizes a mapreduce processing accepts a given ctx.
func WithContext(ctx context.Context) Option {
	return func(opts *mapReduceOptions) {
		opts.ctx = ctx
	}
}

// WithWorkers customizes a mapreduce processing with given workers.
func WithWorkers(workers int) Option {
	return func(opts *mapReduceOptions) {
		if workers < minWorkers {
			opts.workers = minWorkers
		} else {
			opts.workers = workers
		}
	}
}

func buildOptions(opts ...Option) *mapReduceOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

func buildSource(generate GenerateFunc, panicChan *onceChan) chan interface{} {
	source := make(chan interface{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicChan.write(r)
			}
			close(source)
		}()

		generate(source)
	}()

	return source
}

// drain drains the channel.
func drain(channel <-chan interface{}) {
	// drain the channel
	for range channel {
	}
}

// 数据加工
func executeMappers(mCtx mapperContext) {
	var wg sync.WaitGroup
	defer func() {
		// 等待数据加工任务完成
		// 防止数据加工的协程还未处理完数据就直接退出了
		wg.Wait()
		// 关闭数据加工channel
		close(mCtx.collector)
		drain(mCtx.source)
	}()

	var failed int32
	// 根据指定数量创建 worker池,控制数据加工的协程数量
	pool := make(chan lang.PlaceholderType, mCtx.workers)
	// 数据加工writer
	writer := newGuardedWriter(mCtx.ctx, mCtx.collector, mCtx.doneChan)
	for atomic.LoadInt32(&failed) == 0 {
		select {
		case <-mCtx.ctx.Done():
			// 监听到外部结束信号，直接结束
			return
		case <-mCtx.doneChan:
			return

		// 控制数据加工协程数量
		// 缓冲区容量-1
		// 无容量时将会被阻塞，等待释放容量
		case pool <- lang.Placeholder:

			item, ok := <-mCtx.source
			if !ok {

				// 缓冲区容量+1
				<-pool
				//  当通道关闭，结束
				return
			}

			wg.Add(1)
			go func() {
				// 异步执行数据加工，防止panic错误
				defer func() {
					if r := recover(); r != nil {
						atomic.AddInt32(&failed, 1)
						mCtx.panicChan.write(r)
					}
					wg.Done()
					<-pool
				}()

				mCtx.mapper(item, writer)
			}()
		}
	}
}

func newOptions() *mapReduceOptions {
	return &mapReduceOptions{
		ctx:     context.Background(),
		workers: defaultWorkers,
	}
}

func once(fn func(error)) func(error) {
	once := new(sync.Once)
	return func(err error) {
		once.Do(func() {
			fn(err)
		})
	}
}

type guardedWriter struct {
	ctx     context.Context
	channel chan<- interface{}
	done    <-chan lang.PlaceholderType
}

func newGuardedWriter(ctx context.Context, channel chan<- interface{},
	done <-chan lang.PlaceholderType) guardedWriter {
	return guardedWriter{
		ctx:     ctx,
		channel: channel,
		done:    done,
	}
}

func (gw guardedWriter) Write(v interface{}) {
	select {
	case <-gw.ctx.Done():
		return
	case <-gw.done:
		return
	default:
		gw.channel <- v
	}
}

type onceChan struct {
	channel chan interface{}
	wrote   int32
}

func (oc *onceChan) write(val interface{}) {
	if atomic.AddInt32(&oc.wrote, 1) > 1 {
		return
	}

	oc.channel <- val
}

```