## 案例
```go
package main

import "fmt"

func sum(a, b int) {
	c := a + b
	fmt.Println("sum:", c)
}

func f(a, b int) {
	defer sum(a, b)

	fmt.Printf("a: %d, b: %d\n", a, b)
}

func main() {
	a, b := 1, 2
	f(a, b)
}
```
### deferproc 函数
runtime/panic.go
```go
// Create a new deferred function fn with siz bytes of arguments.
// The compiler turns a defer statement into a call to this.
//go:nosplit
func deferproc(siz int32, fn *funcval)
```
    deferproc 函数的第一个参数 siz 是 defered 函数（比如本例中的 sum 函数）的参数以字节为单位的大小，第二个参数 funcval 是一个变长结构体：
proc/runtime2.go
```go
type funcval struct {
    fn uintptr
    // variable-size, fn-specific data here
}
```
例子中的 defer sum(a, b) 大致等价于
```shell
deferproc(16, &funcval{sum})
```
因为 sum 函数有 2 个 int 型的参数共 16 字节，所以在调用 deferproc 函数时第一个参数为16，第二个参数 funcval 结构体对象的 fn 成员为 sum 函数的地址。
我们可以先想一下为什么需要把 sum 函数的参数大小传递给 deferproc() 函数？另外为什么没看到 sum 函数需要的两个参数呢？