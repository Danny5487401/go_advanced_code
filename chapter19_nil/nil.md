<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [nil](#nil)
  - [nil 无默认类型](#nil-%E6%97%A0%E9%BB%98%E8%AE%A4%E7%B1%BB%E5%9E%8B)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# nil

nil 是Go语言内预置的标识符,nil代表了指针pointer、通道channel、函数func、接口interface、map、切片slice类型变量的零值

```go
	// nil is a predeclared identifier representing the zero value for a
	// pointer, channel, func, interface, map, or slice type.
	var nil Type // Type must be a pointer, channel, func, interface, map, or slice type

	// Type is here for the purposes of documentation only. It is a stand-in
	// for any Go type, but represents the same type for any given function
	// invocation.
	type Type int
```


## nil 无默认类型

除了nil以外，所有的Go语言预置标识符都有一个默认类型，如 iota 的默认类型为 int。但nil 是个例外，预置的 nil 是唯一一个无默认类型的值。编译器需要足够的信息来判断一个 nil 值对应的类型。

```go
    _ = (*struct{})(nil)
    _ = []int(nil)
    _ = map[int]bool(nil)
    _ = chan string(nil)
    _ = (func())(nil)
    _ = interface{}(nil)

    // These lines are equivalent to the above lines.
    var _ *struct{} = nil
    var _ []int = nil
    var _ map[int]bool = nil
    var _ chan string = nil
    var _ func() = nil
    var _ interface{} = nil


```

```go
// 无法通过编译
var _ = nil
```