// 不使用sync.pool 的情况
package main

/*
 Get函数：如果没有的话， 就从共享池子里拿，如果还是没有的话，就调用 getSlow 里拿，再不行的话，就调用 New 函数了
 */
import (
	"bytes"
	"sync"
	"testing"
)
var (
	// 声明一个全局变量（或者局部变量也可以）用于存储内存池
	bytesPool = sync.Pool{
		New: func() interface{} { return bytes.Buffer{} },
	}
)

// NewBufferFromPool new bytes.Buffer from sync.Pool
func NewBufferFromPool() bytes.Buffer {
	return bytesPool.Get().(bytes.Buffer) // 通过Get来获得一个
}

// NewBuffer return new bytes.Buffer
func NewBuffer() bytes.Buffer {
	return bytes.Buffer{}
}

func BenchmarkNewBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewBuffer()
		_ = p
	}
}

func BenchmarkNewBufferFromPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewBufferFromPool()
		bytesPool.Put(p)
	}
}

/*
$ go test -bench .
goos: linux
goarch: amd64
pkg: github.com/jiajunhuang/test
BenchmarkNewBuffer-8           	1000000000	         1.14 ns/op
BenchmarkNewBufferFromPool-8   	10978279	       207 ns/op
PASS
ok  	github.com/jiajunhuang/test	3.645s
// 分析
用了Pool比不用还要慢。为啥呢？
	经过我的测试发现主要是在类型转换上比较费时，如果去掉这个，就会快很多，但是 去掉类型转换之后，用了Pool还是比不用更慢，这又是为啥呢？
	因为bytes所使用的内存比较小，使用内存池的效果并不好

// 结论：
1。 类型转换(type casting)很费CPU
2。 对于大块的内存，使用内存池才有意义
 */