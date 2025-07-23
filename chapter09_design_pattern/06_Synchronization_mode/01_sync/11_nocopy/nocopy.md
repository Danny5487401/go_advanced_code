<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [NoCopy机制](#nocopy%E6%9C%BA%E5%88%B6)
  - [go vet 检查](#go-vet-%E6%A3%80%E6%9F%A5)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->



## NoCopy机制

noCopy 是 go1.7 开始引入的一个静态检查机制。它不仅仅工作在运行时或标准库，同时也对用户代码有效
强调no copy 的原因是为了安全，因为结构体对象中包含指针对象的话，直接赋值拷贝是浅拷贝，是极不安全的.

### go vet 检查
go vet工具来检查，那么这个对象必须实现sync.Locker
```go
// A Locker represents an object that can be locked and unlocked.
type Locker interface {
    Lock()
    Unlock()
}

// noCopy 用于嵌入一个结构体中来保证其第一次使用后不会被复制
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
```


```go
func main() {
    var wg sync.WaitGroup
    // 由于是值复制的，run函数内部的wg和main函数内的wg不是同一个，会出现直接打印done的情况
    run(wg)
    wg.Wait()
    println("done")
}

func run(wg sync.WaitGroup) {
    for i := 0; i < 10; i++ {
        go func(num int) {
            wg.Add(1)
            println(num)
            wg.Done()
        }(i)
    }
}
```
