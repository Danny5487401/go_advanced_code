<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [sync.Once](#synconce)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [知识点](#%E7%9F%A5%E8%AF%86%E7%82%B9)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


# sync.Once
始化单例资源，或者并发访问只需初始化一次的共享资源，或者在测试的时候初始化一次测试资源。

sync.Once 只暴露了一个方法 Do，你可以多次调用 Do 方法，但是只有第一次调用 Do 方法时 f 参数才会执行，这里的 f 是一个无参数无返回值的函数。

## 源码分析
```go

//源码分析:sync/once.go

type Once struct {
   done uint32 // 初始值为0表示还未执行过，1表示已经执行过
   m    Mutex
}
func (o *Once) Do(f func()) {
   // 判断done是否为0，若为0，表示未执行过，调用doSlow()方法初始化
   if atomic.LoadUint32(&o.done) == 0 {
      // Outlined slow-path to allow inlining of the fast-path.
      o.doSlow(f)
   }
}

// 加载资源
func (o *Once) doSlow(f func()) {
   o.m.Lock()
   defer o.m.Unlock()
   // 采用双重检测机制 加锁判断done是否为零
   if o.done == 0 {
      // 执行完f()函数后，将done值设置为1
      defer atomic.StoreUint32(&o.done, 1)
      // 执行传入的f()函数
      f()
   }
}

```

### 知识点
1. Do方法中使用atomic.LoadUint32(&o.done) == 0的意义:

如果直接只用o.done = 0 会导致无法及时观察doSlow对o.done的值设置
```go
Programs that modify data being simultaneously accessed by multiple goroutines must serialize such access.

To serialize access, protect the data with channel operations or other synchronization primitives such as those in the sync and sync/atomic packages.   

```
翻译：说的是当一个变量被多个gorouting访问的时候，必须要保证他们是有序的，可以使用sync或者sync/atomic包来实现。用LoadUint32可以保证doSlow设置o.done后可以及时的被取到


2. 但在doSlow函数中却直接使用了o.done==0？

    这里使用了互斥锁.

3. 已经使用锁了，为什么还用StoreUint32来赋值?

   互斥锁只能保证临界区内的操作是可观测的，即只有处于lock和Unlock之间的代码对o.done是可观测的，但在Do中对o.done访问就可以会出现观测不到的情况。因此通过store来保证原子性

4. 这里done字段为什么使用uint32而不适用uint8或者bool?
   在atomic包中没有提供LoadUint8或者LoadBool的操作，通过注释：提到了 hot path，即Do方法的调用是高频的，而每次调用访问done，done位于结构体的第一个字段，可以通过结构体指针直接进行访问。不用增加偏移量。
