<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Goland 随机数](#goland-%E9%9A%8F%E6%9C%BA%E6%95%B0)
  - [math/rand 伪随机](#mathrand-%E4%BC%AA%E9%9A%8F%E6%9C%BA)
    - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [crypto/rand真随机](#cryptorand%E7%9C%9F%E9%9A%8F%E6%9C%BA)
  - [第三方实现](#%E7%AC%AC%E4%B8%89%E6%96%B9%E5%AE%9E%E7%8E%B0)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Goland 随机数

随机数分为真随机和伪随机，

- 一般来说，真随机是指通过物理现象产生的，比如掷骰子、双色球摇奖（前提是机器没问题）、同位素衰变等等，其特征是不可预测、不可重现。

- 而伪随机数一般是通过软件算法产生，看上去是随机的，但是无论是什么算法函数，都会有输入和输出，如果你能得到其输入，自然也可以预测其输出。


## math/rand 伪随机
官方注释
```go
// Package rand implements pseudo-random number generators.
//
// Random numbers are generated by a Source. Top-level functions, such as
// Float64 and Int, user a default shared Source that produces a deterministic
// sequence of values each time a program is run. Use the Seed function to
// initialize the default Source if different behavior is required for each run.
// The default Source is safe for concurrent user by multiple goroutines, but
// Sources created by NewSource are not.
//
// Mathematical interval notation such as [0, n) is used throughout the
// documentation for this package.
//
// For random numbers suitable for security-sensitive work, see the crypto/rand
// package.
```
翻译： 这是一个伪随机数生成器，对于高阶函数比如Float64或者Int，每次运行的时候它使用一个默认共享源来产生一个随机数。

如果你需要每次运行产生的结果不一样，那么就需要使用Seed函数初始化默认源。默认的源是协程并发安全的，但是使用NewSource函数产生的并不是。

最后，对于安全十分敏感的工作，推荐使用crypto/rand包。（注：因为这个包可以产生真随机数）

### 源码分析

结构体
```go
// /go1.20/src/math/rand/rand.go
type Rand struct {
	src Source
	s64 Source64 // non-nil if src is source64

    // readVal包含用于字节的63位整数的remainer，在最近的读取调用期间生成。
    //  它被保存，以便下一个读调用可以从上一个读调用结束的地方开始。
	readVal int64
    //  readPos表示仍然有效的readVal的低位字节数。 
	readPos int8
}
```


```go
// Int returns a non-negative pseudo-random int from the default Source.
func Int() int { return globalRand.Int() }

var globalRand = New(&lockedSource{src: NewSource(1).(*rngSource)})

type lockedSource struct {
    lk  sync.Mutex  //锁
    src *rngSource
}

type rngSource struct {
	tap  int           // index into vec
	feed int           // index into vec
	vec  [rngLen]int64 // current feedback register
}
const (
    rngLen   = 607
    rngTap   = 273
    rngMax   = 1 << 63
    rngMask  = rngMax - 1
    int32max = (1 << 31) - 1
)

func NewSource(seed int64) Source {
    var rng rngSource
    rng.Seed(seed)
    return &rng
}
```
它们默认使用全局的源，而这个源的默认种子是固定的1.


为了每次运行产生不同的结果，我们需要使用一个随机数当种子来初始化源，最常见的做法就是使用当前时间戳：
```go
rand.Seed(time.Now().UnixNano())
```
分析
```go
func (rng *rngSource) Seed(seed int64) {
	rng.tap = 0
	rng.feed = rngLen - rngTap

	seed = seed % int32max
	if seed < 0 {
		seed += int32max
	}
	if seed == 0 {
		seed = 89482311
	}

	x := int32(seed)
	for i := -20; i < rngLen; i++ {
		x = seedrand(x)
		if i >= 0 {
			var u int64
			u = int64(x) << 40
			x = seedrand(x)
			u ^= int64(x) << 20
			x = seedrand(x)
			u ^= int64(x)
			u ^= rngCooked[i]
			// 对 rng.vec 各个位置设置对应的值. rng.vec 的大小是 607.
			rng.vec[i] = u
		}
	}
}


// seed rng x[n+1] = 48271 * x[n] mod (2**31 - 1)
func seedrand(x int32) int32 {
	const (
		A = 48271
		Q = 44488
		R = 3399
	)

	hi := x / Q
	lo := x % Q
	x = A*lo - R*hi
	if x < 0 {
		x += int32max
	}
	return x
}
```

我们在使用不管调用 Intn(), Int31n(), Int63(), Int63n() 等其他函数, 最终调用到就是这个函数 rngSource.Uint64().
```go
// Uint64 returns a non-negative pseudo-random 64-bit integer as an uint64.
func (rng *rngSource) Uint64() uint64 {
	rng.tap--
	if rng.tap < 0 {
		rng.tap += rngLen
	}

	rng.feed--
	if rng.feed < 0 {
		rng.feed += rngLen
	}

	// rng.feed rng.tap 从 rng.vec 中取到两个值相加的结果
	x := rng.vec[rng.feed] + rng.vec[rng.tap]
	rng.vec[rng.feed] = x
	return uint64(x)
}

```

在这里需要注意使用 rng.go 的 rngSource 时, 由于 rng.vec 在获取随机数时会同时设置 rng.vec 的值, 当多 goroutine 同时调用时就会有数据竞争问题. math/rand 采用在调用 rngSource 时加锁  sync.Mutex 解决.
```go
func (r *lockedSource) Uint64() (n uint64) {
    r.lk.Lock()
    n = r.src.Uint64()
    r.lk.Unlock()
    return
}
```

如果你真的要追求极致性能的话，你可能需要自己New一个rand，因为默认的Source为了实现并发安全使用了一个全局的排它锁，必然会带来性能损耗，
如果确实特别在意这点性能消耗的话，可以通过定义一个你的包共享的或者结构体实例共享的 Rand 实例来优化锁的性能消耗（最小化锁的粒度，不跟其他包/代码竞争这个锁）

应该知道使用时需要避开的坑。

（1）相同种子，每次运行的结果是一样的。 因为随机数是从 rng.vec 数组中取出来的，这个数组是根据种子生成的，相同的种子生成的 rng.vec 数组是相同的。

（2）不同种子，每次运行的结果可能一样。 因为根据种子生成 rng.vec 数组时会有一个取模的操作，模后的结果可能相同，导致 rng.vec 数组相同。

（3）rand.New 初始化出来的 rand 不是并发安全的。 因为每次利用 rng.feed, rng.tap 从 rng.vec 中取到随机值后会将随机值重新放入 rng.vec。如果想并发安全，可以使用全局的随机数发生器 rand.globalRand。

（4）不同种子，随机序列发生碰撞的概率高于单个碰撞概率的乘积。 这


## crypto/rand真随机

rand.Reader： 一个全局、共享的加密安全的伪随机数生成器
```go
// Reader is a global, shared instance of a cryptographically
// secure random number generator.
//
// On Linux and FreeBSD, Reader uses getrandom(2) if available, /dev/urandom otherwise.
// On OpenBSD, Reader uses getentropy(2).
// On other Unix-like systems, Reader reads from /dev/urandom.
// On Windows systems, Reader uses the CryptGenRandom API.
// On Wasm, Reader uses the Web Crypto API.
var Reader io.Reader

// Read is a helper function that calls Reader.Read using io.ReadFull.
// On return, n == len(b) if and only if err == nil.
func Read(b []byte) (n int, err error) {
	return io.ReadFull(Reader, b)
}
```
使用系统底层提供的随机数生成器产生加密安全的随机数。这意味着通过rand.Reader生成的随机数在理论上是无法预测的，非常适合用于加密、安全认证等领域。
- 在Linux平台下，使用 getrandom(2) ，在int32位的机器上面上每次最多可以获取2^25-1个字节的数据。
- 在Windows系统中，rand.Reader使用CryptGenRandom函数，这是Windows为开发者提供的用来生成随机数的API。


Linux提供了两种主要的随机数生成设备文件：/dev/random和/dev/urandom。/dev/random是一个阻塞型的随机数生成器。/dev/urandom是一个非阻塞型的随机数生成器.
/dev/random 和 /dev/urandom 都使用熵池来生成随机数，但它们的行为方式有所不同。/dev/random 会在熵池中的熵低于一定值时阻塞等待熵的增加，而 /dev/urandom 不会阻塞等待熵，而是使用伪随机数生成器来生成随机数。
```shell
# 这将生成10个随机字节并将它们转换为Base64编码，以便更容易阅读和使用。
$ head -c 10 /dev/random | base64
$ head -c 10 /dev/urandom | base64
```

```go
// go1.21.5/src/crypto/rand/rand_unix.go

const urandomDevice = "/dev/urandom"

func init() {
	if boring.Enabled {
		Reader = boring.RandReader
		return
	}
	Reader = &reader{}
}


// altGetRandom if non-nil specifies an OS-specific function to get
// urandom-style randomness.
var altGetRandom func([]byte) (err error)

func warnBlocked() {
	println("crypto/rand: blocked for 60 seconds waiting to read random data from the kernel")
}

func (r *reader) Read(b []byte) (n int, err error) {
	boring.Unreachable()
	if r.used.CompareAndSwap(0, 1) {
		// First use of randomness. Start timer to warn about
		// being blocked on entropy not being available.
		t := time.AfterFunc(time.Minute, warnBlocked)
		defer t.Stop()
	}
	// 如果linux,这里指支持 getrandom(2)  
	if altGetRandom != nil && altGetRandom(b) == nil {
		return len(b), nil
	}
	if r.used.Load() != 2 {
		r.mu.Lock()
		if r.used.Load() != 2 {
			f, err := os.Open(urandomDevice)
			if err != nil {
				r.mu.Unlock()
				return 0, err
			}
			r.f = hideAgainReader{f}
			r.used.Store(2)
		}
		r.mu.Unlock()
	}
	return io.ReadFull(r.f, b)
}
```


```go
//go:build dragonfly || freebsd || linux || netbsd || solaris

package rand

import (
	"internal/syscall/unix"
	"runtime"
	"syscall"
)

func init() {
	var maxGetRandomRead int
	switch runtime.GOOS {
	case "linux", "android":
		// Per the manpage:
		//     When reading from the urandom source, a maximum of 33554431 bytes
		//     is returned by a single call to getrandom() on systems where int
		//     has a size of 32 bits.
		maxGetRandomRead = (1 << 25) - 1
	case "dragonfly", "freebsd", "illumos", "netbsd", "solaris":
		maxGetRandomRead = 1 << 8
	default:
		panic("no maximum specified for GetRandom")
	}
	altGetRandom = batched(getRandom, maxGetRandomRead)
}

// If the kernel is too old to support the getrandom syscall(),
// unix.GetRandom will immediately return ENOSYS and we will then fall back to
// reading from /dev/urandom in rand_unix.go. unix.GetRandom caches the ENOSYS
// result so we only suffer the syscall overhead once in this case.
// If the kernel supports the getrandom() syscall, unix.GetRandom will block
// until the kernel has sufficient randomness (as we don't use GRND_NONBLOCK).
// In this case, unix.GetRandom will not return an error.
func getRandom(p []byte) error {
	n, err := unix.GetRandom(p, 0)
	if err != nil {
		return err
	}
	if n != len(p) {
		return syscall.EIO
	}
	return nil
}
```


## 第三方实现
- github.com/joway/fastrand





## 参考

- [优化go生成随机数](https://juejin.cn/post/7071979385787514911)
- [关于 /dev/urandom 的流言终结](https://juejin.cn/post/6844903838256726023) 
