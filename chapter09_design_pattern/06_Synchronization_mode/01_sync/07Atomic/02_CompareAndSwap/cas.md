<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [cas算法](#cas%E7%AE%97%E6%B3%95)
  - [案例](#%E6%A1%88%E4%BE%8B)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->




# cas算法
CAS算法是一种有名的无锁算法。无锁编程，即不使用锁的情况下实现多线程之间的变量同步，也就是在没有线程被阻塞的情况下实现变量的同步，
所以也叫非阻塞同步（Non-blocking Synchronization）。CAS算法涉及到三个操作数
1。需要读写的内存值V
2。进行比较的值A
3。拟写入的新值B
当且仅当 V 的值等于 A时，CAS通过原子方式用新值B来更新V的值，否则不会执行任何操作（比较和替换是一个原子操作）。一般情况下是一个自旋操作，即不断的重试

## 案例
```go
// CompareAndSwapUint32 executes the compare-and-swap operation for a uint32 value.
// Consider using the more ergonomic and less error-prone [Uint32.CompareAndSwap] instead.
func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)
```
函数会先判断参数addr指向的值与参数old是否相等，如果相等，则用参数new替换参数addr的值。最后返回swapped是否替换成功。


