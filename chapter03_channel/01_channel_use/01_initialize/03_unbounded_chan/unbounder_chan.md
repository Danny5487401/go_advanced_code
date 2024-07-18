<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/smallnest/chanx 无限缓存的channel](#githubcomsmallnestchanx-%E6%97%A0%E9%99%90%E7%BC%93%E5%AD%98%E7%9A%84channel)
  - [特点](#%E7%89%B9%E7%82%B9)
  - [实现](#%E5%AE%9E%E7%8E%B0)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/smallnest/chanx 无限缓存的channel


Go语言的channel有两种类型，一种是无缓存的channel，一种是有缓存的buffer，这两种类型的channel大家都比较熟悉了，但是对于有缓存的channel,它的缓存长度在创建channel的时候就已经确定了，中间不能扩缩容，这导致在一些场景下使用有问题，或者说不太适合特定的场景。

## 特点
- 不会阻塞write: 它总是能处理write的数据，或者放入到待读取的channel中，或者放入到缓存中
- 无数据时read会被阻塞:当没有可读的数据时，从channel中读取的goroutine会被阻塞
- 读写都是通过channel操作: 内部的缓存不会暴露出来
- 能够查询当前待读取的数据数量:因为缓存中可能也有待处理的数据，所以需要返回len(buffer)+len(chan)
- 关闭channel后，还未读取的channel还是能够被读取，读取完之后才能发现channel已经完毕。这和正常的channel的逻辑是一样的，这种情况叫"drain"未读的数据

## 实现

结构体
```go
type UnboundedChan[T any] struct {
	bufCount int64
	In       chan<- T       // 写 
	Out      <-chan T       // 读
	buffer   *RingBuffer[T] // buffer
}

```

初始化
```go
func NewUnboundedChanSize[T any](ctx context.Context, initInCapacity, initOutCapacity, initBufCapacity int) *UnboundedChan[T] {
	in := make(chan T, initInCapacity)
	out := make(chan T, initOutCapacity)
	ch := UnboundedChan[T]{In: in, Out: out, buffer: NewRingBuffer[T](initBufCapacity)}

	go process(ctx, in, out, &ch)

	return &ch
}


func NewRingBuffer[T any](initialSize int) *RingBuffer[T] {
	if initialSize <= 0 {
		panic("initial size must be great than zero")
	}
	// initial size must >= 2
	if initialSize == 1 {
		initialSize = 2
	}

	return &RingBuffer[T]{
		buf:         make([]T, initialSize),
		initialSize: initialSize,
		size:        initialSize, // 初始化环形大小
	}
}
```


主要处理流程

```go
func process[T any](ctx context.Context, in, out chan T, ch *UnboundedChan[T]) {
	defer close(out)
	drain := func() {
		for !ch.buffer.IsEmpty() {
			select {
			case out <- ch.buffer.Pop():
				atomic.AddInt64(&ch.bufCount, -1)
			case <-ctx.Done():
				return
			}
		}

		ch.buffer.Reset()
		atomic.StoreInt64(&ch.bufCount, 0)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case val, ok := <-in:
			if !ok { // in is closed
				drain()
				return
			}
			// 数据开始写入

			// make sure values' order
			// buffer has some values
			if atomic.LoadInt64(&ch.bufCount) > 0 {
				//  代表 buffer 不为空
				ch.buffer.Write(val)
				atomic.AddInt64(&ch.bufCount, 1)
			} else {
				// out is not full
				select {
				case out <- val:
					// 直接读取
					continue
				default:
				}
                
				// out is full
				ch.buffer.Write(val)
				atomic.AddInt64(&ch.bufCount, 1)
			}

			for !ch.buffer.IsEmpty() {
				// buffer 里面有数据
				select {
				case <-ctx.Done():
					return
				case val, ok := <-in:
					// 写入优先
					if !ok { // in is closed
						drain()
						return
					}
					ch.buffer.Write(val)
					atomic.AddInt64(&ch.bufCount, 1)

				case out <- ch.buffer.Peek():
					// 移动读取 r+1
					ch.buffer.Pop()
					atomic.AddInt64(&ch.bufCount, -1)
					if ch.buffer.IsEmpty() && ch.buffer.size > ch.buffer.initialSize { // after burst
						ch.buffer.Reset()
						atomic.StoreInt64(&ch.bufCount, 0)
					}
				}
			}
		}
	}
}
```

考虑需要写入 buffer 时的扩容

```go
func (r *RingBuffer[T]) Write(v T) {
	r.buf[r.w] = v
	r.w++

	if r.w == r.size {
		// 到初始化地点
		r.w = 0
	}

	if r.w == r.r { // full
		r.grow()
	}
}

// 扩容
func (r *RingBuffer[T]) grow() {
	var size int
	if r.size < 1024 {
		// 两倍
		size = r.size * 2
	} else {
		// 1/4 增长
		size = r.size + r.size/4
	}

	buf := make([]T, size)

	copy(buf[0:], r.buf[r.r:])
	copy(buf[r.size-r.r:], r.buf[0:r.r])

	r.r = 0
	r.w = r.size
	r.size = size
	r.buf = buf
}

```
