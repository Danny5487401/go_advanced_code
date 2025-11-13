<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [CanSet, CanAddr](#canset-canaddr)
  - [canSet](#canset)
    - [场景一](#%E5%9C%BA%E6%99%AF%E4%B8%80)
    - [场景二](#%E5%9C%BA%E6%99%AF%E4%BA%8C)
    - [场景三](#%E5%9C%BA%E6%99%AF%E4%B8%89)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# CanSet, CanAddr

```go
type flag uintptr

const (
	flagKindWidth        = 5 // there are 27 kinds
	flagKindMask    flag = 1<<flagKindWidth - 1
	flagStickyRO    flag = 1 << 5 // 32
	flagEmbedRO     flag = 1 << 6 // 64
	flagIndir       flag = 1 << 7
	flagAddr        flag = 1 << 8
	flagMethod      flag = 1 << 9
	flagMethodShift      = 10
	flagRO          flag = flagStickyRO | flagEmbedRO
)


func (v Value) CanAddr() bool {
	return v.flag&flagAddr != 0
}

func (v Value) CanSet() bool {
	return v.flag&(flagAddr|flagRO) == flagAddr
}

```
他们的区别就是是否判断 flagRO 的两个位。所以他们的不同换句话说就是“判断这个 Value 是否是私有属性”，私有属性是只读的。不能Set。

## canSet


### 场景一
```go
var x float64 = 3.4
v := reflect.ValueOf(x)
fmt.Println(v.CanSet()) // false
```
golang 里面的所有函数调用都是值复制，所以这里在调用 reflect.ValueOf 的时候，已经复制了一个 x 传递进去了，这里获取到的 v 是一个 x 复制体的 value。那么这个时候，我们就希望知道我能不能通过 v 来设置这里的 x 变量。就需要有个方法来辅助我们做这个事情： CanSet()

但是, 非常明显，由于我们传递的是 x 的一个复制，所以这里根本无法改变 x 的值。这里显示的就是 false。



### 场景二

```go
var x float64 = 3.4
v := reflect.ValueOf(&x)
fmt.Println(v.CanSet()) // false
```

我们将 x 变量的地址传递给 reflect.ValueOf 了。应该是 CanSet 了吧。但是这里却要注意一点，这里的 v 指向的是 x 的指针。所以 CanSet 方法判断的是 x 的指针是否可以设置。指针是肯定不能设置的，所以这里还是返回 false。


### 场景三
```go
var x float64 = 3.4
v := reflect.ValueOf(&x)
fmt.Println(v.Elem().CanSet()) // true
```
reflect 提供了 Elem() 方法来获取这个“指针指向的元素




## 参考

- [一篇理解什么是CanSet, CanAddr](https://studygolang.com/articles/31306?hmsr=joyk.com&utm_source=joyk.com&utm_medium=referral)